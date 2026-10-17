# ADR-0005: Sự kiện giữa module qua RabbitMQ với transactional outbox và inbox

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Các module cần biết sự việc xảy ra ở module khác (phiên sẵn sàng, phiên mất, bài được chấm) mà không gọi đồng bộ vòng
tròn. Sự kiện không được mất khi broker gián đoạn và không được phát khi transaction nghiệp vụ bị hủy (Q8). Tần suất
sự kiện thấp (vài chục mỗi giây khi cao điểm).

## Quyết định
- RabbitMQ một node, một topic exchange `crucible.events`, quorum queue cho mỗi bên nhận, dead-letter queue.
- Bên phát ghi sự kiện vào bảng `outbox` trong cùng transaction; tiến trình chuyển tiếp gửi với publisher confirm.
- Bên nhận loại trùng bằng bảng `inbox`, bỏ qua sự kiện cũ theo `sequence`.
- Định dạng CloudEvents 1.0; danh mục bằng AsyncAPI (CR-DOC-03 mục 5).

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Apache Kafka | Tốn bộ nhớ và công vận hành không tương xứng tần suất sự kiện; lợi thế phát lại lịch sử không cần thiết ở quy mô này |
| NATS JetStream | Tích hợp với Spring kém hơn Spring AMQP |
| Chỉ gọi REST đồng bộ | Tạo phụ thuộc vòng (M1 phải gọi M2 khi phiên sẵn sàng); mất thông báo khi bên nhận không chạy |
| Gửi lên broker trực tiếp không qua outbox | Mất sự kiện khi broker lỗi sau khi transaction đã xác nhận |

## Hệ quả
- (+) Không mất sự kiện, không phụ thuộc vòng.
- (+) Tiêu tốn khoảng 0,3 GB bộ nhớ.
- (−) Nhất quán sau một khoảng trễ; giao diện phải xử lý trạng thái trung gian.
- (−) Mỗi dịch vụ cần bảng `outbox`, `inbox` và tiến trình chuyển tiếp.
