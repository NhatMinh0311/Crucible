# ADR-0002: Biểu diễn phiên bằng tài nguyên Kubernetes tùy biến và quản lý bằng operator viết bằng Go

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Phiên là một tài nguyên hạ tầng sống từ vài chục phút tới vài giờ. Hệ thống phải bảo đảm không có môi trường mồ côi,
không nhân đôi phiên khi thành phần khởi động lại, phát hiện mất node và thu hồi đúng hạn (Q4, Q5). Trạng thái thực tế
của pod có thể thay đổi ngoài ý muốn của hệ thống (node mất, pod bị xóa).

## Quyết định
- Định nghĩa các CustomResourceDefinition `LabSession`, `WarmPool`, `CapacityReservation` (CR-DOC-02 mục 4.2).
- `runtime-operator` viết bằng Go với controller-runtime, đối soát liên tục trạng thái mong muốn trong tài nguyên với
  pod thực tế; pod gắn `ownerReference` tới `LabSession` để bị thu hồi theo.
- `runtime-api` chỉ ghi tài nguyên và đọc trạng thái; không trực tiếp tạo pod.
- CSDL `runtime` lưu lịch sử phiên, hàng đợi, sổ hạn mức, vé; tài nguyên Kubernetes là nguồn sự thật cho trạng thái sống.
- Operator chạy hai bản với bầu chọn leader.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Dịch vụ Spring Boot gọi Kubernetes API tạo pod trực tiếp (mệnh lệnh) | Phải tự xử lý mọi trường hợp lệch trạng thái; khởi động lại giữa chừng dễ để lại pod mồ côi |
| Lưu trạng thái sống trong CSDL, đồng bộ sang Kubernetes | Hai nguồn sự thật; khó bảo đảm nhất quán khi node mất |
| Dùng dự án có sẵn (Kubernetes Jobs, các nền tảng sandbox mã nguồn mở) | Không có khái niệm nhóm ấm, đặt trước, hạn mức theo nghiệp vụ; mất phần học tập chính của module |

## Hệ quả
- (+) Khởi động lại `runtime-api` hoặc `runtime-operator` không ảnh hưởng phiên đang chạy.
- (+) Mô hình đối soát là cách làm chuẩn của hệ sinh thái Kubernetes; kiểm thử được bằng envtest.
- (−) M1 dùng ngôn ngữ khác M2, M3; review chéo giữa thành viên khó hơn.
- (−) Tải lên API server tăng theo số phiên; cần giới hạn tần suất đối soát.
