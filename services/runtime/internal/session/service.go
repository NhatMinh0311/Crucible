// Package session hiện thực vòng đời phiên và các quy tắc BR-M1.* thuộc phạm vi specs/001-runtime-docker. Tầng này
// không biết môi trường chạy ở đâu; mọi thao tác với môi trường đi qua driver.Driver.
package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/api"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/driver"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/events"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/problem"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/store"
)

// Mức tài nguyên (CR-DOC-04 mục 10).
var profiles = map[api.Profile]struct {
	memory int64
	cpus   int64
}{
	api.Small:  {memory: 256 << 20, cpus: 1},
	api.Medium: {memory: 512 << 20, cpus: 2},
}

// Các mốc thời gian của vòng đời (CR-MOD-M1 BR-M1.9, CR-DOC-02 mục 7.1).
const (
	expiringBefore   = 5 * time.Minute
	idleWarnBefore   = 2 * time.Minute
	ticketTTL        = 30 * time.Second
	checkMaxTimeout  = 15 * time.Second
	sampleInterval   = 15 * time.Second
	activityFlushGap = 30 * time.Second
)

// Lỗi nghiệp vụ (CR-DOC-03 mục 3.4).
var (
	errSessionLimit     = problem.New(http.StatusConflict, "SESSION_LIMIT", "Học viên đã có phiên cùng loại đang hoạt động")
	errNotActive        = problem.New(http.StatusConflict, "SESSION_NOT_ACTIVE", "Thao tác yêu cầu phiên ở pha READY hoặc EXPIRING")
	errExecBusy         = problem.New(http.StatusConflict, "EXEC_BUSY", "Phiên đang chạy một lệnh khác")
	errIdempotencyReuse = problem.New(http.StatusUnprocessableEntity, "IDEMPOTENCY_KEY_REUSED", "Khóa lặp lại an toàn đã dùng cho nội dung khác hoặc đang được xử lý")
)

// Notifier gửi khung trạng thái tới terminal đang mở của phiên.
type Notifier interface {
	Notify(sessionID uuid.UUID, frame any)
	CloseSession(sessionID uuid.UUID, reason string)
}

// Options cấu hình Service.
type Options struct {
	SetupTimeout        time.Duration
	ProvisionTimeout    time.Duration
	AllowedCapabilities []string
	ImagePrefixes       []string
	EventsSource        string
	// Now thay đồng hồ trong kiểm thử; mặc định time.Now theo UTC.
	Now func() time.Time
}

// Service điều phối vòng đời phiên.
type Service struct {
	store    *store.Store
	drv      driver.Driver
	events   events.Writer
	log      *slog.Logger
	opts     Options
	notifier Notifier
	now      func() time.Time

	mu         sync.Mutex
	provisions map[uuid.UUID]context.CancelFunc
	busy       map[uuid.UUID]bool
	activity   map[uuid.UUID]time.Time
	wg         sync.WaitGroup
}

// New tạo Service.
func New(st *store.Store, drv driver.Driver, log *slog.Logger, opts Options) *Service {
	now := opts.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{
		store: st, drv: drv, log: log, opts: opts, events: events.Writer{Source: opts.EventsSource},
		notifier: nopNotifier{}, now: now,
		provisions: map[uuid.UUID]context.CancelFunc{}, busy: map[uuid.UUID]bool{}, activity: map[uuid.UUID]time.Time{},
	}
}

// SetNotifier gắn bộ phát khung trạng thái của terminal.
func (s *Service) SetNotifier(n Notifier) { s.notifier = n }

type nopNotifier struct{}

func (nopNotifier) Notify(uuid.UUID, any)          {}
func (nopNotifier) CloseSession(uuid.UUID, string) {}

// Get đọc phiên; trả problem.ErrNotFound khi không có.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*store.Session, error) {
	sess, err := store.GetSession(ctx, s.store.Pool, id)
	if errors.Is(err, store.ErrNotFound) {
		return nil, problem.ErrNotFound
	}
	return sess, err
}

