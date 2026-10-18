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
	"syscall"
	"time"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/config"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/httpapi"
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

	handler, err := httpapi.NewHandler(httpapi.Unimplemented{}, httpapi.Options{Log: log})
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	errc := make(chan error, 1)
	go func() {
		log.Info("runtime-api lắng nghe", "port", cfg.Port, "driver", cfg.Driver)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
