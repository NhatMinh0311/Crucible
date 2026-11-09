# Đặc tả: Quản lý lớp học và nhập danh sách sinh viên bằng CSV

| Thuộc tính | Giá trị |
|---|---|
| Mã | 003-m3-lop-hoc |
| Module | M3 (Học tập và kỳ thi) |
| Mốc | MA |
| Trạng thái | Nháp |
| Ngày lập | 09/11/2026 |
| Issue | #40 |
| Tài liệu liên quan | CR-DOC-01 mục 5.3 (L01, L12), mục 7.3 US-L1 (AC-L1.5); CR-DOC-03 mục 3.2; CR-DOC-10 UC11; CR-DOC-11 SCR-11, SCR-12, SCR-13, SCR-21; CR-MOD-M3 (BR-M3.05, BR-M3.26) |
| Yêu cầu ban đầu | "Lớp học, nhập danh sách CSV, giảng viên, trợ giảng (#40)" |

---

## 1. Mục tiêu

Giúp giảng viên tổ chức không gian lớp học cho từng học phần, thêm trợ giảng hỗ trợ, và nhập danh sách sinh viên nhanh chóng thông qua tệp CSV với cơ chế kiểm tra lỗi chi tiết từng dòng. Sinh viên trong danh sách khi đăng nhập bằng email trường sẽ tự động được ghi danh vào lớp mà không cần thao tác mã mời thủ công, sẵn sàng cho việc làm bài thực hành trong lộ trình luyện tập (chuẩn bị cho mốc MA).

---

## 2. Phạm vi

