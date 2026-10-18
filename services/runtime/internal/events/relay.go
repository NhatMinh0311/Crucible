package events

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Relay chuyển sự kiện từ outbox lên exchange. Sự kiện chỉ được đánh dấu đã gửi sau khi broker xác nhận.
type Relay struct {
	Pool     *pgxpool.Pool
	URL      string
	Exchange string
	Interval time.Duration
	Batch    int
	Log      *slog.Logger

	conn *amqp.Connection
	ch   *amqp.Channel
}

// Run chuyển tiếp cho tới khi ctx kết thúc; mất kết nối thì kết nối lại ở chu kỳ sau.
func (r *Relay) Run(ctx context.Context) {
	r.defaults()
	t := time.NewTicker(r.Interval)
	defer t.Stop()
	defer r.close()
	for {
		for {
			n, err := r.RelayOnce(ctx)
			if err != nil {
				if ctx.Err() == nil {
					r.Log.Warn("chuyển tiếp outbox lỗi", "err", err)
				}
				r.close()
				break
			}
			if n < r.Batch {
				break
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

func (r *Relay) defaults() {
	if r.Interval == 0 {
		r.Interval = time.Second
	}
	if r.Batch == 0 {
		r.Batch = 100
	}
}

func (r *Relay) connect() error {
	if r.ch != nil && !r.ch.IsClosed() {
		return nil
	}
	r.close()
	conn, err := amqp.Dial(r.URL)
	if err != nil {
		return fmt.Errorf("kết nối RabbitMQ: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("mở channel: %w", err)
	}
	if err := ch.ExchangeDeclare(r.Exchange, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		_ = conn.Close()
		return fmt.Errorf("khai báo exchange: %w", err)
	}
	if err := ch.Confirm(false); err != nil {
		_ = conn.Close()
		return fmt.Errorf("bật publisher confirm: %w", err)
	}
	r.conn, r.ch = conn, ch
	return nil
}

// Close đóng kết nối tới broker.
func (r *Relay) Close() { r.close() }

func (r *Relay) close() {
	if r.conn != nil {
		_ = r.conn.Close()
	}
	r.conn, r.ch = nil, nil
}

type pending struct {
	id      uuid.UUID
	typ     string
	payload []byte
}

// RelayOnce gửi một lô sự kiện; trả số sự kiện đã gửi.
func (r *Relay) RelayOnce(ctx context.Context) (int, error) {
	r.defaults()
	if err := r.connect(); err != nil {
		return 0, err
	}
	sent := 0
	err := pgx.BeginFunc(ctx, r.Pool, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `SELECT id, type, payload::text FROM crucible_outbox
			WHERE sent_at IS NULL ORDER BY created_at LIMIT $1 FOR UPDATE SKIP LOCKED`, r.Batch)
		if err != nil {
			return err
		}
		batch, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (pending, error) {
			var p pending
			var payload string
			err := row.Scan(&p.id, &p.typ, &payload)
			p.payload = []byte(payload)
			return p, err
		})
		if err != nil || len(batch) == 0 {
			return err
		}
		confirms := make([]*amqp.DeferredConfirmation, 0, len(batch))
		ids := make([]uuid.UUID, 0, len(batch))
		for _, p := range batch {
			c, err := r.ch.PublishWithDeferredConfirmWithContext(ctx, r.Exchange, p.typ, true, false, amqp.Publishing{
				ContentType:     "application/cloudevents+json",
				ContentEncoding: "utf-8",
				MessageId:       p.id.String(),
				Type:            p.typ,
				DeliveryMode:    amqp.Persistent,
				Body:            p.payload,
			})
			if err != nil {
				return fmt.Errorf("gửi sự kiện %s: %w", p.id, err)
			}
			confirms = append(confirms, c)
			ids = append(ids, p.id)
		}
		waitCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		for _, c := range confirms {
			ok, err := c.WaitContext(waitCtx)
			if err != nil {
				return fmt.Errorf("chờ xác nhận của broker: %w", err)
			}
			if !ok {
				return errors.New("broker từ chối sự kiện")
			}
		}
		if _, err := tx.Exec(ctx, `UPDATE crucible_outbox SET sent_at = now() WHERE id = ANY($1)`, ids); err != nil {
			return err
		}
		sent = len(batch)
		return nil
	})
	return sent, err
}
