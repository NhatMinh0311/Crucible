// Package docker chạy môi trường của phiên bằng container Docker qua Docker Engine API (CR-DOC-03 mục 2.3). Container
// hiện chạy bằng runtime mặc định của Docker (runc), chưa có gVisor.
package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/google/uuid"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/driver"
)

// Nhãn gắn trên container và volume do runtime-api quản lý.
const (
	LabelSession   = "crucible.dev/session"
	LabelManagedBy = "crucible.dev/managed-by"
	managedBy      = "runtime-api"
	toolsDir       = "/opt/crucible/bin"
	pidsLimit      = 512
)

// Driver hiện thực driver.Driver bằng Docker Engine API.
type Driver struct {
	cli        *client.Client
	toolsImage string
	log        *slog.Logger

	toolsMu     sync.Mutex
	toolsVolume string
}

// New kết nối Docker Engine theo biến môi trường DOCKER_HOST (mặc định socket cục bộ).
func New(ctx context.Context, toolsImage string, log *slog.Logger) (*Driver, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("tạo Docker client: %w", err)
	}
	if _, err := cli.Ping(ctx, client.PingOptions{}); err != nil {
		return nil, fmt.Errorf("không kết nối được Docker Engine (Docker đã chạy chưa?): %w", err)
	}
	return &Driver{cli: cli, toolsImage: toolsImage, log: log}, nil
}

// Close đóng kết nối tới Docker Engine.
func (d *Driver) Close() error { return d.cli.Close() }

func (d *Driver) ensureImage(ctx context.Context, ref string) (string, error) {
	img, err := d.cli.ImageInspect(ctx, ref)
	if err == nil {
		return img.ID, nil
	}
	if !cerrdefs.IsNotFound(err) {
		return "", fmt.Errorf("đọc image %s: %w", ref, err)
	}
	d.log.InfoContext(ctx, "kéo image", "image", ref)
	pull, err := d.cli.ImagePull(ctx, ref, client.ImagePullOptions{})
	if err != nil {
		return "", &driver.ProvisionError{Detail: fmt.Sprintf("không kéo được image %s: %v", ref, err)}
	}
	defer pull.Close()
	if err := pull.Wait(ctx); err != nil {
		return "", &driver.ProvisionError{Detail: fmt.Sprintf("không kéo được image %s: %v", ref, err)}
	}
	img, err = d.cli.ImageInspect(ctx, ref)
	if err != nil {
		return "", fmt.Errorf("đọc image %s: %w", ref, err)
	}
	return img.ID, nil
}

// ensureTools chép bộ công cụ từ image vào một volume (một lần cho mỗi phiên bản image) để gắn chỉ đọc vào phiên.
func (d *Driver) ensureTools(ctx context.Context) (string, error) {
	d.toolsMu.Lock()
	defer d.toolsMu.Unlock()
	if d.toolsVolume != "" {
		return d.toolsVolume, nil
	}
	id, err := d.ensureImage(ctx, d.toolsImage)
	if err != nil {
		var pe *driver.ProvisionError
		if errors.As(err, &pe) {
			pe.Detail += "; build image công cụ theo images/README.md"
		}
		return "", err
	}
	name := "crucible-tools-" + strings.TrimPrefix(id, "sha256:")[:12]
	if _, err := d.cli.VolumeInspect(ctx, name, client.VolumeInspectOptions{}); err == nil {
		d.toolsVolume = name
		return name, nil
	} else if !cerrdefs.IsNotFound(err) {
		return "", fmt.Errorf("đọc volume %s: %w", name, err)
	}
	if _, err := d.cli.VolumeCreate(ctx, client.VolumeCreateOptions{
		Name: name, Labels: map[string]string{LabelManagedBy: managedBy},
	}); err != nil {
		return "", fmt.Errorf("tạo volume %s: %w", name, err)
	}
	created, err := d.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image:  d.toolsImage,
			Cmd:    []string{toolsDir + "/cp", "-a", toolsDir + "/.", "/out/"},
			Labels: map[string]string{LabelManagedBy: managedBy},
		},
		HostConfig: &container.HostConfig{
			NetworkMode: "none",
			Mounts:      []mount.Mount{{Type: mount.TypeVolume, Source: name, Target: "/out"}},
		},
	})
	if err != nil {
		return "", fmt.Errorf("tạo container chép công cụ: %w", err)
	}
	defer d.cli.ContainerRemove(context.WithoutCancel(ctx), created.ID, client.ContainerRemoveOptions{Force: true})
	if _, err := d.cli.ContainerStart(ctx, created.ID, client.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("chép công cụ: %w", err)
	}
	wait := d.cli.ContainerWait(ctx, created.ID, client.ContainerWaitOptions{Condition: container.WaitConditionNotRunning})
	select {
	case res := <-wait.Result:
		if res.StatusCode != 0 {
			return "", fmt.Errorf("chép công cụ thoát với mã %d", res.StatusCode)
		}
	case err := <-wait.Error:
		return "", fmt.Errorf("chờ chép công cụ: %w", err)
	}
	d.toolsVolume = name
	return name, nil
}

