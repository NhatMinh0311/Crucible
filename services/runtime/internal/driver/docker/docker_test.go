package docker_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/driver"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/driver/docker"
)

// Image phải được build trước theo images/README.md.
const (
	labImage   = "ghcr.io/crucible-project/crucible-labs/ubuntu-nginx:1.0.0"
	toolsImage = "ghcr.io/crucible-project/crucible/tools:1.0.0"
)

func newDriver(t *testing.T) *docker.Driver {
	t.Helper()
	if os.Getenv("SKIP_DOCKER_TESTS") != "" {
		t.Skip("SKIP_DOCKER_TESTS được đặt")
	}
	d, err := docker.New(context.Background(), toolsImage, slog.Default())
	if err != nil {
		t.Skipf("không có Docker Engine: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	return d
}

func provision(t *testing.T, d *docker.Driver, mutate func(*driver.Spec)) (string, error) {
	t.Helper()
	spec := driver.Spec{
		SessionID: uuid.New(), Image: labImage, MemoryBytes: 256 << 20, NanoCPUs: 1e9, Network: "none",
		SetupTimeout: 30 * time.Second,
	}
	if mutate != nil {
		mutate(&spec)
	}
	h, err := d.Provision(context.Background(), spec)
	if h != "" {
		t.Cleanup(func() { _ = d.Destroy(context.Background(), h) })
	}
	return h, err
}

func check(t *testing.T, d *docker.Driver, h, script string) driver.ExecResult {
	t.Helper()
	res, err := d.Exec(context.Background(), h, driver.ExecSpec{Purpose: driver.PurposeCheck, Script: script, Timeout: 10 * time.Second})
	if err != nil {
		t.Fatalf("exec %q: %v", script, err)
	}
	return res
}

func TestProvisionAndExec(t *testing.T) {
	d := newDriver(t)
	h, err := provision(t, d, func(s *driver.Spec) {
		s.Files = []driver.File{{Path: "/srv/lab/conf/app.conf", Content: "port=8080\n", Mode: 0o640}}
		s.Setup = "useradd -m tester && echo setup-ok > /root/setup.txt"
	})
	if err != nil {
		t.Fatalf("dựng môi trường: %v", err)
	}

	t.Run("D3 tệp, setup, giới hạn tài nguyên, không mạng", func(t *testing.T) {
		res := check(t, d, h, `cat /root/setup.txt; stat -c %a /srv/lab/conf/app.conf; grep -c ^tester: /etc/passwd;
			cat /sys/fs/cgroup/memory.max /sys/fs/cgroup/pids.max; ls /sys/class/net | grep -vc ^lo$`)
		want := "setup-ok\n640\n1\n268435456\n512\n0\n"
		if res.Stdout != want {
			t.Fatalf("mong đợi %q, nhận %q (stderr %q, exit %d)", want, res.Stdout, res.Stderr, res.ExitCode)
		}
	})

	t.Run("D3 BR-M1.31 PID 1 thu dọn tiến trình mồ côi", func(t *testing.T) {
		// Zombie không xuất hiện trong cgroup.procs nhưng vẫn tính trong pids.current,
		// nên PID 1 không thu dọn làm hạn mức tác vụ của phiên hao dần mà không báo lỗi.
		res := check(t, d, h, `cat /proc/1/comm; for i in 1 2 3; do sleep 1 & done; exit 0`)
		if strings.TrimSpace(res.Stdout) != "docker-init" {
			t.Fatalf("PID 1 phải do M1 đặt, không theo image: %q", res.Stdout)
		}
		time.Sleep(3 * time.Second)
		res = check(t, d, h, `awk '{print $3}' /proc/[0-9]*/stat | grep -c '^Z'`)
		if n := strings.TrimSpace(res.Stdout); n != "0" {
			t.Errorf("còn %s tiến trình chưa được thu dọn", n)
		}
	})

	t.Run("D3 không capability bổ sung, công cụ chỉ đọc", func(t *testing.T) {
		res := check(t, d, h, `grep CapEff /proc/self/status; touch /opt/crucible/bin/x`)
		if !strings.Contains(res.Stdout, "00000000a80425fb") && !strings.Contains(res.Stdout, "00000000a80405fb") {
			t.Errorf("CapEff khác mặc định của Docker: %q", res.Stdout)
		}
		if res.ExitCode == 0 || !strings.Contains(res.Stderr, "Read-only") {
			t.Errorf("ghi vào thư mục công cụ phải thất bại: exit %d, %q", res.ExitCode, res.Stderr)
		}
	})

	t.Run("D4 AC-C2.4 lệnh kiểm tra không dùng công cụ của môi trường", func(t *testing.T) {
		restore, err := d.Exec(context.Background(), h, driver.ExecSpec{Purpose: driver.PurposeRestore, Timeout: 10 * time.Second,
			Script: "printf '#!/bin/sh\\nexit 0\\n' > /usr/local/bin/grep && chmod +x /usr/local/bin/grep && echo 'export FOO=1' >> /root/.profile"})
		if err != nil || restore.ExitCode != 0 {
			t.Fatalf("restore: %+v %v", restore, err)
		}
		res := check(t, d, h, `grep -q khong-co /etc/passwd`)
		if res.ExitCode != 1 {
			t.Errorf("grep giả của học viên không được dùng; mong đợi mã 1, nhận %d", res.ExitCode)
		}
		res = check(t, d, h, `env | sort`)
		for _, line := range strings.Split(strings.TrimSpace(res.Stdout), "\n") {
			name, _, _ := strings.Cut(line, "=")
			switch name {
			case "PATH", "LANG", "PWD", "SHLVL":
			default:
				t.Errorf("biến môi trường lạ trong lệnh kiểm tra: %q", line)
			}
		}
	})

	t.Run("D4 quá hạn và cắt đầu ra", func(t *testing.T) {
		start := time.Now()
		res, err := d.Exec(context.Background(), h, driver.ExecSpec{Purpose: driver.PurposeCheck, Script: "sleep 20", Timeout: time.Second})
		if err != nil || !res.TimedOut || res.ExitCode != -1 || time.Since(start) > 8*time.Second {
			t.Errorf("mong đợi hết thời gian sau khoảng 1 giây: %+v %v sau %s", res, err, time.Since(start))
		}
		res = check(t, d, h, `head -c 10000 /dev/zero | tr '\0' a; head -c 6000 /dev/zero | tr '\0' b >&2`)
		if len(res.Stdout) != driver.MaxOutput || len(res.Stderr) != driver.MaxOutput {
			t.Errorf("đầu ra phải cắt ở %d byte: stdout %d, stderr %d", driver.MaxOutput, len(res.Stdout), len(res.Stderr))
		}
	})

	t.Run("terminal tương tác", func(t *testing.T) {
		st, err := d.Attach(context.Background(), h, 80, 24)
		if err != nil {
			t.Fatal(err)
		}
		defer st.Close()
		if err := st.Resize(context.Background(), 120, 40); err != nil {
			t.Errorf("đổi kích thước: %v", err)
		}
		_, _ = io.WriteString(st, "echo xin-$((40+2)); stty size\n")
		var out bytes.Buffer
		buf := make([]byte, 4096)
		deadline := time.Now().Add(10 * time.Second)
		for !strings.Contains(out.String(), "40 120") && time.Now().Before(deadline) {
			n, err := st.Read(buf)
			out.Write(buf[:n])
			if err != nil {
				break
			}
		}
		if !strings.Contains(out.String(), "xin-42") || !strings.Contains(out.String(), "40 120") {
			t.Errorf("không thấy kết quả lệnh trong terminal: %q", out.String())
		}
	})

	t.Run("D11 liệt kê và xóa", func(t *testing.T) {
		list, err := d.List(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, in := range list {
			found = found || (in.Handle == h && in.Running)
		}
		if !found {
			t.Fatalf("không thấy container đang chạy trong danh sách")
		}
	})
}

func TestProvisionFailures(t *testing.T) {
	d := newDriver(t)
	var pe *driver.ProvisionError

	_, err := provision(t, d, func(s *driver.Spec) { s.Setup = "echo dang-cai; echo loi-cau-hinh >&2; exit 3" })
	if !errors.As(err, &pe) || !strings.Contains(pe.Detail, "loi-cau-hinh") || !strings.Contains(pe.Detail, "mã 3") {
		t.Errorf("D10 setup lỗi: %v", err)
	}
	_, err = provision(t, d, func(s *driver.Spec) { s.Setup = "sleep 20"; s.SetupTimeout = time.Second })
	if !errors.As(err, &pe) || !strings.Contains(pe.Detail, "quá") {
		t.Errorf("D10 setup quá hạn: %v", err)
	}
	_, err = provision(t, d, func(s *driver.Spec) { s.Image = "crucible.invalid/khong-co:1" })
	if !errors.As(err, &pe) {
		t.Errorf("image không tồn tại phải là ProvisionError: %v", err)
	}
}

func TestOutOfMemory(t *testing.T) {
	d := newDriver(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	exits := d.Watch(ctx)
	id := uuid.New()
	h, err := provision(t, d, func(s *driver.Spec) { s.SessionID = id; s.MemoryBytes = 64 << 20 })
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_, _ = d.Exec(ctx, h, driver.ExecSpec{Purpose: driver.PurposeRestore, Timeout: 20 * time.Second,
			Script: "head -c 300m /dev/zero | tail > /dev/null"})
	}()
	for {
		select {
		case e := <-exits:
			if e.SessionID == id && e.OOM {
				return
			}
		case <-ctx.Done():
			t.Fatal("D8 không nhận được sự kiện hết bộ nhớ")
		}
	}
}
