// Package config đọc cấu hình triển khai từ biến môi trường (CR-MOD-M1 mục 12).
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config là cấu hình của runtime-api. Giá trị mặc định trỏ tới dịch vụ phụ thuộc của dev/compose.yaml.
type Config struct {
	Driver              string
	Port                int
	DBURL               string
	AMQPURL             string
	EventsExchange      string
	EventsSource        string
	OIDCIssuer          string
	OIDCAudience        string
	TerminalPublicURL   string
	SetupTimeout        time.Duration
	ProvisionTimeout    time.Duration
	AllowedCapabilities []string
	ToolsImage          string
}

// Load đọc cấu hình; trả lỗi khi giá trị không hợp lệ.
func Load() (Config, error) {
	c := Config{
		Driver:            env("RUNTIME_DRIVER", "kubernetes"),
		DBURL:             env("DB_URL", "postgres://runtime:runtime-dev@localhost:5432/runtime?sslmode=disable"),
		AMQPURL:           env("AMQP_URL", "amqp://crucible:rabbitmq-dev@localhost:5672/"),
		EventsExchange:    env("EVENTS_EXCHANGE", "crucible.events"),
		EventsSource:      env("EVENTS_SOURCE", "/crucible/runtime-api"),
		OIDCIssuer:        env("OIDC_ISSUER", "http://localhost:8180/realms/crucible"),
		OIDCAudience:      env("OIDC_AUDIENCE", "runtime-api"),
		TerminalPublicURL: env("TERMINAL_PUBLIC_URL", "ws://localhost:8081/v1/terminal"),
		ToolsImage:        env("TOOLS_IMAGE", "ghcr.io/crucible-project/crucible/tools:1.0.0"),
	}
	var err error
	if c.Port, err = strconv.Atoi(env("PORT", "8081")); err != nil {
		return c, fmt.Errorf("PORT: %w", err)
	}
	if c.SetupTimeout, err = time.ParseDuration(env("SETUP_TIMEOUT", "60s")); err != nil {
		return c, fmt.Errorf("SETUP_TIMEOUT: %w", err)
	}
	if c.ProvisionTimeout, err = time.ParseDuration(env("PROVISION_TIMEOUT", "5m")); err != nil {
		return c, fmt.Errorf("PROVISION_TIMEOUT: %w", err)
	}
	// Mặc định không cấp capability: chế độ docker chạy runc, không có gVisor (BR-M1.26, ADR-0016).
	defaultCaps := ""
	if c.Driver == "kubernetes" {
		defaultCaps = "NET_ADMIN"
	}
	for _, v := range strings.Split(env("SESSION_ALLOWED_CAPABILITIES", defaultCaps), ",") {
		if v = strings.TrimSpace(v); v != "" {
			c.AllowedCapabilities = append(c.AllowedCapabilities, v)
		}
	}
	if c.Driver != "docker" && c.Driver != "kubernetes" {
		return c, fmt.Errorf("RUNTIME_DRIVER phải là docker hoặc kubernetes, nhận %q", c.Driver)
	}
	return c, nil
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
