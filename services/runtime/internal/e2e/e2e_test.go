// Package e2e kiểm thử runtime-api chế độ docker qua HTTP và WebSocket với PostgreSQL, RabbitMQ, Docker thật
// (specs/001-runtime-docker-mode, tiêu chí D2–D11).
package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/auth"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/driver/docker"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/events"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/httpapi"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/session"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/store"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/terminal"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/testenv"
)

const (
	labImage   = "ghcr.io/crucible-project/crucible-labs/ubuntu-nginx:1.0.0"
	toolsImage = "ghcr.io/crucible-project/crucible/tools:1.0.0"
	m2         = "service-account-content-service"
)

// clock là đồng hồ dịch được để kiểm thử hết hạn và không hoạt động.
type clock struct{ offset atomic.Int64 }

func (c *clock) now() time.Time          { return time.Now().UTC().Add(time.Duration(c.offset.Load())) }
func (c *clock) advance(d time.Duration) { c.offset.Add(int64(d)) }

type env struct {
	t        *testing.T
	srv      *httptest.Server
	store    *store.Store
	sessions *session.Service
	drv      *docker.Driver
	clock    *clock
	events   <-chan amqp.Delivery
	cancel   context.CancelFunc
}

// testAuth thay xác thực JWT (đã kiểm thử ở package auth): header X-Sub, X-Scopes, X-Roles.
func testAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub := r.Header.Get("X-Sub")
		if sub == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		p := auth.NewPrincipal(sub, strings.Fields(r.Header.Get("X-Scopes")), strings.Fields(r.Header.Get("X-Roles")))
		next.ServeHTTP(w, r.WithContext(auth.WithPrincipal(r.Context(), p)))
	})
}

func setup(t *testing.T) *env {
	t.Helper()
	if os.Getenv("SKIP_DOCKER_TESTS") != "" {
		t.Skip("SKIP_DOCKER_TESTS được đặt")
	}
	ctx := context.Background()
	drv, err := docker.New(ctx, toolsImage, slog.Default())
	if err != nil {
		t.Skipf("không có Docker Engine: %v", err)
	}
	st, err := store.Open(ctx, testenv.Postgres(t))
	if err != nil {
		t.Fatal(err)
	}
	amqpURL := testenv.RabbitMQ(t)

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatal(err)
	}
	ch, _ := conn.Channel()
	_ = ch.ExchangeDeclare("crucible.events", amqp.ExchangeTopic, true, false, false, false, nil)
	q, _ := ch.QueueDeclare("", false, true, true, false, nil)
	_ = ch.QueueBind(q.Name, "runtime.#", "crucible.events", false, nil)
	msgs, _ := ch.Consume(q.Name, "", true, true, false, false, nil)

	clk := &clock{}
	sessions := session.New(st, drv, slog.Default(), session.Options{
		SetupTimeout: 20 * time.Second, ProvisionTimeout: 2 * time.Minute, EventsSource: "/crucible/runtime-api",
		ImagePrefixes: []string{"ghcr.io/crucible-project/crucible-labs/"},
		Now:           clk.now,
	})
	term := &terminal.Handler{Store: st, Sessions: sessions, Driver: drv, Events: events.Writer{Source: "/crucible/runtime-api"},
		OriginPatterns: []string{"*"}, Log: slog.Default()}
	sessions.SetNotifier(term)
	srv := httptest.NewUnstartedServer(nil)
	handler, err := httpapi.NewHandler(&httpapi.Server{Sessions: sessions, TerminalURL: "ws://" + srv.Listener.Addr().String() + "/v1/terminal"},
		httpapi.Options{Authenticate: testAuth, Terminal: term})
	if err != nil {
		t.Fatal(err)
	}
	srv.Config.Handler = handler
	srv.Start()

	bg, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	relay := &events.Relay{Pool: st.Pool, URL: amqpURL, Exchange: "crucible.events", Interval: 200 * time.Millisecond, Log: slog.Default()}
	wg.Go(func() { relay.Run(bg) })
	wg.Go(func() { _ = sessions.Run(bg) })
	t.Cleanup(func() {
		srv.Close()
		cancel()
		wg.Wait()
		sessions.Wait()
		for _, s := range mustActive(t, st) {
			if s.ContainerID != nil {
				_ = drv.Destroy(context.Background(), *s.ContainerID)
			}
		}
		_ = conn.Close()
		st.Close()
		_ = drv.Close()
	})
	return &env{t: t, srv: srv, store: st, sessions: sessions, drv: drv, clock: clk, events: msgs, cancel: cancel}
}

