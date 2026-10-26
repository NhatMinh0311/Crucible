package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/api"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/events"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/store"
)

// Run nhận lại môi trường sau khởi động rồi chạy các tác vụ định kỳ cho tới khi ctx kết thúc (CR-MOD-M1 mục 10).
func (s *Service) Run(ctx context.Context) error {
	if err := s.Reconcile(ctx); err != nil {
		return err
	}
	exits := s.drv.Watch(ctx)
	timers := time.NewTicker(5 * time.Second)
	defer timers.Stop()
	sample := time.NewTicker(sampleInterval)
	defer sample.Stop()
	cleanup := time.NewTicker(time.Hour)
	defer cleanup.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case e, ok := <-exits:
			if !ok {
				exits = nil
				continue
			}
			reason := api.EndReasonENVIRONMENTFAILED
			if e.OOM {
				reason = api.EndReasonOUTOFMEMORY
			}
			// BR-M1.17: môi trường dừng khi phiên chưa được yêu cầu kết thúc. Phiên đã kết thúc thì Terminate bỏ qua.
			if sess, err := s.Get(ctx, e.SessionID); err == nil && (sess.Phase == store.PhaseReady || sess.Phase == store.PhaseExpiring) {
				s.log.Info("môi trường dừng ngoài ý muốn", "session", e.SessionID, "oom", e.OOM, "exitCode", e.ExitCode)
				if _, err := s.Terminate(ctx, e.SessionID, reason); err != nil {
					s.log.Error("kết thúc phiên", "session", e.SessionID, "err", err)
				}
			}
		case <-timers.C:
			s.Tick(ctx)
		case <-sample.C:
			s.sampleMemory(ctx)
		case <-cleanup.C:
			if err := s.store.Cleanup(ctx); err != nil {
				s.log.Warn("dọn bảng", "err", err)
			}
		}
	}
}

// Reconcile đối chiếu phiên chưa kết thúc trong CSDL với môi trường đang có (D11, BR-M1.12).
func (s *Service) Reconcile(ctx context.Context) error {
	instances, err := s.drv.List(ctx)
	if err != nil {
		return err
	}
	active, err := store.ListActiveSessions(ctx, s.store.Pool)
	if err != nil {
		return err
	}
	keep := map[string]bool{}
	for _, sess := range active {
		switch sess.Phase {
		case store.PhasePending, store.PhaseProvisioning:
			// Dựng môi trường bị gián đoạn khi dịch vụ dừng.
			_, err := s.transition(ctx, sess.ID, func(tx pgx.Tx, sess *store.Session) error {
				reason, detail := string(api.EndReasonPROVISIONFAILED), "runtime-api khởi động lại khi đang dựng môi trường"
				now := s.now()
				sess.Phase, sess.EndReason, sess.FailureDetail, sess.EndedAt = store.PhaseFailed, &reason, &detail, &now
				_, err := s.events.Write(ctx, tx, events.SessionFailed, sess.ID.String(), map[string]any{
					"sessionId": sess.ID, "endReason": reason, "detail": detail, "labels": sess.Labels,
				})
				return err
			})
			if err != nil {
				return err
			}
		default:
			running := false
			for _, in := range instances {
				if in.SessionID == sess.ID && in.Running && sess.ContainerID != nil && in.Handle == *sess.ContainerID {
					running = true
				}
			}
			if running {
				keep[*sess.ContainerID] = true
				continue
			}
			if _, err := s.Terminate(ctx, sess.ID, api.EndReasonADMIN); err != nil {
				return err
			}
		}
	}
	for _, in := range instances {
		if !keep[in.Handle] {
			if err := s.drv.Destroy(ctx, in.Handle); err != nil {
				s.log.Warn("xóa môi trường mồ côi", "handle", in.Handle, "err", err)
			}
		}
	}
	return nil
}

