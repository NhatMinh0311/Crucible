# ADR-0012: Quan sát bằng Prometheus, Grafana và log JSON có trace_id

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận (chờ kết quả thử nghiệm S5 về bộ nhớ) |
| Ngày | 26/10/2026 |

## Bối cảnh
Một yêu cầu tạo phiên đi qua ba dịch vụ của ba thành viên viết bằng hai ngôn ngữ và qua broker. Yêu cầu chất lượng phải
được đo bằng số liệu (Q1–Q9). Thành viên cần tự chẩn đoán lỗi của module mình trên môi trường chuẩn. Bộ nhớ của
`cr-server` có hạn.

## Quyết định
- Mỗi dịch vụ xuất metric dạng Prometheus (`/metrics` của Go, Spring Boot Actuator) theo danh mục tại CR-DOC-04 mục 7.2.
- Prometheus và Grafana chạy trong Docker Compose trên `cr-server`; Grafana đăng nhập qua Keycloak; dashboard và luật
  cảnh báo quản lý dạng mã trong `deploy/`.
- Dịch vụ truyền `traceparent` ở REST và thuộc tính `traceparent` của CloudEvents; log dạng JSON ghi `trace_id`, xem bằng
  `docker compose logs` và lọc theo `trace_id`.
- Không thu trace và log tập trung ở phiên bản này.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| OpenTelemetry, Tempo, Loki | Truy vết qua nhiều module bằng giao diện, nhưng thêm ba thành phần, khoảng 1 GB bộ nhớ và kiến thức phải học |
| Elasticsearch, Kibana | Tốn bộ nhớ vượt khả năng của `cr-server` |
| Dịch vụ quan sát trên cloud | Phụ thuộc bên ngoài, giới hạn gói miễn phí |

## Hệ quả
- (+) Số đo yêu cầu chất lượng lấy trực tiếp từ dashboard.
- (+) Lọc log của mọi dịch vụ theo một `trace_id` vẫn cho phép lần theo một yêu cầu.
- (−) Không có biểu đồ thời gian theo từng bước của một yêu cầu; chẩn đoán độ trễ qua nhiều module làm bằng tay.