func mustActive(t *testing.T, st *store.Store) []*store.Session {
	list, err := store.ListActiveSessions(context.Background(), st.Pool)
	if err != nil {
		t.Logf("đọc phiên: %v", err)
	}
	return list
}

type caller struct {
	sub    string
	scopes string
	roles  string
}

var service = caller{sub: m2, scopes: "runtime.sessions.write runtime.sessions.read runtime.exec"}

func (e *env) do(c caller, method, path string, body any, headers ...string) (int, map[string]any) {
	e.t.Helper()
	var rd io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rd = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, e.srv.URL+path, rd)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Sub", c.sub)
	req.Header.Set("X-Scopes", c.scopes)
	req.Header.Set("X-Roles", c.roles)
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer res.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(res.Body).Decode(&out)
	return res.StatusCode, out
}

func environment(mutate func(map[string]any)) map[string]any {
	envm := map[string]any{
		"image": labImage, "profile": "small", "ttlMinutes": 60, "idleMinutes": 15, "maxExtensionMinutes": 20,
		"network": "none",
		"files":   []map[string]any{{"path": "/srv/lab/readme.txt", "content": "xin chao\n", "mode": "0644"}},
		"setup":   "rm -f /etc/nginx/sites-enabled/default",
	}
	if mutate != nil {
		mutate(envm)
	}
	return envm
}

func (e *env) create(owner uuid.UUID, priority string, envm map[string]any, key ...string) (int, map[string]any) {
	k := uuid.NewString()
	if len(key) > 0 {
		k = key[0]
	}
	return e.do(service, http.MethodPost, "/v1/sessions", map[string]any{
		"ownerUserId": owner, "priority": priority, "quotaAccountRef": "class:1",
		"labels": map[string]string{"attemptId": "a-" + owner.String()[:8]}, "environment": envm,
	}, "Idempotency-Key", k)
}