// CreateResult là phản hồi của tạo phiên; Replayed khi trả lại phản hồi đã lưu theo khóa lặp lại an toàn.
type CreateResult struct {
	Session  api.Session
	Replayed bool
}

// Create ghi phiên PENDING và bắt đầu dựng môi trường (AC-R1.1, BR-M1.1, BR-M1.26, CR-DOC-03 mục 3.3).
func (s *Service) Create(ctx context.Context, caller string, key uuid.UUID, req api.SessionCreate) (CreateResult, error) {
	env := req.Environment
	prof, ok := profiles[env.Profile]
	if !ok {
		return CreateResult{}, problem.New(http.StatusUnprocessableEntity, "ENVIRONMENT_NOT_ALLOWED", "Mức tài nguyên không tồn tại")
	}
	if !slices.ContainsFunc(s.opts.ImagePrefixes, func(p string) bool { return strings.HasPrefix(env.Image, p) }) {
		return CreateResult{}, problem.New(http.StatusUnprocessableEntity, "ENVIRONMENT_NOT_ALLOWED",
			fmt.Sprintf("Image %s ngoài danh sách cho phép (SESSION_IMAGE_PREFIXES)", env.Image))
	}
	if env.Capabilities != nil {
		for _, c := range *env.Capabilities {
			if !slices.Contains(s.opts.AllowedCapabilities, c) {
				return CreateResult{}, problem.New(http.StatusUnprocessableEntity, "ENVIRONMENT_NOT_ALLOWED",
					fmt.Sprintf("Capability %s không được cấp ở môi trường này (SESSION_ALLOWED_CAPABILITIES)", c))
			}
		}
	}
	if _, err := files(env); err != nil {
		return CreateResult{}, problem.Validation(err.Error())
	}
	body, err := json.Marshal(req)
	if err != nil {
		return CreateResult{}, err
	}
	hash := sha256.Sum256(body)
	envJSON, _ := json.Marshal(env)

	now := s.now()
	sess := &store.Session{
		ID: uuid.New(), OwnerUserID: req.OwnerUserId, Priority: string(req.Priority), Profile: string(env.Profile),
		Image: env.Image, QuotaAccountRef: req.QuotaAccountRef, ReservationRef: req.ReservationRef,
		RecordInput: req.RecordInput != nil && *req.RecordInput, Labels: map[string]string{}, Environment: envJSON,
		Phase: store.PhasePending, CreatedAt: now, TTLMinutes: env.TtlMinutes, IdleMinutes: env.IdleMinutes,
		MaxExtensionMinutes: env.MaxExtensionMinutes, MemoryLimitBytes: prof.memory,
	}
	if req.Labels != nil {
		sess.Labels = *req.Labels
	}

	var result CreateResult
	err = s.store.InTx(ctx, func(tx pgx.Tx) error {
		prev, err := store.ClaimIdempotencyKey(ctx, tx, caller, key, "createSession", hash[:])
		if err != nil {
			return err
		}
		if prev != nil {
			if string(prev.RequestHash) != string(hash[:]) || prev.Status == nil {
				return errIdempotencyReuse
			}
			result.Replayed = true
			return json.Unmarshal(prev.Response, &result.Session)
		}
		if err := store.InsertSession(ctx, tx, sess); err != nil {
			if errors.Is(err, store.ErrActiveSessionExists) {
				return errSessionLimit
			}
			return err
		}
		result.Session = ToAPI(sess)
		return store.SaveIdempotentResponse(ctx, tx, caller, key, "createSession", http.StatusAccepted, result.Session)
	})
	if err != nil {
		return CreateResult{}, err
	}
	if !result.Replayed {
		s.startProvision(sess)
	}
	return result, nil
}

func files(env api.RenderedEnvironment) ([]driver.File, error) {
	if env.Files == nil {
		return nil, nil
	}
	out := make([]driver.File, 0, len(*env.Files))
	for _, f := range *env.Files {
		mode, err := strconv.ParseUint(f.Mode, 8, 32)
		if err != nil || mode > 0o7777 {
			return nil, fmt.Errorf("mode của tệp %s không hợp lệ: %q", f.Path, f.Mode)
		}
		out = append(out, driver.File{Path: f.Path, Content: f.Content, Mode: os.FileMode(mode)})
	}
	return out, nil
}

