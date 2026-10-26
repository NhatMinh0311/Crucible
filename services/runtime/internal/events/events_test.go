package events_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/events"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/store"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/testenv"
)

func TestOutboxRelay(t *testing.T) {
	ctx := context.Background()
	s, err := store.Open(ctx, testenv.Postgres(t))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	amqpURL := testenv.RabbitMQ(t)

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare("crucible.events", amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		t.Fatal(err)
	}
	q, _ := ch.QueueDeclare("", false, true, true, false, nil)
	_ = ch.QueueBind(q.Name, "runtime.session.*", "crucible.events", false, nil)
	msgs, _ := ch.Consume(q.Name, "", true, true, false, false, nil)

	w := events.Writer{Source: "/crucible/runtime-api"}
	err = s.InTx(ctx, func(tx pgx.Tx) error {
		for i := 0; i < 2; i++ {
			if _, err := w.Write(ctx, tx, events.SessionReady, "s-1", map[string]any{"sessionId": "s-1", "i": i}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	r := &events.Relay{Pool: s.Pool, URL: amqpURL, Exchange: "crucible.events", Log: slog.Default()}
	n, err := r.RelayOnce(ctx)
	if err != nil || n != 2 {
		t.Fatalf("mong đợi gửi 2 sự kiện, gửi %d: %v", n, err)
	}
	defer r.Close()
	if n, _ := r.RelayOnce(ctx); n != 0 {
		t.Fatalf("sự kiện đã gửi không được gửi lại, gửi %d", n)
	}

	for want := int64(1); want <= 2; want++ {
		select {
		case m := <-msgs:
			if m.ContentType != "application/cloudevents+json" || m.Type != events.SessionReady || m.MessageId == "" ||
				m.RoutingKey != events.SessionReady || m.DeliveryMode != amqp.Persistent {
				t.Fatalf("thuộc tính message sai: %+v", m)
			}
			var env events.Envelope
			if err := json.Unmarshal(m.Body, &env); err != nil {
				t.Fatal(err)
			}
			if env.SpecVersion != "1.0" || env.Subject != "s-1" || env.Sequence != want || env.Source != "/crucible/runtime-api" {
				t.Fatalf("envelope sai: %+v", env)
			}
		case <-time.After(10 * time.Second):
			t.Fatal("không nhận được sự kiện")
		}
	}
}
