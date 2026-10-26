# ADR-0008: Triển khai bằng Docker Compose và script triển khai chạy tay

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 26/10/2026 |

## Bối cảnh
Môi trường chuẩn chạy trên máy ảo sau NAT trong máy tính cá nhân; GitHub Actions không kết nối vào được. Nhánh `main`
được bảo vệ, không cho phép đẩy trực tiếp kể cả từ CI. Môi trường phải dựng lại được từ kho mã (Q12). Nhóm có một người
phụ trách vận hành.

## Quyết định
- Thành phần nền tảng và dịch vụ của ba module chạy bằng Docker Compose trên `cr-server`, cấu hình tại `deploy/`.
- CI build và đẩy image lên GHCR với tag `main-<sha>` khi merge vào `main`.
- Trưởng nhóm chạy `deploy/deploy.sh <tag>` trên `cr-server`: kéo image theo tag, `docker compose up -d`, kiểm tra
  `healthz` của từng dịch vụ; tag đã triển khai được ghi vào nhật ký triển khai trên máy chủ.
- Bí mật nằm trong tệp `.env` trên `cr-server` (quyền 600); kho mã chỉ chứa `deploy/.env.example`.
- Migration CSDL do mỗi dịch vụ chạy khi khởi động (Flyway, golang-migrate).

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| GitOps với Argo CD trên cụm Kubernetes | Cần cụm Kubernetes (ADR-0002) |
| CI tự triển khai qua SSH trong mạng Tailscale | Phải giữ khóa SSH và khóa Tailscale trong secret của kho mã; thêm cấu hình CI trong khi số lần triển khai mỗi tuần thấp |
| Dựng và chạy trực tiếp từ mã nguồn trên máy chủ | Không tái lập được phiên bản đang chạy; build trên máy chủ chiếm tài nguyên |

## Hệ quả
- (+) Một công cụ (Docker Compose) cho cả máy phát triển và môi trường chuẩn.
- (+) Phiên bản đang chạy xác định bằng một tag; quay lui bằng cách chạy lại script với tag cũ.
- (−) Triển khai phụ thuộc người chạy script; không có đồng bộ tự động giữa kho mã và môi trường.
- (−) `docker compose up -d` khởi động lại dịch vụ thay đổi; dịch vụ một bản gián đoạn vài giây.
