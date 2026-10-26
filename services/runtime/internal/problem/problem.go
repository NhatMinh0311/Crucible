// Package problem biểu diễn lỗi API theo RFC 9457 với mã lỗi ổn định (CR-DOC-03 mục 3.4).
package problem

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Error là lỗi nghiệp vụ trả về cho bên gọi.
type Error struct {
	Status int
	Code   string
	Detail string
}

func (e *Error) Error() string { return e.Code + ": " + e.Detail }

// New tạo lỗi với mã HTTP, mã lỗi và mô tả.
func New(status int, code, detail string) *Error {
	return &Error{Status: status, Code: code, Detail: detail}
}

// Các lỗi dùng chung.
var (
	ErrUnauthenticated = New(http.StatusUnauthorized, "UNAUTHENTICATED", "Thiếu hoặc sai token")
	ErrForbidden       = New(http.StatusForbidden, "FORBIDDEN", "Token thiếu vai trò hoặc scope cần thiết")
	ErrNotFound        = New(http.StatusNotFound, "NOT_FOUND", "Không tồn tại hoặc không có quyền xem")
	ErrNotImplemented  = New(http.StatusNotImplemented, "NOT_IMPLEMENTED", "Operation chưa được hiện thực")
)

// Validation tạo lỗi 400 VALIDATION_FAILED.
func Validation(detail string) *Error {
	return New(http.StatusBadRequest, "VALIDATION_FAILED", detail)
}

type body struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Detail  string `json:"detail,omitempty"`
	TraceID string `json:"traceId,omitempty"`
}

// Write ghi lỗi dạng application/problem+json. Lỗi không phải *Error được trả thành 500 INTERNAL.
func Write(w http.ResponseWriter, err error) {
	var p *Error
	if !errors.As(err, &p) {
		p = New(http.StatusInternalServerError, "INTERNAL", "Lỗi nội bộ")
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	_ = json.NewEncoder(w).Encode(body{
		Type:   "about:blank",
		Title:  http.StatusText(p.Status),
		Status: p.Status,
		Code:   p.Code,
		Detail: p.Detail,
	})
}
