# ADR-0008: Triển khai bằng GitOps với Argo CD và Argo CD Image Updater

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Cụm chạy trên máy ảo sau NAT trong máy tính cá nhân; GitHub Actions không kết nối vào cụm được. Nhánh `main` được bảo vệ,
không cho phép đẩy trực tiếp kể cả từ CI. Cụm phải dựng lại được từ kho mã (Q12) và cập nhật không gián đoạn phiên (Q9).

## Quyết định
- Argo CD trong cụm kéo cấu hình từ `deploy/` trên nhánh `main` theo mẫu app-of-apps; mỗi dịch vụ một Helm chart.
- CI chỉ build và đẩy image lên GHCR với tag `main-<sha>`. Argo CD Image Updater phát hiện tag mới và cập nhật tham số
  image của Application (phương thức ghi ngược `argocd`).
- Migration CSDL chạy bằng Job PreSync; kiểm tra sau triển khai bằng Job PostSync.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| CI chạy `kubectl apply` hoặc `helm upgrade` | Cần mở cụm ra Internet hoặc runner tự quản trong cụm; trạng thái cụm không khớp kho mã |
| CI tạo commit cập nhật tag vào `main` | Bị ruleset chặn; tạo pull request tự động thì mỗi lần triển khai cần thêm một lượt duyệt |
| Flux CD | Tương đương về chức năng; Argo CD có giao diện trực quan cho thành viên và cho demo |

## Hệ quả
- (+) Cụm không cần mở cổng vào; mọi thay đổi cấu hình đi qua pull request.
- (+) Lịch sử đồng bộ và quay lui có sẵn trong Argo CD.
- (−) Tag image đang chạy nằm trong Argo CD, không nằm trong kho mã; tái lập chính xác một thời điểm cần đọc lịch sử Argo CD.
- (−) Argo CD và Image Updater chiếm khoảng 0,7 GB bộ nhớ.