// Provision dựng container của phiên (spec plan mục 3).
func (d *Driver) Provision(ctx context.Context, spec driver.Spec) (string, error) {
	if _, err := d.ensureImage(ctx, spec.Image); err != nil {
		return "", err
	}
	tools, err := d.ensureTools(ctx)
	if err != nil {
		return "", err
	}
	network := "none"
	if spec.Network == "mirror" {
		// Máy phát triển không có máy chủ gói nội bộ; bài cần cài gói ra Internet qua bridge (spec mục 4).
		network = "bridge"
	}
	pids := int64(pidsLimit)
	// Tiến trình mồ côi chưa được thu dọn vẫn chiếm một suất của PidsLimit tới khi hủy
	// phiên, nên PID 1 phải thu dọn. Đặt ở đây để bảo đảm không phụ thuộc image (BR-M1.31).
	reapInit := true
	created, err := d.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name: "crucible-session-" + spec.SessionID.String(),
		Config: &container.Config{
			Image:    spec.Image,
			Hostname: "lab",
			Labels:   map[string]string{LabelSession: spec.SessionID.String(), LabelManagedBy: managedBy},
		},
		HostConfig: &container.HostConfig{
			Resources: container.Resources{
				Memory:     spec.MemoryBytes,
				MemorySwap: spec.MemoryBytes,
				NanoCPUs:   spec.NanoCPUs,
				PidsLimit:  &pids,
			},
			Init:          &reapInit,
			CapAdd:        spec.Capabilities,
			SecurityOpt:   []string{"no-new-privileges"},
			NetworkMode:   container.NetworkMode(network),
			RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyDisabled},
			Mounts:        []mount.Mount{{Type: mount.TypeVolume, Source: tools, Target: toolsDir, ReadOnly: true}},
		},
	})
	if err != nil {
		return "", fmt.Errorf("tạo container: %w", err)
	}
	handle := created.ID
	if _, err := d.cli.ContainerStart(ctx, handle, client.ContainerStartOptions{}); err != nil {
		return handle, fmt.Errorf("khởi động container: %w", err)
	}
	if len(spec.Files) > 0 {
		archive, err := tarFiles(spec.Files)
		if err != nil {
			return handle, &driver.ProvisionError{Detail: err.Error()}
		}
		if _, err := d.cli.CopyToContainer(ctx, handle, client.CopyToContainerOptions{
			DestinationPath: "/", Content: archive,
		}); err != nil {
			return handle, &driver.ProvisionError{Detail: fmt.Sprintf("ghi tệp: %v", err)}
		}
	}
	if strings.TrimSpace(spec.Setup) != "" {
		res, err := d.run(ctx, handle, wrapTimeout(spec.SetupTimeout, "/bin/sh", "-c", spec.Setup), nil, "/root", spec.SetupTimeout)
		if err != nil {
			return handle, err
		}
		if res.TimedOut {
			return handle, &driver.ProvisionError{Detail: fmt.Sprintf("setup quá %s", spec.SetupTimeout)}
		}
		if res.ExitCode != 0 {
			return handle, &driver.ProvisionError{Detail: fmt.Sprintf("setup thoát với mã %d: %s", res.ExitCode, lastLine(res.Stderr, res.Stdout))}
		}
	}
	return handle, nil
}

// Exec chạy lệnh kiểm tra hoặc script khôi phục (CR-DOC-03 mục 7).
func (d *Driver) Exec(ctx context.Context, handle string, e driver.ExecSpec) (driver.ExecResult, error) {
	if e.Purpose == driver.PurposeCheck {
		// Môi trường biến rỗng ngoài PATH, LANG; trình thông dịch và công cụ lấy từ volume chỉ đọc (ADR-0011).
		cmd := append([]string{toolsDir + "/env", "-i", "PATH=" + toolsDir, "LANG=C.UTF-8"},
			wrapTimeout(e.Timeout, toolsDir+"/sh", "-c", e.Script)...)
		return d.run(ctx, handle, cmd, nil, "/", e.Timeout)
	}
	return d.run(ctx, handle, wrapTimeout(e.Timeout, "/bin/sh", "-c", e.Script), nil, "/root", e.Timeout)
}

