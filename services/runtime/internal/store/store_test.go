package store_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/store"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/testenv"
)

func newSession(owner uuid.UUID, priority string) *store.Session {
	return &store.Session{
		ID: uuid.New(), OwnerUserID: owner, Priority: priority, Profile: "small", Image: "img",
		QuotaAccountRef: "class:1", Labels: map[string]string{"attemptId": "a1"}, Environment: []byte(`{}`),
		Phase: store.PhasePending, CreatedAt: time.Now().UTC(), TTLMinutes: 60, IdleMinutes: 15,
		MemoryLimitBytes: 256 << 20,
	}
}

func TestStore(t *testing.T) {
	ctx := context.Background()
	s, err := store.Open(ctx, testenv.Postgres(t))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	t.Run("BR-M1.1 một phiên PRACTICE chưa kết thúc", func(t *testing.T) {
		owner := uuid.New()
		first := newSession(owner, "PRACTICE")
		if err := store.InsertSession(ctx, s.Pool, first); err != nil {
			t.Fatal(err)
		}
		if err := store.InsertSession(ctx, s.Pool, newSession(owner, "PRACTICE")); !errors.Is(err, store.ErrActiveSessionExists) {
			t.Fatalf("mong đợi ErrActiveSessionExists, nhận %v", err)
		}
		if err := store.InsertSession(ctx, s.Pool, newSession(owner, "EXAM")); err != nil {
			t.Fatalf("phiên EXAM phải được tạo: %v", err)
		}
		first.Phase = store.PhaseTerminated
		if err := store.UpdateSession(ctx, s.Pool, first); err != nil {
			t.Fatal(err)
		}
		if err := store.InsertSession(ctx, s.Pool, newSession(owner, "PRACTICE")); err != nil {
			t.Fatalf("sau khi phiên cũ kết thúc phải tạo được phiên mới: %v", err)
		}
		got, err := store.GetSession(ctx, s.Pool, first.ID)
		if err != nil || got.Labels["attemptId"] != "a1" || got.Phase != store.PhaseTerminated {
			t.Fatalf("đọc lại phiên sai: %+v %v", got, err)
		}
	})

	t.Run("BR-M1.13 vé dùng một lần, hết hạn", func(t *testing.T) {
		sess := newSession(uuid.New(), "PRACTICE")
		if err := store.InsertSession(ctx, s.Pool, sess); err != nil {
			t.Fatal(err)
		}
		_ = store.InsertTicket(ctx, s.Pool, "ve-1", sess.ID, "INTERACTIVE", sess.OwnerUserID, time.Now().Add(30*time.Second))
		_ = store.InsertTicket(ctx, s.Pool, "ve-het-han", sess.ID, "INTERACTIVE", sess.OwnerUserID, time.Now().Add(-time.Second))
		if tk, err := store.ConsumeTicket(ctx, s.Pool, "ve-1"); err != nil || tk.SessionID != sess.ID {
			t.Fatalf("lần đầu phải dùng được: %v", err)
		}
		if _, err := store.ConsumeTicket(ctx, s.Pool, "ve-1"); !errors.Is(err, store.ErrTicketInvalid) {
			t.Fatalf("dùng lại phải bị từ chối, nhận %v", err)
		}
		if _, err := store.ConsumeTicket(ctx, s.Pool, "ve-het-han"); !errors.Is(err, store.ErrTicketInvalid) {
			t.Fatalf("vé hết hạn phải bị từ chối, nhận %v", err)
		}
	})

	t.Run("khóa lặp lại an toàn", func(t *testing.T) {
		key := uuid.New()
		err := s.InTx(ctx, func(tx pgx.Tx) error {
			prev, err := store.ClaimIdempotencyKey(ctx, tx, "u1", key, "createSession", []byte("h1"))
			if err != nil || prev != nil {
				t.Fatalf("khóa mới: %v %v", prev, err)
			}
			return store.SaveIdempotentResponse(ctx, tx, "u1", key, "createSession", 202, map[string]string{"id": "x"})
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = s.InTx(ctx, func(tx pgx.Tx) error {
			prev, err := store.ClaimIdempotencyKey(ctx, tx, "u1", key, "createSession", []byte("h1"))
			if err != nil || prev == nil || *prev.Status != 202 || string(prev.RequestHash) != "h1" {
				t.Fatalf("khóa đã dùng phải trả phản hồi cũ: %+v %v", prev, err)
			}
			return nil
		})
	})
}
