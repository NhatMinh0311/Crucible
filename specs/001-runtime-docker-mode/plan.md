# Kế hoạch: `runtime-api` chế độ `docker`

Đặc tả: `spec.md`. Kế hoạch dựng phần khung dùng chung cho cả hai chế độ (API, CSDL, sự kiện, vé, terminal) và một
trình điều khiển `docker`; chế độ `kubernetes` sau này thay trình điều khiển, giữ nguyên phần còn lại.

## 1. Cấu trúc mã

```text
services/runtime/
├── go.mod                      module github.com/NhatMinh0311/crucible/services/runtime, go 1.26
├── cmd/runtime-api/main.go     đọc cấu hình, nối các thành phần, dừng có kiểm soát
├── internal/api/               cấu hình oapi-codegen; mã sinh từ contracts/runtime-api.yaml (commit kèm, CI kiểm tra không lệch)
├── internal/
│   ├── config/                 biến môi trường (CR-MOD-M1 mục 12)
│   ├── auth/                   kiểm tra JWT Keycloak (go-oidc), scope, vai trò, chủ phiên
│   ├── httpapi/                hiện thực interface sinh từ hợp đồng; lỗi RFC 9457; Idempotency-Key
│   ├── session/                vòng đời phiên, quy tắc BR-M1.*, bộ hẹn giờ hết hạn và không hoạt động
│   ├── driver/                 interface Driver
│   │   └── docker/             trình điều khiển Docker Engine API
│   ├── terminal/               WebSocket theo CR-DOC-03 mục 6, một kết nối tương tác, phát cho người xem
│   ├── store/                  PostgreSQL (pgx), migration nhúng (golang-migrate)
│   └── events/                 outbox, chuyển tiếp lên RabbitMQ có publisher confirm
images/
├── crucible-tools/Dockerfile   BusyBox, curl, jq liên kết tĩnh trong /opt/crucible/bin
└── lab-ubuntu-nginx/Dockerfile Ubuntu 24.04, nginx, tini; PID 1 là tini chạy sleep infinity
```

## 2. Lựa chọn kỹ thuật

| Vấn đề | Lựa chọn | Lý do |
|---|---|---|
| HTTP | `net/http` (ServeMux của Go 1.22+), oapi-codegen v2 chế độ strict server | Theo CR-MOD-M1 mục 2; mã kiểu dữ liệu sinh từ hợp đồng |
| Docker | `github.com/moby/moby/client` | Thư viện chính thức của Docker Engine API |
| WebSocket | `github.com/coder/websocket` | Theo CR-MOD-M1 mục 2 |
| CSDL | `pgx/v5`, `golang-migrate` với migration nhúng | Theo CR-MOD-M1 mục 2 |
| Sự kiện | `amqp091-go`; bảng `crucible_outbox`, `crucible_event_sequence` cùng cấu trúc với thư viện Java | Định dạng CloudEvents và thuộc tính message giống `OutboxRelay` của Java để bên nhận M2, M3 dùng chung một cách đọc |
| JWT | `go-oidc/v3`, audience `runtime-api` | Realm phát triển đã có client `runtime-api` và mapper audience |
| Kiểm thử tích hợp | Testcontainers for Go (PostgreSQL, RabbitMQ) và Docker Engine thật | Giống M2, M3 |

## 3. Trình điều khiển

```go
type Driver interface {
    Provision(ctx context.Context, s Spec) (Handle, error)   // tạo container, ghi tệp, chạy setup
    Exec(ctx context.Context, h Handle, e ExecSpec) (ExecResult, error)
    Attach(ctx context.Context, h Handle, tty TTY) (Stream, error) // exec bash -l có TTY
    Destroy(ctx context.Context, h Handle) error
    Events(ctx context.Context) <-chan Exit                  // container dừng ngoài ý muốn
    List(ctx context.Context) ([]Handle, error)              // nhận lại khi khởi động
}
```

Container phiên ở chế độ `docker`:

| Thuộc tính | Giá trị |
|---|---|
| Lệnh | Theo image (`tini -- sleep infinity`) |
| Giới hạn | `Memory`, `MemorySwap` bằng nhau theo mức tài nguyên; `NanoCPUs`; `PidsLimit` 512 |
| Bảo mật | `CapDrop` mặc định của Docker, `CapAdd` theo BR-M1.26, `no-new-privileges` |
| Công cụ | Volume `crucible-tools-<phiên bản>` tạo một lần từ image `crucible-tools`, gắn `/opt/crucible/bin` chỉ đọc |
| Mạng | `none` hoặc `bridge` (spec mục 4) |
| Nhãn | `crucible.dev/session`, `crucible.dev/managed-by=runtime-api` |
| Khởi động lại | `RestartPolicy: no` (BR-M1.17) |

Tệp của bài lab được ghi bằng `CopyToContainer` (luồng tar), setup chạy bằng exec `/bin/sh -s` qua stdin.

## 4. Dữ liệu

Migration đầu tiên tạo `sessions` (các cột ở CR-MOD-M1 mục 7 cần cho phạm vi này, thêm `expires_at`,
`extension_minutes`, `last_activity_at`, `container_id`, `environment` jsonb), `terminal_tickets`, `idempotency_keys`,
`crucible_outbox`, `crucible_event_sequence`. Chỉ mục duy nhất có điều kiện cho BR-M1.1:

```sql
CREATE UNIQUE INDEX sessions_one_active ON sessions (owner_user_id, priority)
WHERE phase IN ('PENDING','PROVISIONING','READY','EXPIRING') AND priority <> 'LAB_CI';
```

Mỗi lần đổi pha cập nhật `sessions` và ghi outbox trong cùng transaction (CR-MOD-M1 mục 5.4).

## 5. Thay đổi ngoài `services/runtime`

| Tệp | Thay đổi |
|---|---|
| `contracts/runtime-api.yaml`, CR-DOC-03 mục 3.4 | Mã lỗi 501 `NOT_IMPLEMENTED` cho operation chưa có ở chế độ đang chạy |
| `dev/keycloak/realm-crucible.json` | Kiểm tra secret client `runtime-api`, mapper audience cho token của người dùng |
| `.github/workflows/ci.yml` | Job `runtime`: `go vet`, `staticcheck`, `go test` (có Docker trên runner), build hai image |
| `services/README.md`, CONTRIBUTING mục 5 | Cách chạy chế độ `docker` cho M2, M3 |

## 6. Kiểm thử

| Mức | Nội dung |
|---|---|
| Đơn vị | Máy trạng thái phiên, tính hạn và cảnh báo, cắt đầu ra 4 KB, chuẩn hóa lỗi |
| Tích hợp (Testcontainers và Docker thật) | D2–D11 theo luồng HTTP thật; nhận sự kiện từ RabbitMQ và đối chiếu lược đồ |
| Hợp đồng | Schemathesis trên dịch vụ đang chạy (D12) |
