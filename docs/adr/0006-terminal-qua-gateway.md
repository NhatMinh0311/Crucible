# ADR-0006: Terminal qua runtime-gateway dùng Kubernetes exec, vé dùng một lần và ghi asciicast

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận (chờ kết quả thử nghiệm S3) |
| Ngày | 17/10/2026 |

## Bối cảnh
Học viên cần terminal trong trình duyệt; giảng viên cần xem trực tiếp và xem lại; M3 cần dòng lệnh đã nhập để phân tích
gian lận. Môi trường học viên không được mở cổng nào ra ngoài (ADR-0003). Kết nối WebSocket sống lâu, phải sống sót
qua cập nhật dịch vụ (Q9).

## Quyết định
- `runtime-gateway` là dịch vụ riêng của M1, nhận WebSocket từ trình duyệt và mở luồng exec tới container phiên qua
  Kubernetes API; pod phiên không mở cổng.
- Xác thực bằng vé dùng một lần do `runtime-api` cấp; vé lưu trong CSDL `runtime`, tiêu thụ nguyên tử.
- Gateway ghi luồng ra asciicast v2 vào MinIO theo đoạn; tái dựng dòng lệnh từ luồng nhập cho phiên có `recordInput`.
- Mỗi phiên tối đa một kết nối tương tác; kết nối xem chỉ nhận dữ liệu.
- Khi dừng để cập nhật, gateway gửi mã đóng 4010; giao diện lấy vé mới và kết nối lại.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| ttyd hoặc máy chủ SSH chạy trong mỗi pod phiên | Phải mở cổng vào pod phiên; học viên có root nên kiểm soát được tiến trình phục vụ terminal và bản ghi |
| Trình duyệt gọi Kubernetes API trực tiếp | Lộ API server ra người dùng |
| Gộp gateway vào `runtime-api` | Cập nhật API làm đứt mọi terminal; hai loại tải khác nhau cần mở rộng riêng |

## Hệ quả
- (+) Bản ghi và dòng lệnh được thu ở phía nền tảng, học viên không can thiệp được.
- (+) Pod phiên không có cổng mở.
- (−) Toàn bộ lưu lượng terminal đi qua API server và kubelet; cần đo ở S3 và kiểm thử sức chứa.
- (−) Tái dựng dòng lệnh từ luồng nhập không chính xác tuyệt đối.
