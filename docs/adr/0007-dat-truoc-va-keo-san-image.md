# ADR-0007: Đặt trước sức chứa theo kỳ thi trong CSDL, kéo sẵn image, không giữ môi trường ấm

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận (chờ kết quả thử nghiệm S2) |
| Ngày | 26/10/2026 |

## Bối cảnh
Kỳ thi khiến hàng chục yêu cầu đến trong cùng một phút (Q1, Q2). Máy phiên có sức chứa cố định. Giảng viên cần biết khi
lên lịch rằng kỳ thi có đủ chỗ. Kéo image lớn trong lúc học viên chờ có thể mất hàng chục giây.

## Quyết định
- Đặt trước là bản ghi trong bảng `reservations` của CSDL `runtime`. Kiểm tra sức chứa và ghi đặt trước chạy trong một
  transaction có khóa dòng trên các đặt trước chồng lấn (CR-MOD-M1 mục 5.1).
- Lên lịch kỳ thi bị từ chối nếu tổng đặt trước trong khung giờ vượt sức chứa sau khi trừ phần tối thiểu cho luyện tập
  (BR-M1.7). Chỗ đặt trước chỉ cấp cho phiên `EXAM` mang `reservationRef` tương ứng (BR-M1.6).
- Image của bài lab được kéo sẵn lên máy phiên khi bài lab được xuất bản và khi đặt trước được tạo.
- Mỗi phiên dựng môi trường mới khi được cấp chỗ; không giữ sẵn môi trường ấm.
- Hàng đợi khi thiếu chỗ theo ưu tiên `EXAM` > `PRACTICE` > `LAB_CI`, không dừng phiên đang chạy để nhường chỗ.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Nhóm môi trường ấm theo image, nâng theo lịch đặt trước | Thêm một vòng điều khiển và trạng thái phải đối soát; chỉ cần khi thời gian khởi động sandbox vượt mục tiêu Q1 (đo ở S2) |
| Không đặt trước, cấp theo thứ tự đến | Kỳ thi có thể thiếu chỗ khi lớp khác đang luyện tập |
| Dừng phiên luyện tập để nhường chỗ cho kỳ thi | Mất công sức của học viên; đặt trước tránh được tình huống này |

## Hệ quả
- (+) Giảng viên biết ngay khi lên lịch nếu máy phiên không đủ chỗ.
- (+) Không có tài nguyên bị chiếm khi không có ai dùng.
- (−) Thời gian sẵn sàng gồm khởi động sandbox và setup; script setup của bài lab phải ngắn. Nếu S2 cho thấy không đạt
  Q1, bổ sung nhóm môi trường ấm cố định cho các image của kỳ thi.