// Tick xét hết hạn, sắp hết hạn và không hoạt động của các phiên đang chạy (BR-M1.9, D7).
func (s *Service) Tick(ctx context.Context) {
	active, err := store.ListActiveSessions(ctx, s.store.Pool)
	if err != nil {
		s.log.Warn("đọc phiên đang chạy", "err", err)
		return
	}
	now := s.now()
	for _, sess := range active {
		if sess.ExpiresAt == nil || (sess.Phase != store.PhaseReady && sess.Phase != store.PhaseExpiring) {
			continue
		}
		if !now.Before(*sess.ExpiresAt) {
			s.terminate(ctx, sess.ID, api.EndReasonEXPIRED)
			continue
		}
		if !sess.ExpiringNotified && sess.ExpiresAt.Sub(now) <= expiringBefore {
			s.markExpiring(ctx, sess.ID)
		}
		if sess.Priority != string(api.PRACTICE) {
			continue
		}
		last := s.lastActivity(sess)
		idleAt := last.Add(time.Duration(sess.IdleMinutes) * time.Minute)
		switch {
		case !now.Before(idleAt):
			s.terminate(ctx, sess.ID, api.EndReasonIDLE)
		case idleAt.Sub(now) <= idleWarnBefore && (sess.IdleWarnedAt == nil || sess.IdleWarnedAt.Before(last)):
			s.updateActivity(ctx, sess.ID, last, &now)
			s.notifier.Notify(sess.ID, map[string]any{"type": "warning", "code": "IDLE", "terminateAt": idleAt})
		case sess.LastActivityAt == nil || last.Sub(*sess.LastActivityAt) >= activityFlushGap:
			s.updateActivity(ctx, sess.ID, last, nil)
		}
	}
}

func (s *Service) lastActivity(sess *store.Session) time.Time {
	last := *sess.ReadyAt
	if sess.LastActivityAt != nil && sess.LastActivityAt.After(last) {
		last = *sess.LastActivityAt
	}
	s.mu.Lock()
	if t, ok := s.activity[sess.ID]; ok && t.After(last) {
		last = t
	}
	s.mu.Unlock()
	return last
}

func (s *Service) updateActivity(ctx context.Context, id uuid.UUID, last time.Time, warned *time.Time) {
	_, err := s.transition(ctx, id, func(_ pgx.Tx, sess *store.Session) error {
		sess.LastActivityAt = &last
		if warned != nil {
			sess.IdleWarnedAt = warned
		}
		return nil
	})
	if err != nil {
		s.log.Warn("ghi thời điểm hoạt động", "session", id, "err", err)
	}
}

func (s *Service) markExpiring(ctx context.Context, id uuid.UUID) {
	sess, err := s.transition(ctx, id, func(tx pgx.Tx, sess *store.Session) error {
		if sess.Phase != store.PhaseReady && sess.Phase != store.PhaseExpiring || sess.ExpiringNotified {
			return errSkip
		}
		sess.Phase, sess.ExpiringNotified = store.PhaseExpiring, true
		_, err := s.events.Write(ctx, tx, events.SessionExpiring, id.String(), map[string]any{
			"sessionId": id, "expiresAt": sess.ExpiresAt, "labels": sess.Labels,
		})
		return err
	})
	if err != nil {
		if !errors.Is(err, errSkip) {
			s.log.Warn("ghi EXPIRING", "session", id, "err", err)
		}
		return
	}
	s.notifier.Notify(id, map[string]any{"type": "expiring", "expiresAt": sess.ExpiresAt})
}

func (s *Service) terminate(ctx context.Context, id uuid.UUID, reason api.EndReason) {
	if _, err := s.Terminate(ctx, id, reason); err != nil {
		s.log.Error("kết thúc phiên", "session", id, "reason", reason, "err", err)
	}
}

// sampleMemory ghi bộ nhớ cao nhất của phiên đang chạy (CR-MOD-M1 mục 5.5).
func (s *Service) sampleMemory(ctx context.Context) {
	active, err := store.ListActiveSessions(ctx, s.store.Pool)
	if err != nil {
		return
	}
	for _, sess := range active {
		if sess.ContainerID == nil || (sess.Phase != store.PhaseReady && sess.Phase != store.PhaseExpiring) {
			continue
		}
		used, err := s.drv.MemoryUsage(ctx, *sess.ContainerID)
		if err != nil || (sess.MemoryPeakBytes != nil && used <= *sess.MemoryPeakBytes) {
			continue
		}
		_, err = s.transition(ctx, sess.ID, func(_ pgx.Tx, sess *store.Session) error {
			if sess.MemoryPeakBytes != nil && used <= *sess.MemoryPeakBytes {
				return errSkip
			}
			sess.MemoryPeakBytes = &used
			return nil
		})
		if err != nil && !errors.Is(err, errSkip) {
			s.log.Warn("ghi bộ nhớ cao nhất", "session", sess.ID, "err", err)
		}
	}
}