// wrapTimeout bọc lệnh bằng timeout của BusyBox để tiến trình bị dừng hẳn khi quá hạn.
func wrapTimeout(timeout time.Duration, cmd ...string) []string {
	secs := int(timeout.Round(time.Second) / time.Second)
	if secs < 1 {
		secs = 1
	}
	return append([]string{toolsDir + "/timeout", "-s", "KILL", strconv.Itoa(secs)}, cmd...)
}

func (d *Driver) run(ctx context.Context, handle string, cmd, env []string, workdir string, timeout time.Duration) (driver.ExecResult, error) {
	start := time.Now()
	created, err := d.cli.ExecCreate(ctx, handle, client.ExecCreateOptions{
		User: "root", Cmd: cmd, Env: env, WorkingDir: workdir, AttachStdout: true, AttachStderr: true,
	})
	if err != nil {
		return driver.ExecResult{}, fmt.Errorf("tạo exec: %w", err)
	}
	attached, err := d.cli.ExecAttach(ctx, created.ID, client.ExecAttachOptions{})
	if err != nil {
		return driver.ExecResult{}, fmt.Errorf("gắn exec: %w", err)
	}
	defer attached.Close()

	// Lệnh bị timeout của BusyBox dừng đúng hạn; giới hạn phía Go chỉ để không chờ vô hạn khi có sự cố.
	runCtx, cancel := context.WithTimeout(ctx, timeout+5*time.Second)
	defer cancel()
	stdout, stderr := &capped{max: driver.MaxOutput}, &capped{max: driver.MaxOutput}
	done := make(chan error, 1)
	go func() {
		_, err := stdcopy.StdCopy(stdout, stderr, attached.Reader)
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			return driver.ExecResult{}, fmt.Errorf("đọc đầu ra exec: %w", err)
		}
	case <-runCtx.Done():
		attached.Close()
		return driver.ExecResult{ExitCode: -1, TimedOut: true, Stdout: stdout.String(), Stderr: stderr.String(),
			Duration: time.Since(start)}, nil
	}
	insp, err := d.cli.ExecInspect(ctx, created.ID, client.ExecInspectOptions{})
	if err != nil {
		return driver.ExecResult{}, fmt.Errorf("đọc kết quả exec: %w", err)
	}
	res := driver.ExecResult{ExitCode: insp.ExitCode, Stdout: stdout.String(), Stderr: stderr.String(), Duration: time.Since(start)}
	if insp.ExitCode == 137 && res.Duration >= timeout {
		res.ExitCode, res.TimedOut = -1, true
	}
	return res, nil
}

// Attach mở shell đăng nhập có TTY: bash nếu có, nếu không thì sh.
func (d *Driver) Attach(ctx context.Context, handle string, cols, rows uint) (driver.Stream, error) {
	created, err := d.cli.ExecCreate(ctx, handle, client.ExecCreateOptions{
		User: "root", TTY: true, AttachStdin: true, AttachStdout: true, AttachStderr: true,
		ConsoleSize: client.ConsoleSize{Height: rows, Width: cols},
		Env:         []string{"TERM=xterm-256color"},
		WorkingDir:  "/root",
		Cmd:         []string{"/bin/sh", "-c", "if [ -x /bin/bash ]; then exec /bin/bash -l; else exec /bin/sh -l; fi"},
	})
	if err != nil {
		return nil, fmt.Errorf("tạo shell: %w", err)
	}
	attached, err := d.cli.ExecAttach(ctx, created.ID, client.ExecAttachOptions{
		TTY: true, ConsoleSize: client.ConsoleSize{Height: rows, Width: cols},
	})
	if err != nil {
		return nil, fmt.Errorf("gắn shell: %w", err)
	}
	return &stream{cli: d.cli, execID: created.ID, resp: attached.HijackedResponse}, nil
}

type stream struct {
	cli    *client.Client
	execID string
	resp   client.HijackedResponse
}

func (s *stream) Read(p []byte) (int, error)  { return s.resp.Reader.Read(p) }
func (s *stream) Write(p []byte) (int, error) { return s.resp.Conn.Write(p) }
func (s *stream) Close() error                { s.resp.Close(); return nil }
func (s *stream) Resize(ctx context.Context, cols, rows uint) error {
	_, err := s.cli.ExecResize(ctx, s.execID, client.ExecResizeOptions{Height: rows, Width: cols})
	return err
}

