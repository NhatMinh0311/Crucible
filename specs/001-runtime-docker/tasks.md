# Danh sách việc: `runtime-api` chạy phiên bằng Docker

Mỗi việc là một commit, kiểm chứng được bằng lệnh ghi ở cột cuối.

| # | Việc | Kiểm chứng | Trạng thái |
|---|---|---|---|
| T1 | Image `crucible-tools`: BusyBox, curl, jq liên kết tĩnh, ghim phiên bản và mã băm | `docker run --rm crucible-tools /opt/crucible/bin/jq --version`; `file` báo statically linked | Xong |
| T2 | Image `lab-ubuntu-nginx` với tini làm PID 1 | Container chạy `sleep infinity`; `nginx -t` đạt | Xong |
| T3 | Khung Go module, cấu hình, sinh mã từ hợp đồng, `GET /healthz` | `go build ./...`, `go vet ./...` | Xong |
| T4 | Hợp đồng: mã lỗi 501 `NOT_IMPLEMENTED` | Redocly lint | Xong |
| T5 | CSDL: migration, tầng lưu phiên, vé, khóa lặp lại an toàn, outbox | Kiểm thử tích hợp với PostgreSQL | Xong |
| T6 | Chuyển tiếp outbox lên RabbitMQ | Kiểm thử tích hợp: sự kiện tới hàng đợi với đúng thuộc tính message | Xong |
| T7 | Xác thực JWT, scope, chủ phiên; lỗi RFC 9457 | Kiểm thử đơn vị với JWKS giả | Xong |
| T8 | Trình điều khiển Docker: dựng, exec, xóa, theo dõi container dừng, nhận lại | Kiểm thử tích hợp D3, D4, D8, D10, D11 | Xong |
| T9 | Dịch vụ phiên và handler: tạo, đọc, gia hạn, kết thúc, exec; hết hạn, không hoạt động | Kiểm thử tích hợp D2, D5, D7, D9 | Xong |
| T10 | Terminal WebSocket: vé, tương tác, xem, một kết nối tương tác, khung trạng thái | Kiểm thử tích hợp D6 | Xong |
| T11 | CI job `runtime`; hướng dẫn chạy cho M2, M3 | CI xanh; thành viên làm theo hướng dẫn chạy được D2 | Xong |
| T12 | Schemathesis | D12 | Xong |

Kết quả: toàn bộ kiểm thử Go đạt với `-race` (Testcontainers PostgreSQL, RabbitMQ; Docker Engine thật); Schemathesis 655 ca
đạt; chạy thử với Keycloak phát triển và profile `runtime` của `dev/compose.yaml`.
