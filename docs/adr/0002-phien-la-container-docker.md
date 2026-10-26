# ADR-0002: Phiên là container Docker trên máy phiên riêng, điều khiển bằng một dịch vụ Go

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 26/10/2026 |

## Bối cảnh
Phiên là một môi trường Linux sống từ vài chục phút tới vài giờ. Hệ thống phải bảo đảm không có môi trường mồ côi,
không nhân đôi phiên khi thành phần khởi động lại, phát hiện máy chạy phiên mất kết nối và thu hồi đúng hạn (Q4, Q5).
Quy mô mục tiêu là một lớp học (khoảng 40 phiên đồng thời) trên hạ tầng của một máy tính cá nhân. Người phụ trách M1
cần một thiết kế vận hành và kiểm chứng được trong 11 tuần.

## Quyết định
- Mỗi phiên là một container Docker trên máy phiên `cr-lab`, chạy bằng gVisor (ADR-0003).
- `runtime-api` (Go) điều khiển Docker Engine của máy phiên qua Docker Engine API với `DOCKER_HOST=ssh://`; container và
  mạng của phiên mang nhãn để nhận lại sau khi khởi động lại.
- CSDL `runtime` là nguồn sự thật về phiên; `runtime-api` đối soát với Docker khi khởi động và định kỳ, xóa tài nguyên
  mồ côi (CR-MOD-M1 mục 3.3).
- `runtime-api` chạy một bản, gồm API, vòng phiên và terminal (ADR-0006).
- Tầng phiên làm việc qua giao diện `driver.Driver`; trình điều khiển cho nền tảng điều phối khác bổ sung được mà không
  đổi tầng phiên và hợp đồng.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Cụm Kubernetes (k3s) với tài nguyên tùy biến và operator | Đúng hướng cho nhiều máy chủ, nhưng thêm vận hành cụm, operator, Helm, GitOps; vượt khối lượng của một người trong 11 tuần ở quy mô một lớp |
| Kubernetes, dịch vụ tạo pod trực tiếp không dùng operator | Vẫn phải vận hành cụm; không tận dụng được lợi ích chính của Kubernetes là đối soát |
| Docker Engine cùng máy với nền tảng | Thoát khỏi sandbox chạm tới CSDL và Keycloak; phiên tranh bộ nhớ với nền tảng |

## Hệ quả
- (+) Ít thành phần: một dịch vụ, một Docker Engine; trình điều khiển `docker` dùng chung cho máy phát triển của M2, M3.
- (+) Khởi động lại `runtime-api` không ảnh hưởng môi trường đang chạy; terminal kết nối lại.
- (−) Sức chứa giới hạn trong một máy phiên; mất máy phiên làm mất mọi phiên (BR-M1.11).
- (−) Một bản `runtime-api`: cập nhật dịch vụ làm đứt terminal trong vài giây (Q9).
- (−) Đối soát và dọn mồ côi do mã của nhóm tự viết thay vì có sẵn trong nền tảng điều phối.
