package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/api"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/auth"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/problem"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/session"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/store"
)

// Server hiện thực các operation đã có; operation khác trả 501 qua Unimplemented.
type Server struct {
	Unimplemented
	Sessions    *session.Service
	TerminalURL string
}

var _ api.StrictServerInterface = (*Server)(nil)

func principal(ctx context.Context) (auth.Principal, error) {
	p, ok := auth.FromContext(ctx)
	if !ok {
		return p, problem.ErrUnauthenticated
	}
	return p, nil
}

func isOwner(p auth.Principal, sess *store.Session) bool {
	id, ok := p.UserID()
	return ok && id == sess.OwnerUserID
}

// readable trả phiên khi người gọi là chủ phiên hoặc có scope đọc; người khác nhận 404 như phiên không tồn tại.
func (s *Server) readable(ctx context.Context, id uuid.UUID, scopes ...string) (*store.Session, auth.Principal, error) {
	p, err := principal(ctx)
	if err != nil {
		return nil, p, err
	}
	sess, err := s.Sessions.Get(ctx, id)
	if err != nil {
		return nil, p, err
	}
	if isOwner(p, sess) || p.HasRole(auth.RoleAdmin) {
		return sess, p, nil
	}
	for _, sc := range scopes {
		if p.HasScope(sc) {
			return sess, p, nil
		}
	}
	return nil, p, problem.ErrNotFound
}

func requireScope(ctx context.Context, scope string) (auth.Principal, error) {
	p, err := principal(ctx)
	if err != nil {
		return p, err
	}
	if !p.HasScope(scope) {
		return p, problem.ErrForbidden
	}
	return p, nil
}

// CreateSession tạo phiên (scope runtime.sessions.write).
func (s *Server) CreateSession(ctx context.Context, req api.CreateSessionRequestObject) (api.CreateSessionResponseObject, error) {
	p, err := requireScope(ctx, auth.ScopeSessionsWrite)
	if err != nil {
		return nil, err
	}
	res, err := s.Sessions.Create(ctx, p.Subject, req.Params.IdempotencyKey, *req.Body)
	if err != nil {
		return nil, err
	}
	loc := "/v1/sessions/" + res.Session.Id.String()
	return api.CreateSession202JSONResponse{Body: res.Session, Headers: api.CreateSession202ResponseHeaders{Location: &loc}}, nil
}

// GetSession đọc phiên: chủ phiên hoặc scope runtime.sessions.read.
func (s *Server) GetSession(ctx context.Context, req api.GetSessionRequestObject) (api.GetSessionResponseObject, error) {
	sess, _, err := s.readable(ctx, req.SessionId, auth.ScopeSessionsRead, auth.ScopeSessionsWrite)
	if err != nil {
		return nil, err
	}
	return api.GetSession200JSONResponse(session.ToAPI(sess)), nil
}

// ExtendSession gia hạn phiên: chủ phiên hoặc scope runtime.sessions.write (BR-M1.10).
func (s *Server) ExtendSession(ctx context.Context, req api.ExtendSessionRequestObject) (api.ExtendSessionResponseObject, error) {
	if _, _, err := s.readable(ctx, req.SessionId, auth.ScopeSessionsWrite); err != nil {
		return nil, err
	}
	sess, err := s.Sessions.Extend(ctx, req.SessionId, req.Body.Minutes)
	if err != nil {
		return nil, err
	}
	return api.ExtendSession200JSONResponse(session.ToAPI(sess)), nil
}

// TerminateSession kết thúc phiên: scope runtime.sessions.write hoặc vai trò ADMIN.
func (s *Server) TerminateSession(ctx context.Context, req api.TerminateSessionRequestObject) (api.TerminateSessionResponseObject, error) {
	p, err := principal(ctx)
	if err != nil {
		return nil, err
	}
	if !p.HasScope(auth.ScopeSessionsWrite) && !p.HasRole(auth.RoleAdmin) {
		return nil, problem.ErrForbidden
	}
	sess, err := s.Sessions.Terminate(ctx, req.SessionId, api.EndReason(req.Body.Reason))
	if err != nil {
		return nil, err
	}
	return api.TerminateSession200JSONResponse(session.ToAPI(sess)), nil
}

// ExecInSession chạy lệnh trong phiên (scope runtime.exec).
func (s *Server) ExecInSession(ctx context.Context, req api.ExecInSessionRequestObject) (api.ExecInSessionResponseObject, error) {
	if _, err := requireScope(ctx, auth.ScopeExec); err != nil {
		return nil, err
	}
	res, err := s.Sessions.Exec(ctx, req.SessionId, *req.Body)
	if err != nil {
		return nil, err
	}
	return api.ExecInSession200JSONResponse{
		ExitCode: res.ExitCode, TimedOut: res.TimedOut, Stdout: res.Stdout, Stderr: res.Stderr,
		DurationMs: int(res.Duration.Milliseconds()),
	}, nil
}

// CreateTerminalTicket cấp vé tương tác cho chủ phiên (BR-M1.13).
func (s *Server) CreateTerminalTicket(ctx context.Context, req api.CreateTerminalTicketRequestObject) (api.CreateTerminalTicketResponseObject, error) {
	p, err := principal(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := s.Sessions.Get(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	if !isOwner(p, sess) {
		return nil, problem.ErrNotFound
	}
	token, expires, err := s.Sessions.IssueTicket(ctx, sess.ID, string(api.INTERACTIVE), sess.OwnerUserID)
	if err != nil {
		return nil, err
	}
	return api.CreateTerminalTicket201JSONResponse{Ticket: token, Mode: api.INTERACTIVE, ExpiresAt: expires,
		Url: s.TerminalURL + "?ticket=" + token}, nil
}

// CreateViewTicket cấp vé xem cho giảng viên qua M3 (scope runtime.view).
func (s *Server) CreateViewTicket(ctx context.Context, req api.CreateViewTicketRequestObject) (api.CreateViewTicketResponseObject, error) {
	if _, err := requireScope(ctx, auth.ScopeView); err != nil {
		return nil, err
	}
	token, expires, err := s.Sessions.IssueTicket(ctx, req.SessionId, string(api.VIEW), req.Body.ViewerUserId)
	if err != nil {
		var pe *problem.Error
		// Hợp đồng chỉ có 404 cho operation này: phiên không chạy được coi như không xem được.
		if asProblem(err, &pe) && pe.Status == http.StatusConflict {
			return nil, problem.ErrNotFound
		}
		return nil, err
	}
	return api.CreateViewTicket201JSONResponse{Ticket: token, Mode: api.VIEW, ExpiresAt: expires,
		Url: s.TerminalURL + "?ticket=" + token}, nil
}

func asProblem(err error, target **problem.Error) bool { return errors.As(err, target) }
