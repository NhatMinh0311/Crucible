// Package store lưu trạng thái của runtime-api trong PostgreSQL (CSDL `runtime`, CR-MOD-M1 mục 7).
package store

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Lỗi của tầng lưu trữ.
var (
	ErrNotFound            = errors.New("không tìm thấy")
	ErrActiveSessionExists = errors.New("học viên đã có phiên cùng loại chưa kết thúc")
	ErrTicketInvalid       = errors.New("vé không hợp lệ, đã dùng hoặc hết hạn")
)

// Phase là pha của phiên (CR-DOC-02 mục 7.1).
const (
	PhasePending      = "PENDING"
	PhaseProvisioning = "PROVISIONING"
	PhaseReady        = "READY"
	PhaseExpiring     = "EXPIRING"
	PhaseTerminated   = "TERMINATED"
	PhaseFailed       = "FAILED"
)

// Active cho biết pha chưa kết thúc.
func Active(phase string) bool {
	switch phase {
	case PhasePending, PhaseProvisioning, PhaseReady, PhaseExpiring:
		return true
	}
	return false
}

// Session là một dòng của bảng sessions.
type Session struct {
	ID                  uuid.UUID
	OwnerUserID         uuid.UUID
	Priority            string
	Profile             string
	Image               string
	QuotaAccountRef     string
	ReservationRef      *string
	RecordInput         bool
	Labels              map[string]string
	Environment         json.RawMessage
	Phase               string
	EndReason           *string
	FailureDetail       *string
	CreatedAt           time.Time
	ReadyAt             *time.Time
	ExpiresAt           *time.Time
	EndedAt             *time.Time
	TTLMinutes          int
	IdleMinutes         int
	MaxExtensionMinutes int
	ExtensionMinutes    int
	LastActivityAt      *time.Time
	IdleWarnedAt        *time.Time
	ExpiringNotified    bool
	UsedMinutes         *int
	ContainerID         *string
	MemoryPeakBytes     *int64
	MemoryLimitBytes    int64
	OOMKilled           bool
}

// Store bọc pool kết nối PostgreSQL.
type Store struct {
	Pool *pgxpool.Pool
}

// Open kết nối CSDL và chạy migration.
func Open(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("kết nối CSDL: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("kết nối CSDL: %w", err)
	}
	if err := migrateUp(pool); err != nil {
		pool.Close()
		return nil, err
	}
	return &Store{Pool: pool}, nil
}

// Close đóng pool.
func (s *Store) Close() { s.Pool.Close() }

func migrateUp(pool *pgxpool.Pool) error {
	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("migration: %w", err)
	}
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	drv, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		return fmt.Errorf("migration: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", src, "pgx5", drv)
	if err != nil {
		return fmt.Errorf("migration: %w", err)
	}
	// Close trả kết nối mà driver migration giữ để khóa; không gọi thì pool không đóng được.
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration: %w", err)
	}
	return nil
}

// InTx chạy fn trong một transaction; lỗi hoặc panic làm rollback.
func (s *Store) InTx(ctx context.Context, fn func(pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, s.Pool, fn)
}

const sessionColumns = `id, owner_user_id, priority, profile, image, quota_account_ref, reservation_ref, record_input,
	labels, environment, phase, end_reason, failure_detail, created_at, ready_at, expires_at, ended_at, ttl_minutes,
	idle_minutes, max_extension_minutes, extension_minutes, last_activity_at, idle_warned_at, expiring_notified,
	used_minutes, container_id, memory_peak_bytes, memory_limit_bytes, oom_killed`

func scanSession(row pgx.Row) (*Session, error) {
	var s Session
	var labels []byte
	err := row.Scan(&s.ID, &s.OwnerUserID, &s.Priority, &s.Profile, &s.Image, &s.QuotaAccountRef, &s.ReservationRef,
		&s.RecordInput, &labels, &s.Environment, &s.Phase, &s.EndReason, &s.FailureDetail, &s.CreatedAt, &s.ReadyAt,
		&s.ExpiresAt, &s.EndedAt, &s.TTLMinutes, &s.IdleMinutes, &s.MaxExtensionMinutes, &s.ExtensionMinutes,
		&s.LastActivityAt, &s.IdleWarnedAt, &s.ExpiringNotified, &s.UsedMinutes, &s.ContainerID, &s.MemoryPeakBytes,
		&s.MemoryLimitBytes, &s.OOMKilled)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(labels, &s.Labels); err != nil {
		return nil, fmt.Errorf("đọc labels: %w", err)
	}
	return &s, nil
}