// Destroy xóa container; container không còn thì bỏ qua.
func (d *Driver) Destroy(ctx context.Context, handle string) error {
	_, err := d.cli.ContainerRemove(ctx, handle, client.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
	if err != nil && !cerrdefs.IsNotFound(err) {
		return fmt.Errorf("xóa container: %w", err)
	}
	return nil
}

// MemoryUsage trả bộ nhớ làm việc của container (không gồm page cache có thể thu hồi).
func (d *Driver) MemoryUsage(ctx context.Context, handle string) (int64, error) {
	res, err := d.cli.ContainerStats(ctx, handle, client.ContainerStatsOptions{Stream: false})
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	var st container.StatsResponse
	if err := json.NewDecoder(res.Body).Decode(&st); err != nil {
		return 0, err
	}
	usage := st.MemoryStats.Usage
	if inactive, ok := st.MemoryStats.Stats["inactive_file"]; ok && inactive < usage {
		usage -= inactive
	}
	return int64(usage), nil
}

// Watch theo dõi sự kiện die và oom của container phiên; tự kết nối lại khi luồng sự kiện đứt.
func (d *Driver) Watch(ctx context.Context) <-chan driver.Exit {
	out := make(chan driver.Exit, 16)
	go func() {
		defer close(out)
		for ctx.Err() == nil {
			f := make(client.Filters).Add("type", "container").Add("label", LabelManagedBy+"="+managedBy).
				Add("event", string(events.ActionDie), string(events.ActionOOM))
			res := d.cli.Events(ctx, client.EventsListOptions{Filters: f})
		loop:
			for {
				select {
				case <-ctx.Done():
					return
				case err := <-res.Err:
					if ctx.Err() == nil {
						d.log.Warn("luồng sự kiện Docker đứt, kết nối lại", "err", err)
						time.Sleep(time.Second)
					}
					break loop
				case m := <-res.Messages:
					id, err := uuid.Parse(m.Actor.Attributes[LabelSession])
					if err != nil {
						continue
					}
					exit := driver.Exit{SessionID: id, OOM: m.Action == events.ActionOOM}
					if m.Action == events.ActionDie {
						exit.ExitCode, _ = strconv.Atoi(m.Actor.Attributes["exitCode"])
						if insp, err := d.cli.ContainerInspect(ctx, m.Actor.ID, client.ContainerInspectOptions{}); err == nil && insp.Container.State != nil {
							exit.OOM = insp.Container.State.OOMKilled
						}
					}
					select {
					case out <- exit:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return out
}

// List trả container phiên đang tồn tại (kể cả đã dừng).
func (d *Driver) List(ctx context.Context) ([]driver.Instance, error) {
	res, err := d.cli.ContainerList(ctx, client.ContainerListOptions{
		All: true, Filters: make(client.Filters).Add("label", LabelManagedBy+"="+managedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("liệt kê container: %w", err)
	}
	var out []driver.Instance
	for _, c := range res.Items {
		id, err := uuid.Parse(c.Labels[LabelSession])
		if err != nil {
			continue
		}
		out = append(out, driver.Instance{SessionID: id, Handle: c.ID, Running: c.State == container.StateRunning})
	}
	return out, nil
}

// tarFiles đóng gói tệp của bài lab thành luồng tar, kèm thư mục cha.
func tarFiles(files []driver.File) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	dirs := map[string]bool{}
	now := time.Now()
	for _, f := range files {
		p := path.Clean(f.Path)
		if !path.IsAbs(p) || p == "/" {
			return nil, fmt.Errorf("đường dẫn tệp phải tuyệt đối: %q", f.Path)
		}
		var parents []string
		for dir := path.Dir(p); dir != "/" && !dirs[dir]; dir = path.Dir(dir) {
			parents = append([]string{dir}, parents...)
		}
		for _, dir := range parents {
			dirs[dir] = true
			if err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: strings.TrimPrefix(dir, "/") + "/", Mode: 0o755, ModTime: now}); err != nil {
				return nil, err
			}
		}
		if err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeReg, Name: strings.TrimPrefix(p, "/"),
			Mode: int64(f.Mode.Perm()), Size: int64(len(f.Content)), ModTime: now}); err != nil {
			return nil, err
		}
		if _, err := io.WriteString(tw, f.Content); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

// capped giữ tối đa max byte đầu, đọc và bỏ phần còn lại để luồng không bị nghẽn.
type capped struct {
	max int
	buf bytes.Buffer
}

func (c *capped) Write(p []byte) (int, error) {
	if room := c.max - c.buf.Len(); room > 0 {
		if len(p) > room {
			c.buf.Write(p[:room])
		} else {
			c.buf.Write(p)
		}
	}
	return len(p), nil
}

func (c *capped) String() string { return c.buf.String() }

func lastLine(outputs ...string) string {
	for _, o := range outputs {
		lines := strings.Split(strings.TrimRight(o, "\n"), "\n")
		if l := strings.TrimSpace(lines[len(lines)-1]); l != "" {
			return l
		}
	}
	return "không có đầu ra"
}

var _ driver.Driver = (*Driver)(nil)