func (s *Service) startProvision(sess *store.Session) {
	ctx, cancel := context.WithTimeout(context.Background(), s.opts.ProvisionTimeout)
	s.mu.Lock()
	s.provisions[sess.ID] = cancel
	s.mu.Unlock()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			cancel()
			s.mu.Lock()
			delete(s.provisions, sess.ID)
			s.mu.Unlock()
		}()
		s.provision(ctx, sess.ID)
	}()
}

func (s *Service) provision(ctx context.Context, id uuid.UUID) {
	log := s.log.With("session", id)
	var env api.RenderedEnvironment
	sess, err := s.transition(ctx, id, func(_ pgx.Tx, sess *store.Session) error {
		if sess.Phase != store.PhasePending {
			return errSkip
		}
		sess.Phase = store.PhaseProvisioning
		return json.Unmarshal(sess.Environment, &env)
	})
	if err != nil {
		if !errors.Is(err, errSkip) {
			log.Error("chuyển PROVISIONING", "err", err)
		}
		return
	}
	fs, _ := files(env)
	spec := driver.Spec{
		SessionID: id, Image: env.Image, MemoryBytes: sess.MemoryLimitBytes, NanoCPUs: profiles[env.Profile].cpus * 1e9,
		Network: string(env.Network), Files: fs, SetupTimeout: s.opts.SetupTimeout,
	}
	if env.Capabilities != nil {
		spec.Capabilities = *env.Capabilities
	}
	if env.Setup != nil {
		spec.Setup = *env.Setup
	}
	start := s.now()
	handle, perr := s.drv.Provision(ctx, spec)

	cleanup := func() {
		if handle != "" {
			if err := s.drv.Destroy(context.Background(), handle); err != nil {
				log.Warn("xóa môi trường", "err", err)
			}
		}
	}
	if perr != nil {
		if ctx.Err() != nil && errors.Is(context.Cause(ctx), context.Canceled) {
			// Phiên đã được kết thúc trong lúc dựng (Terminate hủy ctx).
			cleanup()
			return
		}
		detail := "Lỗi hệ thống khi dựng môi trường"
		var pe *driver.ProvisionError
		switch {
		case errors.As(perr, &pe):
			detail = pe.Detail
		case errors.Is(perr, context.DeadlineExceeded):
			detail = fmt.Sprintf("Dựng môi trường quá %s", s.opts.ProvisionTimeout)
		default:
			log.Error("dựng môi trường", "err", perr)
		}
		_, err := s.transition(context.Background(), id, func(tx pgx.Tx, sess *store.Session) error {
			if sess.Phase != store.PhaseProvisioning {
				return errSkip
			}
			reason := string(api.EndReasonPROVISIONFAILED)
			sess.Phase, sess.EndReason, sess.FailureDetail = store.PhaseFailed, &reason, &detail
			ended := s.now()
			sess.EndedAt = &ended
			_, err := s.events.Write(context.Background(), tx, events.SessionFailed, id.String(), map[string]any{
				"sessionId": id, "endReason": reason, "detail": detail, "labels": sess.Labels,
			})
			return err
		})
		if err != nil && !errors.Is(err, errSkip) {
			log.Error("ghi FAILED", "err", err)
		}
		cleanup()
		return
	}

	_, err = s.transition(context.Background(), id, func(tx pgx.Tx, sess *store.Session) error {
		if sess.Phase != store.PhaseProvisioning {
			return errSkip
		}
		ready := s.now()
		expires := ready.Add(time.Duration(sess.TTLMinutes) * time.Minute)
		sess.Phase, sess.ReadyAt, sess.ExpiresAt, sess.LastActivityAt, sess.ContainerID = store.PhaseReady, &ready, &expires, &ready, &handle
		_, err := s.events.Write(context.Background(), tx, events.SessionReady, id.String(), map[string]any{
			"sessionId": id, "ownerUserId": sess.OwnerUserID, "readyAt": ready, "expiresAt": expires, "labels": sess.Labels,
		})
		return err
	})
	if err != nil {
		if !errors.Is(err, errSkip) {
			log.Error("ghi READY", "err", err)
		}
		cleanup()
		return
	}
	log.Info("phiên sẵn sàng", "provisionMs", s.now().Sub(start).Milliseconds())
}