// Querier là pool hoặc transaction.
type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// InsertSession thêm phiên mới; trả ErrActiveSessionExists khi vi phạm BR-M1.1.
func InsertSession(ctx context.Context, q Querier, s *Session) error {
	labels, err := json.Marshal(s.Labels)
	if err != nil {
		return err
	}
	_, err = q.Exec(ctx, `INSERT INTO sessions (id, owner_user_id, priority, profile, image, quota_account_ref,
		reservation_ref, record_input, labels, environment, phase, created_at, ttl_minutes, idle_minutes,
		max_extension_minutes, memory_limit_bytes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		s.ID, s.OwnerUserID, s.Priority, s.Profile, s.Image, s.QuotaAccountRef, s.ReservationRef, s.RecordInput,
		labels, s.Environment, s.Phase, s.CreatedAt, s.TTLMinutes, s.IdleMinutes, s.MaxExtensionMinutes,
		s.MemoryLimitBytes)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "sessions_one_active" {
		return ErrActiveSessionExists
	}
	return err
}

// GetSession đọc một phiên.
func GetSession(ctx context.Context, q Querier, id uuid.UUID) (*Session, error) {
	return scanSession(q.QueryRow(ctx, `SELECT `+sessionColumns+` FROM sessions WHERE id = $1`, id))
}

// LockSession đọc một phiên và khóa dòng tới hết transaction.
func LockSession(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Session, error) {
	return scanSession(tx.QueryRow(ctx, `SELECT `+sessionColumns+` FROM sessions WHERE id = $1 FOR UPDATE`, id))
}

// ListActiveSessions trả các phiên chưa kết thúc.
func ListActiveSessions(ctx context.Context, q Querier) ([]*Session, error) {
	rows, err := q.Query(ctx, `SELECT `+sessionColumns+` FROM sessions
		WHERE phase IN ('PENDING', 'PROVISIONING', 'READY', 'EXPIRING') ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Session
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// UpdateSession ghi các trường thay đổi được của phiên.
func UpdateSession(ctx context.Context, q Querier, s *Session) error {
	tag, err := q.Exec(ctx, `UPDATE sessions SET phase = $2, end_reason = $3, failure_detail = $4, ready_at = $5,
		expires_at = $6, ended_at = $7, extension_minutes = $8, last_activity_at = $9, idle_warned_at = $10,
		expiring_notified = $11, used_minutes = $12, container_id = $13, memory_peak_bytes = $14, oom_killed = $15
		WHERE id = $1`,
		s.ID, s.Phase, s.EndReason, s.FailureDetail, s.ReadyAt, s.ExpiresAt, s.EndedAt, s.ExtensionMinutes,
		s.LastActivityAt, s.IdleWarnedAt, s.ExpiringNotified, s.UsedMinutes, s.ContainerID, s.MemoryPeakBytes,
		s.OOMKilled)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

// HashToken trả mã băm SHA-256 của vé; CSDL chỉ lưu mã băm (BR-M1.13).
func HashToken(token string) []byte {
	h := sha256.Sum256([]byte(token))
	return h[:]
}

// InsertTicket lưu vé terminal.
func InsertTicket(ctx context.Context, q Querier, token string, sessionID uuid.UUID, mode string, userID uuid.UUID, expiresAt time.Time) error {
	_, err := q.Exec(ctx, `INSERT INTO terminal_tickets (token_hash, session_id, mode, user_id, expires_at)
		VALUES ($1, $2, $3, $4, $5)`, HashToken(token), sessionID, mode, userID, expiresAt)
	return err
}

// Ticket là vé đã được tiêu thụ.
type Ticket struct {
	SessionID uuid.UUID
	Mode      string
	UserID    uuid.UUID
}

// ConsumeTicket đánh dấu vé đã dùng nếu còn hạn và chưa dùng; vé chỉ dùng được một lần.
func ConsumeTicket(ctx context.Context, q Querier, token string) (Ticket, error) {
	var t Ticket
	err := q.QueryRow(ctx, `UPDATE terminal_tickets SET used_at = now()
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now()
		RETURNING session_id, mode, user_id`, HashToken(token)).Scan(&t.SessionID, &t.Mode, &t.UserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return t, ErrTicketInvalid
	}
	return t, err
}

// InsertViewAccess ghi nhật ký cấp vé xem cho giảng viên.
func InsertViewAccess(ctx context.Context, q Querier, sessionID, viewerUserID uuid.UUID) error {
	_, err := q.Exec(ctx, `INSERT INTO view_access_log (session_id, viewer_user_id) VALUES ($1, $2)`, sessionID, viewerUserID)
	return err
}

// Idempotent là phản hồi đã lưu của một khóa lặp lại an toàn.
type Idempotent struct {
	RequestHash []byte
	Status      *int
	Response    json.RawMessage
}

// ClaimIdempotencyKey ghi nhận khóa trong transaction hiện tại. Trả nil khi khóa mới; trả bản ghi đã có khi khóa
// đã được dùng.
func ClaimIdempotencyKey(ctx context.Context, tx pgx.Tx, caller string, key uuid.UUID, operation string, requestHash []byte) (*Idempotent, error) {
	tag, err := tx.Exec(ctx, `INSERT INTO idempotency_keys (caller, key, operation, request_hash)
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`, caller, key, operation, requestHash)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 1 {
		return nil, nil
	}
	var r Idempotent
	err = tx.QueryRow(ctx, `SELECT request_hash, status, response FROM idempotency_keys
		WHERE caller = $1 AND key = $2 AND operation = $3`, caller, key, operation).Scan(&r.RequestHash, &r.Status, &r.Response)
	return &r, err
}

// SaveIdempotentResponse lưu phản hồi cho khóa đã ghi nhận.
func SaveIdempotentResponse(ctx context.Context, q Querier, caller string, key uuid.UUID, operation string, status int, response any) error {
	body, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = q.Exec(ctx, `UPDATE idempotency_keys SET status = $4, response = $5
		WHERE caller = $1 AND key = $2 AND operation = $3`, caller, key, operation, status, body)
	return err
}

// Cleanup xóa vé, khóa lặp lại an toàn và outbox đã gửi quá hạn lưu giữ (CR-MOD-M1 mục 10).
func (s *Store) Cleanup(ctx context.Context) error {
	for _, q := range []string{
		`DELETE FROM terminal_tickets WHERE expires_at < now() - interval '1 hour'`,
		`DELETE FROM idempotency_keys WHERE created_at < now() - interval '24 hours'`,
		`DELETE FROM crucible_outbox WHERE sent_at < now() - interval '7 days'`,
	} {
		if _, err := s.Pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("dọn bảng: %w", err)
		}
	}
	return nil
}
