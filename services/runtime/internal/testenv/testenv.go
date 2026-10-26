// Package testenv dựng PostgreSQL và RabbitMQ bằng Testcontainers cho kiểm thử tích hợp.
package testenv

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcrabbit "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Cùng phiên bản với dev/compose.yaml.
const (
	PostgresImage = "postgres:18.6-alpine"
	RabbitImage   = "rabbitmq:4.3.6-management-alpine"
)

// Postgres khởi động PostgreSQL và trả chuỗi kết nối; container bị xóa khi kiểm thử kết thúc.
func Postgres(t testing.TB) string {
	t.Helper()
	ctx := context.Background()
	c, err := tcpostgres.Run(ctx, PostgresImage,
		tcpostgres.WithDatabase("runtime"), tcpostgres.WithUsername("runtime"), tcpostgres.WithPassword("runtime"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(60*time.Second)))
	if err != nil {
		t.Fatalf("khởi động PostgreSQL: %v", err)
	}
	t.Cleanup(func() { _ = c.Terminate(context.Background()) })
	url, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("chuỗi kết nối PostgreSQL: %v", err)
	}
	return url
}

// RabbitMQ khởi động RabbitMQ và trả URL AMQP.
func RabbitMQ(t testing.TB) string {
	t.Helper()
	ctx := context.Background()
	c, err := tcrabbit.Run(ctx, RabbitImage, tcrabbit.WithAdminUsername("crucible"), tcrabbit.WithAdminPassword("crucible"))
	if err != nil {
		t.Fatalf("khởi động RabbitMQ: %v", err)
	}
	t.Cleanup(func() { _ = c.Terminate(context.Background()) })
	url, err := c.AmqpURL(ctx)
	if err != nil {
		t.Fatalf("URL RabbitMQ: %v", err)
	}
	return fmt.Sprint(url)
}
