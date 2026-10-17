# ADR-0007: Nhóm môi trường ấm, đặt trước tài nguyên theo kỳ thi và kéo sẵn image

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận (chờ kết quả thử nghiệm S2) |
| Ngày | 17/10/2026 |

## Bối cảnh
Tạo pod mới với gVisor, kéo image và khởi động mất từ vài giây tới vài chục giây. Kỳ thi khiến hàng chục yêu cầu đến
trong cùng một phút (Q1, Q2). Cụm có sức chứa cố định, không tự thêm node được.

## Quyết định
- `WarmPool` giữ sẵn pod đã khởi động cho mỗi cặp (image, mức tài nguyên) phổ biến; pod ấm chưa có dữ liệu riêng của học viên.
  Khi cấp phiên, operator chọn một pod ấm, gắn nhãn sở hữu, ghi tệp và chạy setup của biến thể; nhóm được bổ sung bất đồng bộ.
- `CapacityReservation` giữ số chỗ cho một kỳ thi trong khung giờ. Từ 15 phút trước giờ bắt đầu, operator nâng nhóm ấm
  bằng số chỗ đặt trước của các bài lab trong kỳ thi. Chỗ đặt trước chỉ cấp cho phiên `EXAM` mang `reservationRef` tương ứng.
- Lên lịch kỳ thi bị từ chối nếu tổng đặt trước trong khung giờ vượt sức chứa sau khi trừ phần tối thiểu cho luyện tập.
- DaemonSet trên node phiên kéo sẵn mọi image bài lab đã xuất bản.
- Hàng đợi khi thiếu chỗ theo ưu tiên `EXAM` > `PRACTICE` > `LAB_CI`, không dừng phiên đang chạy để nhường chỗ.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Chỉ tạo pod khi có yêu cầu | Không đạt Q1, Q2 |
| Nhóm ấm cố định lớn | Chiếm tài nguyên cả khi không có kỳ thi |
| Dừng phiên luyện tập để nhường chỗ cho kỳ thi | Mất công sức của học viên; đặt trước tránh được tình huống này |

## Hệ quả
- (+) Thời gian sẵn sàng chỉ còn thời gian chạy setup.
- (+) Giảng viên biết ngay khi lên lịch nếu cụm không đủ chỗ.
- (−) Pod ấm chiếm tài nguyên khi không có ai dùng.
- (−) Setup của bài lab chạy trong thời gian chờ của học viên, nên script setup phải ngắn.
