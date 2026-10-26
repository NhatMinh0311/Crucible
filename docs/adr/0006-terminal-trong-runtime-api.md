# ADR-0006: Terminal trong runtime-api qua Docker exec, vé dùng một lần và ghi asciicast

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận (chờ kết quả thử nghiệm S3) |
| Ngày | 26/10/2026 |

## Bối cảnh
Học viên cần terminal trong trình duyệt; giảng viên cần xem trực tiếp và xem lại; M3 cần dòng lệnh đã nhập để phân tích
gian lận. Môi trường học viên không được mở cổng nào ra ngoài (ADR-0003). Kết nối WebSocket sống lâu; cập nhật dịch vụ
không được làm mất môi trường của học viên (Q9).

## Quyết định
- `runtime-api` nhận WebSocket từ trình duyệt và mở exec có TTY tới container phiên qua Docker Engine API; container phiên
  không mở cổng.
- Xác thực bằng vé dùng một lần do `runtime-api` cấp; vé lưu trong CSDL `runtime`, tiêu thụ nguyên tử.
- Người xem trực tiếp nhận đầu ra của kết nối tương tác trong cùng tiến trình.
- Ghi luồng ra asciicast v2 vào SeaweedFS theo đoạn; tái dựng dòng lệnh từ luồng nhập cho phiên có `recordInput`.
- Mỗi phiên tối đa một kết nối tương tác; kết nối xem chỉ nhận dữ liệu.
- Khi dừng để cập nhật, `runtime-api` gửi mã đóng 4010; giao diện lấy vé mới và kết nối lại.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| ttyd hoặc máy chủ SSH chạy trong mỗi container phiên | Phải mở cổng vào môi trường; học viên có root nên kiểm soát được tiến trình phục vụ terminal và bản ghi |
| Trình duyệt gọi Docker Engine API trực tiếp | Lộ quyền điều khiển máy phiên ra người dùng |
| Dịch vụ gateway riêng, nhiều bản, phát cho người xem qua RabbitMQ | Cập nhật API không làm đứt terminal, nhưng thêm một dịch vụ và cơ chế phân phối giữa các bản; chưa cần ở quy mô một lớp |

## Hệ quả
- (+) Bản ghi và dòng lệnh được thu ở phía nền tảng, học viên không can thiệp được.
- (+) Container phiên không có cổng mở; một dịch vụ duy nhất giữ trạng thái kết nối nên quy tắc một kết nối tương tác và
  người xem xử lý trong bộ nhớ.
- (−) Toàn bộ lưu lượng terminal đi qua kết nối SSH tới Docker Engine của máy phiên; cần đo ở S3.
- (−) Cập nhật `runtime-api` làm đứt mọi terminal trong vài giây.
- (−) Tái dựng dòng lệnh từ luồng nhập không chính xác tuyệt đối.
