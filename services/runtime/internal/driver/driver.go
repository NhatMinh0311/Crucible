// Package driver định nghĩa giao diện giữa tầng phiên và nơi chạy môi trường (Docker khi phát triển, Kubernetes trên
// cụm). Tầng phiên không biết môi trường chạy ở đâu.
package driver

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

// File là một tệp ghi vào môi trường trước khi chạy setup.
type File struct {
	Path    string
	Content string
	Mode    os.FileMode
}

// Spec mô tả môi trường cần dựng cho một phiên.
type Spec struct {
	SessionID    uuid.UUID
	Image        string
	MemoryBytes  int64
	NanoCPUs     int64
	Network      string // none hoặc mirror
	Capabilities []string
	Files        []File
	Setup        string
	SetupTimeout time.Duration
}

// Purpose của lệnh exec (CR-DOC-03 mục 7).
const (
	PurposeCheck   = "CHECK"
	PurposeRestore = "RESTORE"
)

// ExecSpec là một lệnh chạy trong môi trường.
type ExecSpec struct {
	Purpose string
	Script  string
	Timeout time.Duration
}

// ExecResult là kết quả lệnh; Stdout, Stderr tối đa MaxOutput byte.
type ExecResult struct {
	ExitCode int
	TimedOut bool
	Stdout   string
	Stderr   string
	Duration time.Duration
}

// MaxOutput là số byte tối đa giữ lại của mỗi luồng đầu ra.
const MaxOutput = 4096

// Stream là luồng terminal có TTY tới shell trong môi trường.
type Stream interface {
	io.ReadWriteCloser
	Resize(ctx context.Context, cols, rows uint) error
}

// Exit báo môi trường dừng ngoài ý muốn (BR-M1.17).
type Exit struct {
	SessionID uuid.UUID
	OOM       bool
	ExitCode  int
}

// Instance là môi trường đang tồn tại, dùng khi nhận lại sau khởi động.
type Instance struct {
	SessionID uuid.UUID
	Handle    string
	Running   bool
}

// ProvisionError là lỗi dựng môi trường do bài lab (setup lỗi, quá hạn, image không có); Detail hiển thị cho bên gọi.
type ProvisionError struct {
	Detail string
}

func (e *ProvisionError) Error() string { return e.Detail }

// Driver dựng và điều khiển môi trường của phiên. Handle là mã của môi trường ở nơi chạy (mã container).
type Driver interface {
	// Provision dựng môi trường, ghi tệp, chạy setup. Khi lỗi, handle (nếu khác rỗng) vẫn cần được Destroy.
	Provision(ctx context.Context, spec Spec) (handle string, err error)
	Exec(ctx context.Context, handle string, e ExecSpec) (ExecResult, error)
	Attach(ctx context.Context, handle string, cols, rows uint) (Stream, error)
	Destroy(ctx context.Context, handle string) error
	MemoryUsage(ctx context.Context, handle string) (int64, error)
	// Watch phát Exit khi môi trường dừng hoặc bị hết bộ nhớ; kênh đóng khi ctx kết thúc.
	Watch(ctx context.Context) <-chan Exit
	List(ctx context.Context) ([]Instance, error)
}
