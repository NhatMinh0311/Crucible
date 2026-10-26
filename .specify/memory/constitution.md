# Hiến pháp dự án Crucible

Tài liệu này là hiến pháp mà các lệnh Spec Kit (`specify`, `clarify`, `plan`, `tasks`, `analyze`, `implement`,
`converge`) đọc khi lập đặc tả, kế hoạch và danh sách việc. Nội dung tóm tắt các quy tắc đã có trong `AGENTS.md`,
`CONTRIBUTING.md` và CR-DOC-07; khi có khác biệt, `AGENTS.md` được ưu tiên.

## Nguyên tắc

### I. Tài liệu đã duyệt là nguồn yêu cầu

Chức năng, tiêu chí chấp nhận `AC-*` và câu chuyện người dùng `US-*` lấy từ `docs/01-bai-toan-va-pham-vi.md`; quy tắc
nghiệp vụ `BR-*`, mô hình dữ liệu và máy trạng thái lấy từ `docs/modules/<module>.md`. Đặc tả dẫn mã của các mục này
thay cho việc viết lại. Trường hợp tài liệu không trả lời được thì đánh dấu `[NEEDS CLARIFICATION: <câu hỏi>]` để
thành viên quyết định; trợ lý AI không tự chọn phương án nghiệp vụ.

### II. Hợp đồng trước, mã sau

API và sự kiện giữa các module nằm ở `contracts/` tại gốc kho mã và được quản lý theo CR-DOC-03 mục 2. Đặc tả và kế
hoạch chỉ dẫn tới hợp đồng; không tạo thư mục `contracts/` trong `specs/<số>-<tên>/`. Tính năng cần đổi hợp đồng thì
kế hoạch ghi rõ phần thay đổi, và thay đổi hợp đồng được đưa vào một pull request riêng, merge trước pull request mã.

### III. Kiểm thử là bắt buộc

Mỗi tiêu chí chấp nhận và quy tắc nghiệp vụ được hiện thực phải có kiểm thử tự động mang mã của nó trong tên hoặc chú
thích (ví dụ `AC-C2.4`, `BR-M1.3`). Danh sách việc luôn có việc viết kiểm thử, kể cả khi đặc tả không nhắc tới.
Kiểm thử hợp đồng của bên cung cấp phải đạt; không tắt hoặc nới kiểm thử để CI đạt.

### IV. Ranh giới module

Mỗi tính năng thuộc một module (M1, M2, M3 hoặc nền tảng) và chỉ sửa mã trong thư mục của module đó. Module không đọc
dữ liệu của module khác ngoài API và sự kiện trong `contracts/`. Việc cần làm ở module khác được ghi thành phụ thuộc
trong kế hoạch, không đưa vào danh sách việc.

### V. Vận hành và bảo mật

Không đưa bí mật vào mã; log dạng JSON có `trace_id`, không ghi token, vé terminal, nội dung lời giải. Lưu thời gian
theo UTC, hiển thị theo `Asia/Ho_Chi_Minh`. Migration đã merge không được sửa; thay đổi lược đồ theo nguyên tắc mở rộng
rồi thu hẹp. Chức năng chỉ được coi là hoàn thành khi chạy trên môi trường chuẩn qua script triển khai và có metric,
log đủ để chẩn đoán (CR-DOC-06 mục 6).

### VI. Ngôn ngữ và quy ước

Đặc tả, kế hoạch, danh sách việc viết bằng tiếng Việt, văn phong chính thức; tên định danh trong mã, khóa cấu hình và
mã lỗi viết bằng tiếng Anh. Văn bản hiển thị trên giao diện đi qua khóa dịch, đủ tiếng Việt và tiếng Anh. Quy ước mã
nguồn theo `AGENTS.md` mục 4.

## Tệp của một tính năng

| Tệp | Nội dung | Ghi chú |
|---|---|---|
| `specs/<số>-<tên>/spec.md` | Mục tiêu, phạm vi, câu chuyện người dùng, yêu cầu, tiêu chí thành công, trường hợp biên | Mục `Clarifications` của Spec Kit tương ứng với mục "Điểm đã làm rõ" của mẫu |
| `specs/<số>-<tên>/plan.md` | Thiết kế kỹ thuật trong module, lựa chọn kỹ thuật kèm lý do, dữ liệu, ảnh hưởng tới hợp đồng, kiểm thử | Nội dung nghiên cứu, mô hình dữ liệu và cách kiểm chứng ghi trong `plan.md`; không tạo `research.md`, `data-model.md`, `quickstart.md` riêng |
| `specs/<số>-<tên>/tasks.md` | Danh sách việc có thứ tự, mỗi việc kiểm chứng được | Định dạng `- [ ] T001 [P] [US1] ...` để lệnh `implement`, `converge` đánh dấu |

Thư mục `checklists/` do lệnh `specify` sinh ra chỉ dùng trên máy, không commit. Số thứ tự của tính năng là số tiếp
theo trong `specs/` sau khi đã cập nhật `main`.

## Quản trị

Pull request đặc tả và kế hoạch được chủ module và trưởng nhóm duyệt trước khi viết mã; pull request hiện thực dẫn tới
thư mục `specs/` tương ứng (CR-DOC-07 mục 7). Thay đổi hiến pháp đi qua pull request như mã nguồn và được ghi vào lịch
sử dưới đây.

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 26/10/2026 | Ban hành lần đầu, dùng với Spec Kit 1.2.3 |