| Thuộc phạm vi | Ngoài phạm vi |
|---|---|
| Tạo, sửa thông tin, lưu trữ và đóng/lưu trữ lớp học (`ACTIVE`, `ARCHIVED`) | Lộ trình luyện tập của lớp (Issue #41) |
| Phân bổ hạn mức giờ máy của lớp từ trần hạn mức giảng viên (L12, BR-M3.26) | Giao bài tập cho lớp (Issue #42) |
| Quản lý thành viên lớp: thêm giảng viên phụ trách, thêm/gỡ trợ giảng | Lên lịch kỳ thi cho lớp (Issue #45) |
| Tải lên tệp CSV danh sách sinh viên: kiểm tra hợp lệ, báo lỗi chi tiết theo dòng (AC-L1.5, BR-M3.05) | Bảng điểm tổng kết học phần (Issue #48) |
| Tự động liên kết `user_id` của sinh viên khi đăng nhập lần đầu bằng email trường (AC-L1.5) | Quản lý tài khoản và mật khẩu (Keycloak đảm nhiệm) |
| Giao diện SCR-11 (Danh sách lớp GV), SCR-12 tab Thành viên, SCR-13 (Nhập CSV), SCR-21 (Danh sách lớp SV) | |

---

## 3. Câu chuyện người dùng

### US1 — Giảng viên tạo và quản lý lớp học (Ưu tiên P1; US-L1, L01, L12)

Là giảng viên, tôi muốn tạo một lớp học mới kèm mã lớp, học kỳ và phân bổ hạn mức giờ máy, để có không gian tổ chức bài thực hành cho sinh viên.

**Vì sao ưu tiên này**: Là điều kiện tiên quyết bắt buộc để có thể nhập sinh viên, giao bài lab hoặc tạo kỳ thi.

**Kiểm thử độc lập**: Tạo một lớp học qua API/giao diện, xác nhận lớp xuất hiện trong danh sách lớp của giảng viên với đầy đủ thông tin và hạn mức.

**Kịch bản chấp nhận**:
1. **Cho** giảng viên đã đăng nhập vào hệ thống, **khi** gửi yêu cầu tạo lớp với mã lớp hợp lệ (ví dụ `IT4060-2026.1`), tên lớp, học kỳ và hạn mức giờ máy, **thì** hệ thống lưu lớp học ở trạng thái `ACTIVE`, tự động gán giảng viên tạo lớp làm `INSTRUCTOR` của lớp và trả về thông tin chi tiết lớp học.
2. **Cho** giảng viên gửi yêu cầu tạo lớp, **khi** mã lớp đã tồn tại trong hệ thống, **thì** hệ thống từ chối với lỗi 409 `CLASS_CODE_ALREADY_EXISTS`.
3. **Cho** sinh viên đã đăng nhập, **khi** cố gắng gọi API tạo lớp, **thì** hệ thống từ chối với lỗi 403 `FORBIDDEN`.

---

### US2 — Giảng viên nhập danh sách sinh viên bằng tệp CSV (Ưu tiên P1; US-L1, AC-L1.5, BR-M3.05)

Là giảng viên, tôi muốn tải lên tệp CSV chứa danh sách sinh viên (mã sinh viên, họ tên, email trường) và nhận được báo cáo lỗi chi tiết theo từng dòng nếu có dữ liệu sai, để danh sách sinh viên được đưa vào lớp chính xác.

**Vì sao ưu tiên này**: Giảng viên không thể nhập thủ công hàng chục đến hàng trăm sinh viên mỗi lớp; cần cơ chế kiểm tra nghiêm ngặt để tránh nhầm lẫn dữ liệu.

**Kiểm thử độc lập**: Tải lên tệp CSV gồm các dòng đúng và các dòng sai (thiếu cột, sai định dạng email, trùng mã số sinh viên), kiểm chứng hệ thống chỉ lưu các dòng hợp lệ hoặc báo danh sách dòng lỗi chi tiết.

**Kịch bản chấp nhận**:
1. **Cho** giảng viên tải lên tệp CSV hợp lệ với tiêu đề `student_id,full_name,email` và 50 dòng dữ liệu đúng, **khi** xử lý nhập danh sách, **thì** hệ thống lưu 50 sinh viên vào bảng thành viên lớp với vai trò `STUDENT`, trạng thái chưa liên kết tài khoản (`user_id` rỗng), và trả về kết quả 50 bản ghi thành công.
2. **Cho** giảng viên tải lên tệp CSV có dòng 3 thiếu cột email, dòng 7 email không đúng định dạng RFC 5322, dòng 12 trùng mã sinh viên với dòng 2 trong cùng tệp, **khi** phân tích tệp, **thì** hệ thống trả về danh sách lỗi chi tiết theo từng số dòng kèm nguyên nhân lỗi rõ ràng (AC-L1.5).
3. **Cho** tệp CSV chứa sinh viên có email hoặc mã số sinh viên đã tồn tại sẵn trong lớp từ trước, **khi** nhập danh sách, **thì** hệ thống báo lỗi xung đột cho dòng đó và không ghi đè dữ liệu.

---

### US3 — Giảng viên phân quyền trợ giảng và giảng viên đồng phụ trách (Ưu tiên P2; L01)

Là giảng viên phụ trách lớp, tôi muốn thêm trợ giảng hoặc giảng viên khác vào lớp bằng email trường, để họ cùng quản lý bài tập và giám sát sinh viên.

**Vì sao ưu tiên này**: Nhiều lớp học đại học có trợ giảng phụ trách chấm bài hoặc hướng dẫn thực hành.

**Kiểm thử độc lập**: Giảng viên thêm một email với vai trò `TEACHING_ASSISTANT`; kiểm tra tài khoản đó thấy lớp học với quyền trợ giảng.

**Kịch bản chấp nhận**:
1. **Cho** giảng viên phụ trách lớp, **khi** thêm một thành viên bằng email trường với vai trò `TEACHING_ASSISTANT`, **thì** thành viên đó được gán vào lớp với vai trò trợ giảng.
2. **Cho** trợ giảng của lớp, **khi** cố gắng xóa lớp hoặc đổi hạn mức giờ máy của lớp, **thì** hệ thống từ chối với lỗi 403 `FORBIDDEN` (chỉ giảng viên phụ trách mới có quyền).

---

### US4 — Sinh viên tự động liên kết vào lớp khi đăng nhập (Ưu tiên P1; US-L1, AC-L1.5)

Là sinh viên có tên trong danh sách lớp, tôi muốn khi đăng nhập lần đầu vào hệ thống bằng tài khoản trường, tôi lập tức nhìn thấy các lớp mình đã được ghi danh mà không cần nhập mã tham gia.

**Vì sao ưu tiên này**: Đảm bảo trải nghiệm sinh viên liền mạch, không phụ thuộc vào việc phát tán mã mời hay sinh viên vào nhầm lớp.

**Kiểm thử độc lập**: Nhập một email sinh viên vào lớp; sau đó mô phỏng tài khoản sinh viên với email đó đăng nhập vào hệ thống; kiểm chứng `user_id` được cập nhật và sinh viên thấy lớp học tại SCR-21.

**Kịch bản chấp nhận**:
1. **Cho** một bản ghi sinh viên trong `class_members` có email `student1@sis.hust.edu.vn` và `user_id` là `null`, **khi** sinh viên đăng nhập qua Keycloak và gọi API lấy danh sách lớp của tôi (`GET /api/v1/classes`), **thì** hệ thống tự động cập nhật `user_id` bằng `sub` trong JWT token, và trả về danh sách lớp học mà sinh viên đã được ghi danh.
2. **Cho** sinh viên có email không nằm trong bất kỳ lớp nào, **khi** truy cập trang danh sách lớp học, **thì** hệ thống hiển thị danh sách trống kèm thông báo hướng dẫn liên hệ giảng viên.

---

## 4. Yêu cầu chi tiết

| Mã | Yêu cầu nghiệp vụ | Nguồn |
|---|---|---|
| **FR-001** | Giảng viên có quyền tạo lớp học mới gồm: mã lớp (duy nhất), tên lớp, học kỳ, hạn mức giờ máy phân bổ cho lớp, giới hạn số phiên đồng thời mỗi sinh viên (mặc định 1). | `L01`, `L12`, `AC-R2.8` |
| **FR-002** | Hệ thống hỗ trợ tải lên tệp CSV định dạng UTF-8; phân tích cú pháp tiêu đề gồm 3 cột bắt buộc: `student_id` (mã SV), `full_name` (họ tên), `email` (email trường). | `AC-L1.5`, `BR-M3.05` |
| **FR-003** | Hệ thống phải kiểm tra dữ liệu từng dòng trong CSV: kiểm tra đủ 3 cột, kiểm tra cú pháp email hợp lệ, kiểm tra không trùng mã sinh viên hoặc email trong nội bộ tệp CSV và không trùng với sinh viên đã có trong lớp. Báo cáo lỗi chi tiết theo chỉ số dòng (ví dụ: `Dòng 4: Email không hợp lệ`, `Dòng 8: Trùng mã sinh viên với dòng 2`). | `AC-L1.5`, `BR-M3.05` |
| **FR-004** | Hỗ trợ hai chế độ nhập CSV: chế độ xem trước (preview - phân tích lỗi không lưu vào CSDL) và chế độ xác nhận nhập (commit - lưu toàn bộ các dòng hợp lệ vào CSDL nếu không có lỗi chặn). | `SCR-13`, `BR-M3.05` |
| **FR-005** | Khi sinh viên đăng nhập hệ thống, `learning-service` đối chiếu claim `email` trong JWT token: nếu trùng với bản ghi thành viên lớp đang có `user_id IS NULL`, hệ thống tự động gán `user_id` bằng claim `sub` của token. | `AC-L1.5` |
| **FR-006** | Giảng viên có quyền thêm/gỡ trợ giảng (`TEACHING_ASSISTANT`) và giảng viên khác (`INSTRUCTOR`) vào lớp; một lớp có ít nhất một giảng viên phụ trách. | `L01` |
| **FR-007** | Sinh viên trong lớp chỉ xem được danh sách lớp học của mình, thông tin giảng viên và danh sách bạn cùng lớp (nếu được phép), không xem được danh sách lớp khác hoặc quyền quản trị. | `CR-DOC-02` |

### Thực thể chính

* **LopHoc (`classes`)**: Không gian học tập của một học phần trong một học kỳ. Thuộc sở hữu của giảng viên tạo lớp; có hạn mức giờ máy riêng.
* **ThanhVienLop (`class_members`)**: Mối liên kết giữa một người dùng (hoặc một bản ghi chờ đăng nhập) với một lớp học; mang vai trò `INSTRUCTOR`, `TEACHING_ASSISTANT` hoặc `STUDENT`.

---

## 5. Tiêu chí thành công

| Mã | Tiêu chí | Cách đo |
|---|---|---|
| **SC-001** | Tệp CSV 200 sinh viên được phân tích cú pháp, kiểm tra lỗi từng dòng và trả kết quả xem trước trong thời gian ≤ 1.5 giây. | Đo thời gian xử lý API nhập CSV |
| **SC-002** | 100% các dòng lỗi trong CSV (thiếu trường, email sai cú pháp, trùng lặp) được chỉ đúng số dòng và nguyên nhân lỗi trong phản hồi. | Bộ kiểm thử đơn vị CSV parser |
| **SC-003** | Sinh viên đăng nhập lần đầu được tự động gán vào đúng các lớp có email trùng lặp mà không cần bất kỳ thao tác nào từ phía người dùng. | Kiểm thử tích hợp tự động liên kết tài khoản |

---

## 6. Trường hợp biên

| Trường hợp biên | Xử lý mong đợi |
|---|---|
| Tệp CSV rỗng (chỉ có dòng tiêu đề hoặc 0 byte) | Trả về lỗi 422 `CSV_EMPTY`: Tệp tải lên không chứa dữ liệu sinh viên. |
| Tiêu đề CSV sai tên cột (ví dụ ghi `mssv` thay vì `student_id`) | Trả về lỗi 422 `CSV_HEADER_INVALID` kèm danh sách cột mong đợi (`student_id,full_name,email`). |
| Tệp CSV chứa ký tự phân tách là dấu chấm phẩy (`;`) hoặc tab | Tự động nhận diện dấu phân tách hoặc chuẩn hóa theo dấu phẩy (`,`). |
| Sinh viên có tên trong danh sách của 3 lớp khác nhau | Khi đăng nhập lần đầu, toàn bộ 3 bản ghi `class_members` ở cả 3 lớp đều được gán `user_id` của sinh viên đó. |
| Giảng viên nhập danh sách sinh viên trùng email nhưng khác mã sinh viên | Báo lỗi dòng đó với thông báo trùng email đã có trong lớp. |
| File CSV có dung lượng vượt quá 2 MB | Bị chặn từ tầng HTTP với mã 413 `PAYLOAD_TOO_LARGE`. |

---

## 7. Giả định và phụ thuộc

* **Keycloak**: Token JWT chứa claim `sub` (UUID định danh người dùng) và claim `email` của trường.
* **Hạn mức**: Dịch vụ M1 chưa bắt buộc phải có mặt để tạo lớp; hạn mức giờ máy của lớp được lưu tại M3 và sẽ đối soát khi sinh viên tạo phiên.
* **Giao diện**: Sử dụng các thành phần giao diện chuẩn từ `src/shell/ui` (Table, Dialog, Button, Input, Badge) theo hướng thiết kế "Bản vẽ kỹ thuật".

