// Package events ghi sự kiện vào outbox trong transaction nghiệp vụ và chuyển tiếp lên RabbitMQ (ADR-0005,
// CR-DOC-03 mục 5). Định dạng và thuộc tính message giống thư viện dùng chung của M2, M3.
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Kiểu sự kiện của M1 (CR-DOC-03 mục 5.1).
const (
	SessionReady      = "runtime.session.ready"
	SessionExpiring   = "runtime.session.expiring"
	SessionTerminated = "runtime.session.terminated"
	SessionFailed     = "runtime.session.failed"
	TerminalConnected = "runtime.terminal.connected"
)

// Envelope là sự kiện theo CloudEvents 1.0 dạng JSON có cấu trúc, kèm thuộc tính mở rộng sequence.
type Envelope struct {
	SpecVersion     string    `json:"specversion"`
	ID              uuid.UUID `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	Time            time.Time `json:"time"`
	Subject         string    `json:"subject"`
	Sequence        int64     `json:"sequence"`
	DataContentType string    `json:"datacontenttype"`
	Data            any       `json:"data"`
}

// Writer ghi sự kiện vào outbox.
type Writer struct {
	Source string
}

// Write ghi một sự kiện trong transaction tx; sequence tăng dần theo subject.
func (w Writer) Write(ctx context.Context, tx pgx.Tx, typ, subject string, data any) (uuid.UUID, error) {
	var seq int64
	err := tx.QueryRow(ctx, `INSERT INTO crucible_event_sequence (subject, last_sequence) VALUES ($1, 1)
		ON CONFLICT (subject) DO UPDATE SET last_sequence = crucible_event_sequence.last_sequence + 1
		RETURNING last_sequence`, subject).Scan(&seq)
	if err != nil {
		return uuid.Nil, fmt.Errorf("số thứ tự sự kiện: %w", err)
	}
	env := Envelope{
		SpecVersion:     "1.0",
		ID:              uuid.New(),
		Source:          w.Source,
		Type:            typ,
		Time:            time.Now().UTC(),
		Subject:         subject,
		Sequence:        seq,
		DataContentType: "application/json",
		Data:            data,
	}
	payload, err := json.Marshal(env)
	if err != nil {
		return uuid.Nil, fmt.Errorf("mã hóa sự kiện: %w", err)
	}
	_, err = tx.Exec(ctx, `INSERT INTO crucible_outbox (id, type, subject, sequence, payload) VALUES ($1, $2, $3, $4, $5)`,
		env.ID, typ, subject, seq, payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ghi outbox: %w", err)
	}
	return env.ID, nil
}
