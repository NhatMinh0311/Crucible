// Chương trình runtime-api của module M1 (CR-MOD-M1 mục 2).
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/auth"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/config"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/driver/docker"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/events"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/httpapi"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/session"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/store"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/terminal"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(log)
	if err := run(log); err != nil {
		log.Error("runtime-api dừng vì lỗi", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("cấu hình: %w", err)
	}
	if cfg.Driver != "docker" {
		return errors.New("chế độ kubernetes chưa được hiện thực; chạy với RUNTIME_DRIVER=docker")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.Open(ctx, cfg.DBURL)
	if err != nil {
		return err
	}
	defer st.Close()
	drv, err := docker.New(ctx, cfg.ToolsImage, log)
	if err != nil {
		return err
	}
	defer drv.Close()

	sessions := session.New(st, drv, log, session.Options{
		SetupTimeout: cfg.SetupTimeout, ProvisionTimeout: cfg.ProvisionTimeout,
		AllowedCapabilities: cfg.AllowedCapabilities, EventsSource: cfg.EventsSource,
	})
	term := &terminal.Handler{Store: st, Sessions: sessions, Driver: drv, Events: events.Writer{Source: cfg.EventsSource},
		OriginPatterns: cfg.AllowedOrigins, Log: log}
	sessions.SetNotifier(term)

	handler, err := httpapi.NewHandler(&httpapi.Server{Sessions: sessions, TerminalURL: cfg.TerminalPublicURL}, httpapi.Options{
		Authenticate: auth.Middleware(auth.NewVerifier(ctx, cfg.OIDCIssuer, cfg.OIDCAudience)),
		Terminal:     term,
		Log:          log,
	})
	if err != nil {
		return err
	}

	bg, cancelBg := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	relay := &events.Relay{Pool: st.Pool, URL: cfg.AMQPURL, Exchange: cfg.EventsExchange, Log: log}
	wg.Go(func() { relay.Run(bg) })
	runErr := make(chan error, 1)
	wg.Go(func() { runErr <- sessions.Run(bg) })

	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	srvErr := make(chan error, 1)
	go func() {
		log.Info("runtime-api lắng nghe", "port", cfg.Port, "driver", cfg.Driver)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
	case err = <-runErr:
	case <-ctx.Done():
	}
	log.Info("đang dừng runtime-api")
	term.Shutdown()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	cancelBg()
	wg.Wait()
	sessions.Wait()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