var errSkip = errors.New("bỏ qua: pha đã thay đổi")

// transition khóa dòng phiên, áp fn rồi ghi lại trong cùng transaction.
func (s *Service) transition(ctx context.Context, id uuid.UUID, fn func(pgx.Tx, *store.Session) error) (*store.Session, error) {
	var out *store.Session
	err := s.store.InTx(ctx, func(tx pgx.Tx) error {
		sess, err := store.LockSession(ctx, tx, id)
		if err != nil {
			return err
		}
		if err := fn(tx, sess); err != nil {
			return err
		}
		out = sess
		return store.UpdateSession(ctx, tx, sess)
	})
	return out, err
}

// Terminate kết thúc phiên; phiên đã kết thúc trả trạng thái hiện tại (idempotent).
func (s *Service) Terminate(ctx context.Context, id uuid.UUID, reason api.EndReason) (*store.Session, error) {
	sess, err := s.transition(ctx, id, func(tx pgx.Tx, sess *store.Session) error {
		if !store.Active(sess.Phase) {
			return errSkip
		}
		now := s.now()
		r := string(reason)
		sess.Phase, sess.EndReason, sess.EndedAt = store.PhaseTerminated, &r, &now
		used := 0
		if sess.ReadyAt != nil {
			used = int(math.Ceil(now.Sub(*sess.ReadyAt).Minutes()))
		}
		sess.UsedMinutes = &used
		if reason == api.EndReasonOUTOFMEMORY {
			sess.OOMKilled = true
		}
		data := map[string]any{"sessionId": id, "endReason": r, "endedAt": now, "usedMinutes": used, "labels": sess.Labels}
		if sess.ReadyAt != nil {
			data["resourceUsage"] = usage(sess)
		}
		_, err := s.events.Write(ctx, tx, events.SessionTerminated, id.String(), data)
		return err
	})
	if errors.Is(err, errSkip) {
		return s.Get(ctx, id)
	}
	if errors.Is(err, store.ErrNotFound) {
		return nil, problem.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	if cancel, ok := s.provisions[id]; ok {
		cancel()
	}
	delete(s.activity, id)
	s.mu.Unlock()
	s.notifier.CloseSession(id, string(reason))
	if sess.ContainerID != nil {
		if err := s.drv.Destroy(context.WithoutCancel(ctx), *sess.ContainerID); err != nil {
			s.log.Warn("xóa môi trường", "session", id, "err", err)
		}
	}
	return sess, nil
}

// Extend gia hạn phiên luyện tập (BR-M1.10).
func (s *Service) Extend(ctx context.Context, id uuid.UUID, minutes int) (*store.Session, error) {
	sess, err := s.transition(ctx, id, func(_ pgx.Tx, sess *store.Session) error {
		if sess.Phase != store.PhaseReady && sess.Phase != store.PhaseExpiring {
			return errNotActive
		}
		if sess.Priority != string(api.PRACTICE) || sess.ExtensionMinutes+minutes > sess.MaxExtensionMinutes {
			return problem.New(http.StatusConflict, "EXTENSION_LIMIT",
				fmt.Sprintf("Đã gia hạn %d/%d phút; phiên thi không gia hạn", sess.ExtensionMinutes, sess.MaxExtensionMinutes))
		}
		sess.ExtensionMinutes += minutes
		expires := sess.ExpiresAt.Add(time.Duration(minutes) * time.Minute)
		sess.ExpiresAt = &expires
		if expires.Sub(s.now()) > expiringBefore {
			sess.Phase, sess.ExpiringNotified = store.PhaseReady, false
		}
		return nil
	})
	if errors.Is(err, store.ErrNotFound) {
		return nil, problem.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	s.notifier.Notify(id, StatusFrame(sess))
	return sess, nil
}

// Exec chạy lệnh kiểm tra hoặc script khôi phục (BR-M1.16).
func (s *Service) Exec(ctx context.Context, id uuid.UUID, req api.ExecRequest) (driver.ExecResult, error) {
	timeout := time.Duration(req.TimeoutSeconds) * time.Second
	if req.Purpose == api.CHECK && timeout > checkMaxTimeout {
		return driver.ExecResult{}, problem.Validation("timeoutSeconds tối đa 15 với purpose CHECK")
	}
	sess, err := s.Get(ctx, id)
	if err != nil {
		return driver.ExecResult{}, err
	}
	if (sess.Phase != store.PhaseReady && sess.Phase != store.PhaseExpiring) || sess.ContainerID == nil {
		return driver.ExecResult{}, errNotActive
	}
	s.mu.Lock()
	if s.busy[id] {
		s.mu.Unlock()
		return driver.ExecResult{}, errExecBusy
	}
	s.busy[id] = true
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.busy, id)
		s.mu.Unlock()
	}()
	return s.drv.Exec(ctx, *sess.ContainerID, driver.ExecSpec{Purpose: string(req.Purpose), Script: req.Script, Timeout: timeout})
}

