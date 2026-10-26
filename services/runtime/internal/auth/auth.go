// Package auth kiểm tra access token của Keycloak và gắn người gọi vào context (CR-DOC-03 mục 3.2).
package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/problem"
)

// Scope và vai trò dùng ở M1.
const (
	ScopeSessionsWrite = "runtime.sessions.write"
	ScopeSessionsRead  = "runtime.sessions.read"
	ScopeExec          = "runtime.exec"
	ScopeView          = "runtime.view"
	RoleAdmin          = "ADMIN"
)

// Principal là người gọi đã xác thực: người dùng hoặc tài khoản dịch vụ.
type Principal struct {
	Subject  string
	ClientID string
	scopes   map[string]bool
	roles    map[string]bool
}

// UserID trả `sub` dạng UUID; tài khoản Keycloak luôn có `sub` là UUID.
func (p Principal) UserID() (uuid.UUID, bool) {
	id, err := uuid.Parse(p.Subject)
	return id, err == nil
}

// HasScope cho biết token có scope.
func (p Principal) HasScope(s string) bool { return p.scopes[s] }

// HasRole cho biết người dùng có vai trò realm.
func (p Principal) HasRole(r string) bool { return p.roles[r] }

// NewPrincipal tạo Principal; dùng trong kiểm thử.
func NewPrincipal(subject string, scopes, roles []string) Principal {
	p := Principal{Subject: subject, scopes: map[string]bool{}, roles: map[string]bool{}}
	for _, s := range scopes {
		p.scopes[s] = true
	}
	for _, r := range roles {
		p.roles[r] = true
	}
	return p
}

type ctxKey struct{}

// WithPrincipal gắn người gọi vào context.
func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, ctxKey{}, p)
}

// FromContext lấy người gọi; chỉ có khi yêu cầu đã qua Middleware.
func FromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(ctxKey{}).(Principal)
	return p, ok
}

// Verifier kiểm tra chữ ký, issuer, audience, hạn của token.
type Verifier struct {
	v *oidc.IDTokenVerifier
}

// NewVerifier tạo bộ kiểm tra lấy khóa công khai từ JWKS của realm. Không gọi discovery khi khởi động để dịch vụ chạy
// được khi Keycloak chưa sẵn sàng.
func NewVerifier(ctx context.Context, issuer, jwksURL, audience string) *Verifier {
	keys := oidc.NewRemoteKeySet(ctx, jwksURL)
	return &Verifier{v: oidc.NewVerifier(issuer, keys, &oidc.Config{ClientID: audience})}
}

type claims struct {
	Scope       string `json:"scope"`
	AZP         string `json:"azp"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

// Verify kiểm tra token và trả người gọi.
func (v *Verifier) Verify(ctx context.Context, raw string) (Principal, error) {
	tok, err := v.v.Verify(ctx, raw)
	if err != nil {
		return Principal{}, err
	}
	var c claims
	if err := tok.Claims(&c); err != nil {
		return Principal{}, fmt.Errorf("đọc claim: %w", err)
	}
	p := NewPrincipal(tok.Subject, strings.Fields(c.Scope), c.RealmAccess.Roles)
	p.ClientID = c.AZP
	return p, nil
}

// Middleware từ chối yêu cầu không có token hợp lệ (401 UNAUTHENTICATED) và gắn người gọi vào context.
func Middleware(v *Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || raw == "" {
				problem.Write(w, problem.ErrUnauthenticated)
				return
			}
			p, err := v.Verify(r.Context(), raw)
			if err != nil {
				problem.Write(w, problem.ErrUnauthenticated)
				return
			}
			next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), p)))
		})
	}
}
