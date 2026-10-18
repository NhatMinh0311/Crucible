# Đặc tả: `runtime-api` chế độ `docker`

| Thuộc tính | Giá trị |
|---|---|
| Mã | 001-runtime-docker-mode |
| Module | M1 |
| Mốc | MA |
| Trạng thái | Đã duyệt (18/10/2026) |
| Tài liệu liên quan | CR-DOC-03 mục 2.3, 3, 5, 6, 7; CR-MOD-M1 mục 2, 4, 6; ADR-0011; backlog: "runtime-api chế độ docker", "Image nền bài lab ubuntu-nginx và image bộ công cụ kiểm tra" |

## 1. Mục tiêu

Thành viên phát triển M2, M3 chạy được `runtime-api` trên máy của mình (Linux, macOS, Windows với WSL2) mà không cần
Kubernetes hay gVisor. Phiên là một container Docker trên máy; API và sự kiện giống bản chạy trên cụm, để mã của M2, M3
viết một lần và chạy được với cả hai.

Với chế độ này, M2 kiểm thử được luồng thật: dựng môi trường từ đặc tả bài lab, terminal tương tác, lệnh kiểm tra chạy
bằng bộ công cụ chỉ đọc, script khôi phục, sự kiện `runtime.session.*`.

## 2. Phạm vi

| Thuộc phạm vi | Ngoài phạm vi (trả 501 `NOT_IMPLEMENTED` ở chế độ `docker`) |
|---|---|
| `createSession`, `getSession`, `extendSession`, `terminateSession` | Đặt trước, sức chứa, gợi ý khung giờ |
| `execInSession` (`CHECK`, `RESTORE`) | Tài khoản hạn mức, sổ giờ máy |
| `createTerminalTicket`, `createViewTicket`, terminal WebSocket (tương tác và xem) | Lần phát hiện dùng CPU bất thường |
| Sự kiện `runtime.session.ready`, `.expiring`, `.terminated`, `.failed`, `runtime.terminal.connected` | Tham số hệ thống `/v1/settings` |
| Hết hạn, không hoạt động, gia hạn | Bản ghi phiên, sự kiện `runtime.terminal.input-line`, `runtime.recording.available` |
| Image bộ công cụ kiểm tra `/opt/crucible/bin` | Nhóm ấm, hàng đợi ưu tiên, `runtime.session.lost`, `resource-alert` |
| Image nền bài lab `ubuntu-nginx` | |

M2 dùng mock Prism cho các operation ngoài phạm vi khi cần. Các phần ngoài phạm vi được hiện thực ở chế độ `kubernetes`
hoặc ở các đặc tả sau.

## 3. Tiêu chí chấp nhận

| Mã | Tiêu chí |
|---|---|
| D1 | Chạy `docker compose -f dev/compose.yaml up -d --wait` rồi `go run ./cmd/runtime-api` trong `services/runtime` với `RUNTIME_DRIVER=docker`; dịch vụ lắng nghe cổng 8081 và tự tạo bảng trong CSDL `runtime` |
| D2 | `content-service` gọi `createSession` bằng token client credentials của Keycloak phát triển với bài `linux-user-nginx`; nhận 202 pha `PENDING`; trong 15 giây (image đã có trên máy) phiên chuyển `READY` và sự kiện `runtime.session.ready` tới exchange `crucible.events` với `labels` nguyên vẹn |
| D3 | Container phiên có giới hạn bộ nhớ, CPU theo mức tài nguyên (`small` 256 MiB, 1 CPU; `medium` 512 MiB, 2 CPU), giới hạn 512 tiến trình, `/opt/crucible/bin` chỉ đọc, không capability bổ sung ngoài `SESSION_ALLOWED_CAPABILITIES` (BR-M1.26) |
| D4 | `execInSession` với `CHECK` chạy bằng `/opt/crucible/bin/sh`, môi trường biến chỉ gồm `PATH=/opt/crucible/bin`, `LANG=C.UTF-8`; kết quả không đổi khi học viên thay `/usr/bin/curl` bằng script luôn thành công (AC-C2.4). Hết thời gian trả `exitCode = -1`, `timedOut = true`; `stdout`, `stderr` cắt ở 4096 byte |
| D5 | Hai yêu cầu exec đồng thời vào một phiên: một yêu cầu nhận 409 `EXEC_BUSY` (BR-M1.16) |
| D6 | Chủ phiên lấy vé, mở WebSocket theo CR-DOC-03 mục 6: gõ lệnh, đổi kích thước, Ctrl-C hoạt động; vé dùng lại hoặc quá 30 giây bị đóng với mã 4001 (BR-M1.13); kết nối tương tác thứ hai đóng kết nối thứ nhất với mã 4009 (BR-M1.14); kết nối xem chỉ nhận đầu ra |
| D7 | Hết `ttlMinutes` (cộng gia hạn) phiên chuyển `TERMINATED` lý do `EXPIRED`; trước 5 phút phát `runtime.session.expiring` và khung `expiring` trên terminal. Phiên `PRACTICE` không có dữ liệu nhập trong `idleMinutes` kết thúc lý do `IDLE`, có cảnh báo trước 2 phút (BR-M1.9) |
| D8 | Tiến trình trong phiên dùng vượt bộ nhớ: phiên chuyển `TERMINATED` lý do `OUT_OF_MEMORY`, `resourceUsage.oomKilled = true`; container dừng vì lý do khác: `ENVIRONMENT_FAILED` (BR-M1.17) |
| D9 | Một học viên tạo phiên `PRACTICE` thứ hai khi phiên đầu chưa kết thúc nhận 409 `SESSION_LIMIT` (BR-M1.1); gửi lại cùng `Idempotency-Key` và cùng nội dung nhận lại phản hồi cũ, khác nội dung nhận 422 `IDEMPOTENCY_KEY_REUSED` |
| D10 | Setup lỗi hoặc quá `SETUP_TIMEOUT`: phiên chuyển `FAILED` lý do `PROVISION_FAILED`, `failureDetail` có dòng cuối của đầu ra; container bị xóa |
| D11 | Khởi động lại `runtime-api`: phiên đang chạy được nhận lại theo nhãn container; container không còn phiên chưa kết thúc tương ứng bị xóa |
| D12 | Kiểm thử hợp đồng Schemathesis với các operation thuộc phạm vi không báo lỗi lệch lược đồ |

## 4. Trường hợp biên

| Trường hợp | Xử lý |
|---|---|
| Image chưa có trên máy | Kéo image khi dựng phiên; thời gian kéo tính vào thời hạn dựng (mặc định 5 phút); quá hạn là `PROVISION_FAILED` |
| `network: mirror` | Chế độ `docker` không có máy chủ gói nội bộ; container dùng mạng `bridge` mặc định của Docker (có Internet) |
| `network: none` | Container không có mạng |
| Docker Engine không chạy khi khởi động | Dịch vụ thoát với thông báo rõ ràng |
| Gọi `terminateSession` khi phiên đang dựng | Dừng dựng, xóa container, phiên `TERMINATED` |
| Lệnh `RESTORE` | Chạy bằng `/bin/sh -c` với `PATH` chuẩn của môi trường, thời hạn tối đa 60 giây |

## 5. Điểm đã làm rõ

| Câu hỏi | Quyết định |
|---|---|
| Phạm vi lần giao đầu | Như mục 2; đặt trước, hạn mức, tham số hệ thống trả 501 |
| Xác thực khi phát triển | Dùng token thật của Keycloak phát triển như hai dịch vụ còn lại |
| `network: mirror` ở chế độ `docker` | Cho ra Internet qua mạng `bridge` |