// IssueTicket cấp vé terminal cho phiên đang chạy (BR-M1.13); vé chỉ trả về một lần, CSDL lưu mã băm.
func (s *Service) IssueTicket(ctx context.Context, id uuid.UUID, mode string, userID uuid.UUID) (string, time.Time, error) {
	sess, err := s.Get(ctx, id)
	if err != nil {
		return "", time.Time{}, err
	}
	if sess.Phase != store.PhaseReady && sess.Phase != store.PhaseExpiring {
		return "", time.Time{}, errNotActive
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", time.Time{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	expires := s.now().Add(ticketTTL)
	err = s.store.InTx(ctx, func(tx pgx.Tx) error {
		if err := store.InsertTicket(ctx, tx, token, id, mode, userID, expires); err != nil {
			return err
		}
		if mode == string(api.VIEW) {
			return store.InsertViewAccess(ctx, tx, id, userID)
		}
		return nil
	})
	return token, expires, err
}

// Touch ghi nhận dữ liệu nhập từ terminal (BR-M1.9).
func (s *Service) Touch(id uuid.UUID) {
	s.mu.Lock()
	s.activity[id] = s.now()
	s.mu.Unlock()
}

// ToAPI chuyển phiên sang lược đồ của hợp đồng.
func ToAPI(sess *store.Session) api.Session {
	out := api.Session{
		Id: sess.ID, Phase: api.Phase(sess.Phase), Priority: api.Priority(sess.Priority), OwnerUserId: sess.OwnerUserID,
		CreatedAt: sess.CreatedAt, ReadyAt: sess.ReadyAt, ExpiresAt: sess.ExpiresAt, EndedAt: sess.EndedAt,
		FailureDetail: sess.FailureDetail, Labels: sess.Labels,
	}
	if out.Labels == nil {
		out.Labels = map[string]string{}
	}
	if sess.EndReason != nil {
		r := api.EndReason(*sess.EndReason)
		out.EndReason = &r
	}
	if sess.ReadyAt != nil {
		u := usage(sess)
		out.ResourceUsage = &u
	}
	return out
}

func usage(sess *store.Session) api.ResourceUsage {
	u := api.ResourceUsage{MemoryLimitBytes: sess.MemoryLimitBytes, OomKilled: sess.OOMKilled}
	if sess.MemoryPeakBytes != nil {
		u.MemoryPeakBytes = *sess.MemoryPeakBytes
	}
	return u
}

// StatusFrame là khung status gửi tới terminal (CR-DOC-03 mục 6).
func StatusFrame(sess *store.Session) map[string]any {
	f := map[string]any{"type": "status", "phase": sess.Phase}
	if sess.ExpiresAt != nil {
		f["expiresAt"] = sess.ExpiresAt
	}
	return f
}

// Wait chờ các goroutine dựng môi trường kết thúc; dùng khi dừng dịch vụ.
func (s *Service) Wait() { s.wg.Wait() }
