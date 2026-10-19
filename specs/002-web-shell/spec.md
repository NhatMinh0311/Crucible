# Đặc tả: khung giao diện chung và giao diện hai ngôn ngữ

| Thuộc tính | Giá trị |
|---|---|
| Mã | 002-web-shell |
| Module | Nền tảng |
| Mốc | MA |
| Trạng thái | Đã duyệt; bổ sung W12 ngày 20/10/2026 |
| Issue | #3, #4 |
| Tài liệu liên quan | CR-DOC-01 P02, P05, US-P2 (AC-P2.1–AC-P2.3); CR-DOC-03 mục 3.4; CR-DOC-11 mục 2, 3, 4, 4.1, 5 |

## 1. Mục tiêu

Thành viên M2, M3 dựng màn hình trên một khung chung đã có sẵn: đăng nhập, menu theo vai trò, menu tài khoản, chế độ
sáng tối, hai ngôn ngữ tiếng Việt và tiếng Anh, cách hiển thị lỗi API và định dạng ngày giờ dùng chung. Mọi văn bản đi
qua khóa dịch ngay từ màn hình đầu tiên; CI chặn khóa thiếu ở một ngôn ngữ và văn bản viết thẳng trong JSX, để việc
dịch không bị dồn về cuối dự án.

## 2. Phạm vi

