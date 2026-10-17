# ADR-0012: Quan sát bằng OpenTelemetry, Prometheus, Loki, Tempo và Grafana

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận (chờ kết quả thử nghiệm S5 về bộ nhớ) |
| Ngày | 17/10/2026 |

## Bối cảnh
Một yêu cầu tạo phiên đi qua ba dịch vụ của ba thành viên viết bằng hai ngôn ngữ, qua broker và operator. Yêu cầu chất
lượng phải được đo bằng số liệu (Q1–Q10). Thành viên cần tự chẩn đoán lỗi của module mình trên cụm chuẩn.

## Quyết định
- Dịch vụ Java dùng OpenTelemetry Java agent; dịch vụ Go dùng OpenTelemetry SDK; trace truyền qua `traceparent` ở REST
  và thuộc tính `traceparent` của CloudEvents.
- Metric nghiệp vụ theo danh mục tại CR-DOC-04 mục 7.2, xuất qua OTLP hoặc endpoint Prometheus.
- Log dạng JSON có `trace_id`, thu bằng Grafana Alloy vào Loki.
- kube-prometheus-stack, Loki và Tempo chế độ một tiến trình, Grafana đăng nhập qua Keycloak; dashboard và luật cảnh
  báo quản lý dạng mã trong `deploy/`.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Elasticsearch, Kibana | Tốn bộ nhớ vượt khả năng của node nền tảng |
| Chỉ metric, không trace | Không chẩn đoán được độ trễ qua nhiều module (Q10) |
| Dịch vụ quan sát trên cloud | Phụ thuộc bên ngoài, giới hạn gói miễn phí |

## Hệ quả
- (+) Truy vết một yêu cầu qua mọi module bằng một trace ID.
- (+) Số đo yêu cầu chất lượng lấy trực tiếp từ dashboard.
- (−) Khoảng 2,3 GB bộ nhớ cho giám sát; thời gian lưu giữ ngắn.
