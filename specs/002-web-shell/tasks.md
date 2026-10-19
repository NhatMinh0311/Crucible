# Danh sách việc: khung giao diện chung và giao diện hai ngôn ngữ

Mỗi việc là một commit, kiểm chứng được bằng cách ghi ở cột cuối.

| # | Việc | Kiểm chứng | Trạng thái |
|---|---|---|---|
| T1 | Realm phát triển: quyền `account` cho người dùng; hướng dẫn nhập lại realm | Token của `sv1` có audience `account`; `POST /realms/crucible/account` trả 204 | |
| T2 | Nền i18n: thư viện, bốn không gian, khai báo kiểu, chọn ngôn ngữ, thuộc tính `lang` | Kiểm thử đơn vị W1, W6 | |
| T3 | Chuyển văn bản hiện có sang khóa dịch; kiểm thử tập khóa và chữ viết thẳng | Kiểm thử W3; cố ý xóa một khóa `en` thì kiểm thử báo lỗi | |
| T4 | Định dạng ngày giờ, số; thông điệp lỗi API, mã lỗi dùng chung và của M1 | Kiểm thử W4, W5 với `TZ=America/New_York` | |
| T5 | Chế độ màu Sáng, Tối, Theo hệ thống, không nháy màu | Kiểm thử W7 | |
| T6 | Menu tài khoản: ngôn ngữ lưu vào Keycloak, chế độ màu, đăng xuất | Kiểm thử thành phần W2; trình duyệt thật | |
| T7 | Điều hướng trên điện thoại, bỏ qua tới nội dung, chuông thông báo, trang lỗi, `RequireDesktop` | Kiểm thử W8–W10 | |
| T8 | Tài liệu: `web/README.md` (thêm khóa dịch, định dạng, lỗi API), CR-DOC-11 nếu thiết kế đổi | Đọc lại | |
| T9 | Kiểm thử trên trình duyệt thật với Keycloak phát triển | W1, W2, W7 đạt; job CI `web` đạt (W11) | |