| Thuộc phạm vi | Ngoài phạm vi |
|---|---|
| i18next với bốn không gian `shell`, `runtime`, `content`, `learning`; tệp dịch `vi`, `en` | Bản dịch lỗi riêng của M2, M3 (mỗi module điền vào không gian của mình) |
| Chọn ngôn ngữ lần đầu theo trình duyệt; đổi ở menu tài khoản; lưu thuộc tính `locale` qua Account API của Keycloak | Mẫu email, thông báo hai ngôn ngữ (M3, #43) |
| Trang đăng nhập Keycloak hiển thị theo ngôn ngữ đang chọn | Theme riêng cho trang đăng nhập (dùng theme mặc định, đã có tiếng Việt) |
| Hiển thị lỗi API theo `code`, kèm mã tra cứu | Ngôn ngữ nội dung bài lab (M2, #35) |
| Định dạng ngày giờ, số theo ngôn ngữ, múi giờ Việt Nam | Nội dung trang thông báo và số thông báo chưa đọc (M3, #42) |
| Menu tài khoản: tên, email, vai trò, ngôn ngữ, chế độ màu, đăng xuất | Nguyên mẫu SCR-12, SCR-13, SCR-24, SCR-27 (CR-DOC-11 mục 9) |
| Chế độ sáng, tối, theo hệ thống | |
| Điều hướng trên điện thoại, liên kết bỏ qua tới nội dung, chuông thông báo dẫn tới `/notifications` | |
| Trang lỗi của route, thành phần chỉ cho máy tính (`RequireDesktop`) | |
| Realm phát triển: quyền `manage-account` cho người dùng | |
| Giao diện theo hướng "Bản vẽ kỹ thuật" (W12): phông, token màu, thành phần, bố cục khung | Màn hình riêng của module (dùng token và thành phần chung) |

## 3. Tiêu chí chấp nhận

| Mã | Tiêu chí |
|---|---|
| W1 | (AC-P2.2) Token không có claim `locale` hợp lệ và máy chưa lưu lựa chọn: ngôn ngữ ưu tiên đầu tiên của trình duyệt bắt đầu bằng `en` thì dùng tiếng Anh, còn lại tiếng Việt. Trang đăng nhập Keycloak hiển thị cùng ngôn ngữ với giao diện |
| W2 | (AC-P2.2) Đổi ngôn ngữ ở menu tài khoản: giao diện đổi ngay, không tải lại trang; thuộc tính `locale` của tài khoản được lưu; đăng nhập trên trình duyệt khác hiển thị ngôn ngữ đã chọn. Lưu thất bại: giao diện vẫn dùng ngôn ngữ mới trên máy này và báo không lưu được vào tài khoản |
| W3 | (AC-P2.1) Kiểm thử báo lỗi khi: hai ngôn ngữ khác tập khóa ở một không gian; JSX chứa chữ viết thẳng; thuộc tính `aria-label`, `title`, `placeholder`, `alt` là chuỗi chữ viết thẳng. Dùng khóa không tồn tại là lỗi kiểu TypeScript |
| W4 | (AC-P2.3) Lỗi API hiển thị thông điệp theo `code`, tra ở không gian của module trả lỗi rồi tới `shell`; không có bản dịch thì hiển thị thông điệp chung kèm `code`. Có `traceId` thì kèm "Mã tra cứu" gồm 8 ký tự đầu. Không hiển thị `detail` |
| W5 | (AC-P2.3) Ngày giờ và số định dạng theo ngôn ngữ đang chọn; giờ luôn theo múi giờ `Asia/Ho_Chi_Minh` bất kể múi giờ của máy |
| W6 | Thuộc tính `lang` của trang theo ngôn ngữ đang chọn |
| W7 | Chế độ màu: lần đầu theo cài đặt của hệ điều hành; người dùng chọn Sáng, Tối hoặc Theo hệ thống ở menu tài khoản; lựa chọn lưu trên trình duyệt; tải trang không nháy sai màu; ở chế độ Theo hệ thống, giao diện đổi theo khi hệ điều hành đổi |
| W8 | Menu chỉ hiện mục người dùng có vai trò; gõ thẳng đường dẫn của route khai báo vai trò khác hiển thị trang "không có quyền" trong khung. Dưới 768 px menu ẩn sau nút, mở ra có lớp phủ, đóng khi bấm lớp phủ, nhấn Esc hoặc chọn mục. Liên kết "Bỏ qua tới nội dung" là phần tử nhận tiêu điểm đầu tiên. Mọi thao tác trên khung dùng được bằng bàn phím (UX7) |
| W9 | `RequireDesktop`: màn hình hẹp hơn 1024 px hiển thị thông báo yêu cầu dùng máy tính thay cho nội dung (CR-DOC-11 mục 3) |
| W10 | Đường dẫn không tồn tại và lỗi khi hiển thị route có trang lỗi đã dịch, kèm nút về trang chủ |
| W11 | Job CI `web` đạt: kiểm tra kiểu, lint, kiểm thử, build |
| W12 | Giao diện theo hướng "Bản vẽ kỹ thuật" đã chọn: Space Grotesk cho tiêu đề, IBM Plex Sans cho nội dung, IBM Plex Mono cho nhãn, số liệu, terminal; góc vuông; đường kẻ 1 px thay cho bóng đổ; nền lưới 24 px; xanh bản vẽ làm màu chính, cam tín hiệu cho cảnh báo và thời gian; menu đánh số. Chữ thường đạt độ tương phản 4,5:1 ở cả chế độ sáng và tối |

## 4. Trường hợp biên

| Trường hợp | Xử lý |
|---|---|
| Claim `locale` khác `vi`, `en` | Bỏ qua claim, chọn theo mục W1 |
| Ngôn ngữ lưu trên máy khác claim `locale` | Claim thắng, vì là lựa chọn gần nhất được lưu vào tài khoản |
| Account API lỗi (mất mạng, token hết hạn) | Giữ ngôn ngữ mới trên máy, báo không lưu được vào tài khoản |
| Người dùng gửi họ tên, email khác qua Account API | Keycloak trả 400 `error-user-attribute-read-only`; chỉ quản trị sửa họ tên, email (User Profile); `locale` vẫn do người dùng sửa |
| Trình duyệt chặn `localStorage` | Giao diện vẫn chạy, không nhớ lựa chọn giữa các lần tải trang |
| Lỗi không theo RFC 9457 (mất mạng, phản hồi không phải JSON) | Thông điệp chung "Không kết nối được tới máy chủ" hoặc "Đã xảy ra lỗi" |

## 5. Điểm đã làm rõ

| Câu hỏi | Quyết định |
|---|---|
| Chế độ màu mặc định, nơi lưu | Theo hệ điều hành; lựa chọn lưu trên trình duyệt (không lưu theo tài khoản) |
| Nội dung hiển thị kèm lỗi API | Thông điệp theo `code` và mã tra cứu từ `traceId` |
| Tổ chức công việc | Một đặc tả, một pull request đóng #3 và #4 |
| Hướng giao diện | So sánh ba hướng (phòng lab, lò luyện, bản vẽ kỹ thuật) dựng từ các DESIGN.md của designmd.app, getdesign.md; chọn "Bản vẽ kỹ thuật", giữ nguyên chi tiết bản dựng; đưa vào cùng pull request vì #3 gồm theme |
| Trang đăng nhập | Theme mặc định của Keycloak 26.8.1 (có `messages_vi`); ngôn ngữ truyền bằng tham số `ui_locales` |
| Quyền gọi Account API | Người dùng có vai trò `manage-account`, `view-profile` của client `account`; đã kiểm chứng với Keycloak phát triển: ghi `locale` trả 204, CORS cho `http://localhost:5173`, token sau đó có claim `locale` |
