// Package httpapi nối API sinh từ hợp đồng với tầng nghiệp vụ: định tuyến, kiểm tra yêu cầu theo lược đồ, xác thực,
// chuyển lỗi sang RFC 9457.
package httpapi

//go:generate python3 -I gen_unimplemented.py

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/api"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/problem"
)

// Options cấu hình bộ định tuyến.
type Options struct {
	// Authenticate bọc các operation của API; xác thực token và gắn người gọi vào context.
	Authenticate func(http.Handler) http.Handler
	// Terminal phục vụ WebSocket tại /v1/terminal; xác thực bằng vé nên đứng ngoài Authenticate.
	Terminal http.Handler
	Log      *slog.Logger
}

// NewHandler trả http.Handler của runtime-api.
func NewHandler(server api.StrictServerInterface, opts Options) (http.Handler, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("nạp hợp đồng: %w", err)
	}
	validate := middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		DoNotValidateServers: true,
		Options: openapi3filter.Options{
			// Xác thực do Authenticate đảm nhiệm; bộ kiểm tra chỉ kiểm tra cú pháp và lược đồ.
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
		ErrorHandlerWithOpts: func(_ context.Context, err error, w http.ResponseWriter, _ *http.Request, o middleware.ErrorHandlerOpts) {
			if o.StatusCode == http.StatusNotFound || o.MatchedRoute == nil {
				problem.Write(w, problem.ErrNotFound)
				return
			}
			problem.Write(w, problem.Validation(err.Error()))
		},
	})

	log := opts.Log
	if log == nil {
		log = slog.Default()
	}
	writeErr := func(w http.ResponseWriter, r *http.Request, err error) {
		var p *problem.Error
		if !errors.As(err, &p) {
			log.ErrorContext(r.Context(), "lỗi xử lý yêu cầu", "method", r.Method, "path", r.URL.Path, "err", err)
		}
		problem.Write(w, err)
	}
	apiHandler := api.HandlerWithOptions(
		api.NewStrictHandlerWithOptions(server, nil, api.StrictHTTPServerOptions{
			RequestErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
				problem.Write(w, problem.Validation(err.Error()))
			},
			ResponseErrorHandlerFunc: writeErr,
		}),
		api.StdHTTPServerOptions{ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			problem.Write(w, problem.Validation(err.Error()))
		}},
	)

	var protected http.Handler = validate(apiHandler)
	if opts.Authenticate != nil {
		protected = opts.Authenticate(protected)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	if opts.Terminal != nil {
		mux.Handle("GET /v1/terminal", opts.Terminal)
	}
	mux.Handle("/", protected)
	return mux, nil
}