func (e *env) waitPhase(id string, phases ...string) map[string]any {
	e.t.Helper()
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		_, s := e.do(service, http.MethodGet, "/v1/sessions/"+id, nil)
		for _, p := range phases {
			if s["phase"] == p {
				return s
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	e.t.Fatalf("phiên %s không tới pha %v", id, phases)
	return nil
}

func (e *env) waitEvent(typ, subject string) map[string]any {
	e.t.Helper()
	timeout := time.After(30 * time.Second)
	for {
		select {
		case m := <-e.events:
			var ev map[string]any
			_ = json.Unmarshal(m.Body, &ev)
			if ev["type"] == typ && ev["subject"] == subject {
				return ev["data"].(map[string]any)
			}
		case <-timeout:
			e.t.Fatalf("không nhận được sự kiện %s của %s", typ, subject)
		}
	}
}

func (e *env) exec(id, purpose, script string, timeout int) (int, map[string]any) {
	return e.do(service, http.MethodPost, "/v1/sessions/"+id+"/exec",
		map[string]any{"purpose": purpose, "script": script, "timeoutSeconds": timeout})
}

func TestSessionLifecycle(t *testing.T) {
	e := setup(t)
	owner := uuid.New()
	student := caller{sub: owner.String()}
	key := uuid.NewString()

	code, created := e.create(owner, "PRACTICE", environment(nil), key)
	if code != http.StatusAccepted || created["phase"] != "PENDING" {
		t.Fatalf("D2 tạo phiên: %d %v", code, created)
	}
	id := created["id"].(string)

	t.Run("D9 lặp lại an toàn và một phiên mỗi loại", func(t *testing.T) {
		code, again := e.create(owner, "PRACTICE", environment(nil), key)
		if code != http.StatusAccepted || again["id"] != id {
			t.Errorf("cùng khóa, cùng nội dung phải trả phản hồi cũ: %d %v", code, again)
		}
		code, body := e.create(owner, "PRACTICE", environment(func(m map[string]any) { m["ttlMinutes"] = 30 }), key)
		if code != http.StatusUnprocessableEntity || body["code"] != "IDEMPOTENCY_KEY_REUSED" {
			t.Errorf("cùng khóa, khác nội dung: %d %v", code, body)
		}
		code, body = e.create(owner, "PRACTICE", environment(nil))
		if code != http.StatusConflict || body["code"] != "SESSION_LIMIT" {
			t.Errorf("BR-M1.1 phiên thứ hai: %d %v", code, body)
		}
	})

	ready := e.waitPhase(id, "READY")
	data := e.waitEvent(events.SessionReady, id)
	if data["labels"].(map[string]any)["attemptId"] != "a-"+owner.String()[:8] || ready["expiresAt"] == nil {
		t.Fatalf("D2 sự kiện ready thiếu labels hoặc phiên thiếu expiresAt: %v %v", data, ready)
	}

	t.Run("phân quyền", func(t *testing.T) {
		if code, _ := e.do(student, http.MethodGet, "/v1/sessions/"+id, nil); code != http.StatusOK {
			t.Errorf("chủ phiên đọc được phiên: %d", code)
		}
		if code, _ := e.do(caller{sub: uuid.NewString()}, http.MethodGet, "/v1/sessions/"+id, nil); code != http.StatusNotFound {
			t.Errorf("người khác nhận 404: %d", code)
		}
		if code, _ := e.do(student, http.MethodPost, "/v1/sessions/"+id+"/exec",
			map[string]any{"purpose": "CHECK", "script": "true", "timeoutSeconds": 5}); code != http.StatusForbidden {
			t.Errorf("học viên không được exec: %d", code)
		}
		if code, _ := e.do(caller{sub: uuid.NewString()}, http.MethodPost, "/v1/sessions/"+id+"/terminal-tickets", nil); code != http.StatusNotFound {
			t.Errorf("người khác không lấy được vé: %d", code)
		}
		if code, body := e.do(service, http.MethodGet, "/v1/capacity?from=2026-10-18T00:00:00Z&to=2026-10-18T01:00:00Z&profile=small", nil); code != http.StatusNotImplemented {
			t.Errorf("operation ngoài phạm vi trả 501: %d %v", code, body)
		}
	})

	t.Run("D4 D5 exec", func(t *testing.T) {
		code, res := e.exec(id, "CHECK", "cat /srv/lab/readme.txt; test ! -e /etc/nginx/sites-enabled/default", 5)
		if code != http.StatusOK || res["exitCode"].(float64) != 0 || res["stdout"] != "xin chao\n" {
			t.Fatalf("CHECK: %d %v", code, res)
		}
		if code, body := e.exec(id, "CHECK", "true", 30); code != http.StatusBadRequest {
			t.Errorf("CHECK quá 15 giây phải bị từ chối: %d %v", code, body)
		}
		var wg sync.WaitGroup
		codes := make([]int, 2)
		for i := range 2 {
			wg.Go(func() {
				time.Sleep(time.Duration(i) * 300 * time.Millisecond)
				codes[i], _ = e.exec(id, "CHECK", "sleep 2", 5)
			})
		}
		wg.Wait()
		if !(codes[0] == http.StatusOK && codes[1] == http.StatusConflict) {
			t.Errorf("BR-M1.16 exec đồng thời: %v", codes)
		}
	})

	t.Run("D6 terminal", func(t *testing.T) {
		ticket := func() string {
			code, tk := e.do(student, http.MethodPost, "/v1/sessions/"+id+"/terminal-tickets", nil)
			if code != http.StatusCreated {
				t.Fatalf("lấy vé: %d %v", code, tk)
			}
			return tk["url"].(string)
		}
		dial := func(url string) *websocket.Conn {
			c, _, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{Subprotocols: []string{terminal.Subprotocol}})
			if err != nil {
				t.Fatal(err)
			}
			return c
		}
		url := ticket()
		c1 := dial(url)
		out := readUntil(t, c1, func(b string) bool { return strings.Contains(b, `"type":"status"`) })
		if !strings.Contains(out, `"phase":"READY"`) {
			t.Errorf("khung status đầu tiên: %q", out)
		}
		_ = c1.Write(context.Background(), websocket.MessageText, []byte(`{"type":"resize","cols":100,"rows":30}`))
		_ = c1.Write(context.Background(), websocket.MessageBinary, []byte("echo ket-qua-$((6*7)); stty size\n"))
		readUntil(t, c1, func(b string) bool { return strings.Contains(b, "ket-qua-42") && strings.Contains(b, "30 100") })

		reused := dial(url)
		if code := closeCode(reused); code != terminal.CloseBadTicket {
			t.Errorf("vé dùng lại phải đóng với 4001, nhận %d", code)
		}

		vcode, vt := e.do(caller{sub: "service-account-learning-service", scopes: "runtime.view"}, http.MethodPost,
			"/v1/sessions/"+id+"/view-tickets", map[string]any{"viewerUserId": uuid.NewString()})
		if vcode != http.StatusCreated {
			t.Fatalf("vé xem: %d %v", vcode, vt)
		}
		viewer := dial(vt["url"].(string))
		readUntil(t, viewer, func(b string) bool { return strings.Contains(b, "ket-qua-42") }) // phát lại đầu ra gần nhất

		c2 := dial(ticket())
		if code := closeCode(c1); code != terminal.CloseReplaced {
			t.Errorf("BR-M1.14 kết nối cũ phải đóng với 4009, nhận %d", code)
		}
		_ = c2.Write(context.Background(), websocket.MessageBinary, []byte("echo tu-c2\n"))
		readUntil(t, viewer, func(b string) bool { return strings.Contains(b, "tu-c2") })
		_ = viewer.Write(context.Background(), websocket.MessageBinary, []byte("echo khong-duoc-chay\n"))
		time.Sleep(500 * time.Millisecond)
		_ = c2.Write(context.Background(), websocket.MessageBinary, []byte("echo moc-ket-thuc\n"))
		if out := readUntil(t, c2, func(b string) bool { return strings.Contains(b, "moc-ket-thuc\r\n") }); strings.Contains(out, "khong-duoc-chay") {
			t.Errorf("dữ liệu từ kết nối xem không được đưa vào terminal: %q", out)
		}

		t.Run("D7 gia hạn, sắp hết hạn, hết hạn", func(t *testing.T) {
			if code, body := e.do(student, http.MethodPost, "/v1/sessions/"+id+"/extend", map[string]any{"minutes": 15}); code != http.StatusOK {
				t.Fatalf("gia hạn: %d %v", code, body)
			}
			if code, body := e.do(student, http.MethodPost, "/v1/sessions/"+id+"/extend", map[string]any{"minutes": 10}); code != http.StatusConflict || body["code"] != "EXTENSION_LIMIT" {
				t.Errorf("vượt maxExtensionMinutes: %d %v", code, body)
			}
			// Gõ phím sau mỗi lần tua đồng hồ để phiên không bị kết thúc vì không hoạt động (BR-M1.9).
			typeKey := func() {
				_ = c2.Write(context.Background(), websocket.MessageBinary, []byte("\n"))
				time.Sleep(300 * time.Millisecond)
			}
			e.clock.advance(71 * time.Minute) // 60 + 15 phút gia hạn, còn 4 phút
			typeKey()
			e.sessions.Tick(context.Background())
			if s := e.waitPhase(id, "EXPIRING"); s["phase"] != "EXPIRING" {
				t.Fatal("phải chuyển EXPIRING")
			}
			e.waitEvent(events.SessionExpiring, id)
			readUntil(t, c2, func(b string) bool { return strings.Contains(b, `"type":"expiring"`) })
			e.clock.advance(5 * time.Minute)
			typeKey()
			e.sessions.Tick(context.Background())
			s := e.waitPhase(id, "TERMINATED")
			if s["endReason"] != "EXPIRED" {
				t.Errorf("lý do kết thúc: %v", s["endReason"])
			}
			ev := e.waitEvent(events.SessionTerminated, id)
			if ev["endReason"] != "EXPIRED" || ev["resourceUsage"] == nil {
				t.Errorf("sự kiện terminated: %v", ev)
			}
			if code := closeCode(c2); code != terminal.CloseEnded {
				t.Errorf("terminal phải đóng với 1000 khi phiên kết thúc, nhận %d", code)
			}
		})
	})

}

func TestFailuresAndExits(t *testing.T) {
	e := setup(t)

	t.Run("D10 setup lỗi", func(t *testing.T) {
		code, s := e.create(uuid.New(), "PRACTICE", environment(func(m map[string]any) { m["setup"] = "echo sai-cau-hinh >&2; exit 2" }))
		if code != http.StatusAccepted {
			t.Fatal(code, s)
		}
		id := s["id"].(string)
		f := e.waitPhase(id, "FAILED")
		if f["endReason"] != "PROVISION_FAILED" || !strings.Contains(fmt.Sprint(f["failureDetail"]), "sai-cau-hinh") {
			t.Errorf("phiên lỗi: %v", f)
		}
		e.waitEvent(events.SessionFailed, id)
	})

	t.Run("image ngoài danh sách cho phép", func(t *testing.T) {
		code, body := e.create(uuid.New(), "PRACTICE", environment(func(m map[string]any) { m["image"] = "nginx:latest" }))
		if code != http.StatusUnprocessableEntity || body["code"] != "ENVIRONMENT_NOT_ALLOWED" {
			t.Errorf("mong đợi 422 ENVIRONMENT_NOT_ALLOWED: %d %v", code, body)
		}
	})

	t.Run("BR-M1.26 capability ngoài danh sách", func(t *testing.T) {
		code, body := e.create(uuid.New(), "PRACTICE", environment(func(m map[string]any) { m["capabilities"] = []string{"NET_ADMIN"} }))
		if code != http.StatusUnprocessableEntity || body["code"] != "ENVIRONMENT_NOT_ALLOWED" {
			t.Errorf("mong đợi 422 ENVIRONMENT_NOT_ALLOWED: %d %v", code, body)
		}
	})

	t.Run("D8 hết bộ nhớ", func(t *testing.T) {
		_, s := e.create(uuid.New(), "EXAM", environment(nil))
		id := s["id"].(string)
		e.waitPhase(id, "READY")
		_, _ = e.exec(id, "RESTORE", "head -c 600m /dev/zero | tail > /dev/null", 30)
		f := e.waitPhase(id, "TERMINATED")
		usage, _ := f["resourceUsage"].(map[string]any)
		if f["endReason"] != "OUT_OF_MEMORY" || usage["oomKilled"] != true {
			t.Errorf("phiên hết bộ nhớ: %v", f)
		}
	})

	t.Run("BR-M1.9 không hoạt động và ADMIN kết thúc", func(t *testing.T) {
		_, s := e.create(uuid.New(), "PRACTICE", environment(nil))
		id := s["id"].(string)
		e.waitPhase(id, "READY")
		e.clock.advance(16 * time.Minute)
		e.sessions.Tick(context.Background())
		f := e.waitPhase(id, "TERMINATED")
		e.clock.advance(-16 * time.Minute)
		if f["endReason"] != "IDLE" {
			t.Errorf("lý do: %v", f["endReason"])
		}
		code, again := e.do(service, http.MethodPost, "/v1/sessions/"+id+"/terminate", map[string]any{"reason": "CANCELLED"})
		if code != http.StatusOK || again["endReason"] != "IDLE" {
			t.Errorf("kết thúc phiên đã kết thúc trả trạng thái hiện tại: %d %v", code, again)
		}
	})

	t.Run("D11 nhận lại sau khởi động", func(t *testing.T) {
		_, s := e.create(uuid.New(), "PRACTICE", environment(nil))
		id := s["id"].(string)
		e.waitPhase(id, "READY")
		sess, _ := store.GetSession(context.Background(), e.store.Pool, uuid.MustParse(id))
		restarted := session.New(e.store, e.drv, slog.Default(), session.Options{EventsSource: "/x"})
		if err := restarted.Reconcile(context.Background()); err != nil {
			t.Fatal(err)
		}
		if s := e.waitPhase(id, "READY"); s["phase"] != "READY" {
			t.Errorf("phiên đang chạy phải được giữ")
		}
		_ = e.drv.Destroy(context.Background(), *sess.ContainerID)
		if err := restarted.Reconcile(context.Background()); err != nil {
			t.Fatal(err)
		}
		if f := e.waitPhase(id, "TERMINATED"); f["endReason"] != "ADMIN" && f["endReason"] != "ENVIRONMENT_FAILED" {
			t.Errorf("phiên mất môi trường: %v", f["endReason"])
		}
	})
}

func readUntil(t *testing.T, c *websocket.Conn, done func(string) bool) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var all strings.Builder
	for {
		_, data, err := c.Read(ctx)
		if err != nil {
			t.Fatalf("đọc terminal: %v; đã nhận %q", err, all.String())
		}
		all.Write(data)
		if done(all.String()) {
			return all.String()
		}
	}
}

func closeCode(c *websocket.Conn) websocket.StatusCode {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		if _, _, err := c.Read(ctx); err != nil {
			var ce websocket.CloseError
			if errors.As(err, &ce) {
				return ce.Code
			}
			return -1
		}
	}
}
