// Package terminal phục vụ terminal WebSocket theo CR-DOC-03 mục 6: vé dùng một lần, một kết nối tương tác mỗi phiên,
// kết nối xem chỉ nhận đầu ra, khung trạng thái từ máy chủ. Ở chế độ docker, terminal chạy trong runtime-api thay cho
// runtime-gateway.
package terminal

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/api"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/driver"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/events"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/session"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/store"
)

// Subprotocol và mã đóng (CR-DOC-03 mục 6).
const (
	Subprotocol = "crucible.terminal.v1"

	CloseEnded       websocket.StatusCode = 1000
	CloseBadTicket   websocket.StatusCode = 4001
	CloseNoSession   websocket.StatusCode = 4004
	CloseReplaced    websocket.StatusCode = 4009
	CloseRestarting  websocket.StatusCode = 4010
	pingInterval                          = 20 * time.Second
	pongTimeout                           = 60 * time.Second
	replaySize                            = 64 << 10
	maxShellRestarts                      = 20
)

// Handler là http.Handler của /v1/terminal và đồng thời là session.Notifier.
type Handler struct {
	Store          *store.Store
	Sessions       *session.Service
	Driver         driver.Driver
	Events         events.Writer
	OriginPatterns []string
	Log            *slog.Logger

	mu    sync.Mutex
	rooms map[uuid.UUID]*room
}

type room struct {
	interactive *websocket.Conn
	viewers     map[*websocket.Conn]bool
	replay      []byte
}

func (h *Handler) room(id uuid.UUID) *room {
	if h.rooms == nil {
		h.rooms = map[uuid.UUID]*room{}
	}
	r, ok := h.rooms[id]
	if !ok {
		r = &room{viewers: map[*websocket.Conn]bool{}}
		h.rooms[id] = r
	}
	return r
}

func (h *Handler) conns(id uuid.UUID) []*websocket.Conn {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[id]
	if !ok {
		return nil
	}
	var out []*websocket.Conn
	if r.interactive != nil {
		out = append(out, r.interactive)
	}
	for c := range r.viewers {
		out = append(out, c)
	}
	return out
}

func writeJSON(ctx context.Context, c *websocket.Conn, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, b)
}

// Notify gửi khung trạng thái tới mọi kết nối của phiên.
func (h *Handler) Notify(id uuid.UUID, frame any) {
	for _, c := range h.conns(id) {
		_ = writeJSON(context.Background(), c, frame)
	}
}

// CloseSession gửi khung closing rồi đóng mọi kết nối của phiên đã kết thúc.
func (h *Handler) CloseSession(id uuid.UUID, reason string) {
	for _, c := range h.conns(id) {
		_ = writeJSON(context.Background(), c, map[string]any{"type": "closing", "reason": reason})
		_ = c.Close(CloseEnded, reason)
	}
	h.mu.Lock()
	delete(h.rooms, id)
	h.mu.Unlock()
}

// Shutdown báo mọi kết nối máy chủ đang dừng (mã 4010) để giao diện kết nối lại.
func (h *Handler) Shutdown() {
	h.mu.Lock()
	var all []*websocket.Conn
	for _, r := range h.rooms {
		if r.interactive != nil {
			all = append(all, r.interactive)
		}
		for c := range r.viewers {
			all = append(all, c)
		}
	}
	h.mu.Unlock()
	for _, c := range all {
		_ = c.Close(CloseRestarting, "server restarting")
	}
}

func clientIP(r *http.Request) string {
	if f := r.Header.Get("X-Forwarded-For"); f != "" {
		return strings.TrimSpace(strings.Split(f, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{Subprotocol}, OriginPatterns: h.OriginPatterns,
	})
	if err != nil {
		return
	}
	ctx := r.Context()
	ticket, err := store.ConsumeTicket(ctx, h.Store.Pool, r.URL.Query().Get("ticket"))
	if err != nil {
		_ = c.Close(CloseBadTicket, "invalid ticket")
		return
	}
	sess, err := h.Sessions.Get(ctx, ticket.SessionID)
	if err != nil || (sess.Phase != store.PhaseReady && sess.Phase != store.PhaseExpiring) || sess.ContainerID == nil {
		_ = c.Close(CloseNoSession, "session not active")
		return
	}
	err = h.Store.InTx(ctx, func(tx pgx.Tx) error {
		_, err := h.Events.Write(ctx, tx, events.TerminalConnected, sess.ID.String(), map[string]any{
			"sessionId": sess.ID, "mode": ticket.Mode, "userId": ticket.UserID, "clientIp": clientIP(r),
			"connectedAt": time.Now().UTC(), "labels": sess.Labels,
		})
		return err
	})
	if err != nil {
		h.Log.Warn("ghi sự kiện terminal.connected", "session", sess.ID, "err", err)
	}
	_ = writeJSON(ctx, c, session.StatusFrame(sess))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go h.keepAlive(ctx, c, cancel)
	if ticket.Mode == string(api.VIEW) {
		h.serveViewer(ctx, sess.ID, c)
		return
	}
	h.serveInteractive(ctx, sess, c)
}

// keepAlive gửi ping mỗi 20 giây; không có pong sau 60 giây thì đóng kết nối.
func (h *Handler) keepAlive(ctx context.Context, c *websocket.Conn, cancel context.CancelFunc) {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			pctx, pcancel := context.WithTimeout(ctx, pongTimeout)
			err := c.Ping(pctx)
			pcancel()
			if err != nil && ctx.Err() == nil {
				_ = c.Close(websocket.StatusGoingAway, "ping timeout")
				cancel()
				return
			}
		}
	}
}

