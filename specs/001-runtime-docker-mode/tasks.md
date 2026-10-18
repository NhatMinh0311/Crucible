# Danh sách việc: `runtime-api` chế độ `docker`

Mỗi việc là một commit, kiểm chứng được bằng lệnh ghi ở cột cuối.

| # | Việc | Kiểm chứng |
|---|---|---|
| T1 | Image `crucible-tools`: BusyBox, curl, jq liên kết tĩnh, ghim phiên bản và mã băm | `docker run --rm crucible-tools /opt/crucible/bin/jq --version`; `file` báo statically linked |
| T2 | Image `lab-ubuntu-nginx` với tini làm PID 1 | Container chạy `sleep infinity`; `nginx -t` đạt |
| T3 | Khung Go module, cấu hình, sinh mã từ hợp đồng, `GET /healthz` | `go build ./...`, `go vet ./...` |
| T4 | Hợp đồng: mã lỗi 501 `NOT_IMPLEMENTED` | Redocly lint |
| T5 | CSDL: migration, tầng lưu phiên, vé, khóa lặp lại an toàn, outbox | Kiểm thử tích hợp với PostgreSQL |
| T6 | Chuyển tiếp outbox lên RabbitMQ | Kiểm thử tích hợp: sự kiện tới hàng đợi với đúng thuộc tính message |
| T7 | Xác thực JWT, scope, chủ phiên; lỗi RFC 9457 | Kiểm thử đơn vị với JWKS giả |
| T8 | Trình điều khiển Docker: dựng, exec, xóa, theo dõi container dừng, nhận lại | Kiểm thử tích hợp D3, D4, D8, D10, D11 |
| T9 | Dịch vụ phiên và handler: tạo, đọc, gia hạn, kết thúc, exec; hết hạn, không hoạt động | Kiểm thử tích hợp D2, D5, D7, D9 |
| T10 | Terminal WebSocket: vé, tương tác, xem, một kết nối tương tác, khung trạng thái | Kiểm thử tích hợp D6 |
| T11 | CI job `runtime`; hướng dẫn chạy cho M2, M3 | CI xanh; thành viên làm theo hướng dẫn chạy được D2 |
| T12 | Schemathesis | D12 |
