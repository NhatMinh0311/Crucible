package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/auth"
)

type issuer struct {
	srv    *httptest.Server
	signer jose.Signer
}

func newIssuer(t *testing.T) *issuer {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	jwk := jose.JSONWebKey{Key: &key.PublicKey, KeyID: "k1", Algorithm: "RS256", Use: "sig"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}})
	}))
	t.Cleanup(srv.Close)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key},
		(&jose.SignerOptions{}).WithHeader("kid", "k1").WithType("JWT"))
	if err != nil {
		t.Fatal(err)
	}
	return &issuer{srv: srv, signer: signer}
}

func (i *issuer) token(t *testing.T, claims map[string]any) string {
	body, _ := json.Marshal(claims)
	jws, err := i.signer.Sign(body)
	if err != nil {
		t.Fatal(err)
	}
	s, _ := jws.CompactSerialize()
	return s
}

func TestMiddleware(t *testing.T) {
	is := newIssuer(t)
	v := auth.NewVerifier(context.Background(), is.srv.URL, is.srv.URL, "runtime-api")
	var got auth.Principal
	h := auth.Middleware(v)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ = auth.FromContext(r.Context())
	}))
	call := func(tok string) int {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/sessions/x", nil)
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		h.ServeHTTP(rec, req)
		return rec.Code
	}
	base := func() map[string]any {
		return map[string]any{
			"iss": is.srv.URL, "aud": []string{"runtime-api", "account"}, "sub": "7d0c4a9e-4a3b-4c1d-9e57-3b1f8e2a6c10",
			"exp": time.Now().Add(time.Minute).Unix(), "iat": time.Now().Unix(), "azp": "content-service",
			"scope": "runtime.sessions.write runtime.exec", "realm_access": map[string]any{"roles": []string{"ADMIN"}},
		}
	}

	if code := call(is.token(t, base())); code != http.StatusOK {
		t.Fatalf("token hợp lệ phải qua, nhận %d", code)
	}
	if !got.HasScope(auth.ScopeExec) || got.HasScope(auth.ScopeView) || !got.HasRole(auth.RoleAdmin) || got.ClientID != "content-service" {
		t.Fatalf("người gọi sai: %+v", got)
	}
	if id, ok := got.UserID(); !ok || id.String() != "7d0c4a9e-4a3b-4c1d-9e57-3b1f8e2a6c10" {
		t.Fatalf("sub sai: %v", id)
	}

	cases := map[string]func(map[string]any){
		"sai audience": func(c map[string]any) { c["aud"] = "content-service" },
		"hết hạn":      func(c map[string]any) { c["exp"] = time.Now().Add(-time.Minute).Unix() },
		"sai issuer":   func(c map[string]any) { c["iss"] = "http://khac" },
	}
	for name, mutate := range cases {
		c := base()
		mutate(c)
		if code := call(is.token(t, c)); code != http.StatusUnauthorized {
			t.Errorf("%s: mong đợi 401, nhận %d", name, code)
		}
	}
	if code := call(""); code != http.StatusUnauthorized {
		t.Errorf("thiếu token: mong đợi 401, nhận %d", code)
	}
}