func (h *Handler) serveViewer(ctx context.Context, id uuid.UUID, c *websocket.Conn) {
	h.mu.Lock()
	rm := h.room(id)
	replay := append([]byte(nil), rm.replay...)
	rm.viewers[c] = true
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.room(id).viewers, c)
		h.mu.Unlock()
	}()
	if len(replay) > 0 {
		_ = c.Write(ctx, websocket.MessageBinary, replay)
	}
	// Kết nối xem không gửi dữ liệu vào terminal; vòng đọc chỉ để xử lý ping và đóng kết nối.
	for {
		if _, _, err := c.Read(ctx); err != nil {
			return
		}
	}
}

type resize struct {
	Type string `json:"type"`
	Cols uint   `json:"cols"`
	Rows uint   `json:"rows"`
}

func (h *Handler) serveInteractive(ctx context.Context, sess *store.Session, c *websocket.Conn) {
	id := sess.ID
	h.mu.Lock()
	rm := h.room(id)
	old := rm.interactive
	rm.interactive = c
	h.mu.Unlock()
	if old != nil {
		_ = old.Close(CloseReplaced, "replaced by newer connection") // BR-M1.14
	}
	defer func() {
		h.mu.Lock()
		if r, ok := h.rooms[id]; ok && r.interactive == c {
			r.interactive = nil
		}
		h.mu.Unlock()
	}()

	var (
		smu    sync.Mutex
		stream driver.Stream
		size   = resize{Cols: 80, Rows: 24}
	)
	current := func() driver.Stream {
		smu.Lock()
		defer smu.Unlock()
		return stream
	}
	open := func() error {
		st, err := h.Driver.Attach(ctx, *sess.ContainerID, size.Cols, size.Rows)
		if err != nil {
			return err
		}
		smu.Lock()
		stream = st
		smu.Unlock()
		return nil
	}
	if err := open(); err != nil {
		h.Log.Warn("mở shell", "session", id, "err", err)
		_ = c.Close(CloseNoSession, "cannot open shell")
		return
	}
	defer func() {
		if st := current(); st != nil {
			_ = st.Close()
		}
	}()

	// Đầu ra của shell tới trình duyệt và người xem. Shell thoát (học viên gõ exit) thì mở shell mới.
	go func() {
		buf := make([]byte, 32<<10)
		for restarts := 0; ctx.Err() == nil; {
			n, err := current().Read(buf)
			if n > 0 {
				h.broadcast(ctx, id, c, buf[:n])
			}
			if err == nil {
				continue
			}
			if ctx.Err() != nil || !errors.Is(err, io.EOF) || restarts >= maxShellRestarts {
				_ = c.Close(CloseNoSession, "shell closed")
				return
			}
			restarts++
			_ = current().Close()
			if err := open(); err != nil {
				_ = c.Close(CloseNoSession, "cannot open shell")
				return
			}
		}
	}()

	for {
		typ, data, err := c.Read(ctx)
		if err != nil {
			return
		}
		switch typ {
		case websocket.MessageBinary:
			h.Sessions.Touch(id)
			if _, err := current().Write(data); err != nil {
				return
			}
		case websocket.MessageText:
			var m resize
			if json.Unmarshal(data, &m) == nil && m.Type == "resize" && m.Cols > 0 && m.Rows > 0 {
				smu.Lock()
				size = m
				smu.Unlock()
				_ = current().Resize(ctx, m.Cols, m.Rows)
			}
		}
	}
}

func (h *Handler) broadcast(ctx context.Context, id uuid.UUID, interactive *websocket.Conn, p []byte) {
	h.mu.Lock()
	rm := h.room(id)
	rm.replay = append(rm.replay, p...)
	if over := len(rm.replay) - replaySize; over > 0 {
		rm.replay = append([]byte(nil), rm.replay[over:]...)
	}
	viewers := make([]*websocket.Conn, 0, len(rm.viewers))
	for v := range rm.viewers {
		viewers = append(viewers, v)
	}
	h.mu.Unlock()
	_ = interactive.Write(ctx, websocket.MessageBinary, p)
	for _, v := range viewers {
		wctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		_ = v.Write(wctx, websocket.MessageBinary, p)
		cancel()
	}
}

var _ session.Notifier = (*Handler)(nil)
