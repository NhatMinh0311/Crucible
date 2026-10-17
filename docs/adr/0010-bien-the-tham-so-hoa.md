# ADR-0010: Đề bài do giảng viên soạn có tham số, biến thể sinh tất định bằng HMAC

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Bài thi giống nhau cho mọi học viên dễ bị chép kết quả. Đề bài do hệ thống tự sinh hoàn toàn không bảo đảm giải được và
khó kiểm soát chất lượng. Biến thể của một học viên phải ổn định qua nhiều lần tiếp tục và không đoán được từ biến thể
của học viên khác.

## Quyết định
- Giảng viên soạn đặc tả bài lab với tham số và miền giá trị; có thể khai báo danh sách lỗi cài sẵn (CR-DOC-03 mục 4).
- Biến thể sinh tất định: seed = HMAC-SHA256 với khóa bí mật của M2 trên (ngữ cảnh, học viên, bài lab, số thứ tự lượt làm bài); giá trị từng tham
  số và lỗi cài sẵn suy ra từ seed theo thuật toán tại CR-DOC-03 mục 4.4.
- Phiên bản chỉ được xuất bản khi lời giải mẫu đạt với 5 seed và môi trường mới dựng không đạt (CR-DOC-03 mục 4.5).
- Giá trị biến thể được cung cấp cho M3 sau kỳ thi để phát hiện học viên dùng giá trị thuộc biến thể của người khác.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Một đề chung | Không chống được chép kết quả |
| Sinh đề bằng mô hình AI | Không bảo đảm giải được, lệnh kiểm tra có thể sai; vẫn cần kiểm thử như trên |
| Sinh ngẫu nhiên và lưu lại | Cần lưu và đồng bộ trạng thái; HMAC cho cùng kết quả mà không cần lưu |

## Hệ quả
- (+) Mỗi học viên một biến thể; lộ biến thể là tín hiệu gian lận mạnh.
- (+) Bài lab được kiểm thử tự động trước khi giao.
- (−) Đổi khóa ký làm đổi toàn bộ biến thể; khóa chỉ được xoay vòng ngoài thời gian thi và lượt làm bài lưu lại biến thể đã dùng.
- (−) Giảng viên phải viết lệnh kiểm tra phụ thuộc tham số, khó hơn viết cho đề cố định.
