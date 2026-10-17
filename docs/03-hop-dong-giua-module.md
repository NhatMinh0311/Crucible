# Hợp đồng giữa các module

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-DOC-03 |
| Phiên bản | 1.0 |
| Ngày ban hành | 17/10/2026 |
| Tài liệu liên quan | CR-DOC-02 Kiến trúc hệ thống, `contracts/`, ADR-0001, ADR-0005, ADR-0006, ADR-0010, ADR-0011 |

## 1. Phạm vi

Tài liệu quy định các hợp đồng mà một module cung cấp cho module khác, quy trình thay đổi hợp đồng và các quy ước
dùng chung. Thiết kế bên trong của mỗi module, bao gồm API chỉ phục vụ giao diện của chính module, không thuộc
phạm vi tài liệu này.

| Hợp đồng | Tệp | Bên cung cấp | Bên sử dụng |
|---|---|---|---|
| Đặc tả bài lab | `contracts/lab-spec/lab-spec.schema.json`, `validate.py`, `examples/` | M2 | Giảng viên (soạn bài), M1 (môi trường đã kết xuất) |
| Runtime API | `contracts/runtime-api.yaml` | M1 | M2, M3, giao diện |
| Content API | `contracts/content-api.yaml` | M2 | M3 |
| Sự kiện | `contracts/events.asyncapi.yaml` | Mỗi module cho sự kiện mình phát | Như ghi tại `x-consumers` |
| Giao thức terminal | Mục 6 của tài liệu này | M1 | Giao diện |

## 2. Quản lý hợp đồng

### 2.1. Quy trình thay đổi

1. Thay đổi hợp đồng được thực hiện bằng pull request sửa tệp trong `contracts/` trước khi sửa mã nguồn.
2. Pull request sửa hợp đồng cần được duyệt bởi chủ module cung cấp và chủ các module sử dụng (CODEOWNERS của `contracts/`).
3. CI kiểm tra tính hợp lệ của mọi tệp hợp đồng (CR-DOC-07 mục 6).

### 2.2. Tương thích và phiên bản

| Loại thay đổi | Ví dụ | Cách thực hiện |
|---|---|---|
| Tương thích | Thêm operation, thêm trường không bắt buộc ở phản hồi, thêm giá trị lỗi mới, thêm sự kiện mới | Sửa trực tiếp; tăng phiên bản phụ của tệp (`1.0.0` → `1.1.0`) |
| Không tương thích | Đổi tên hoặc xóa trường, thêm trường bắt buộc ở yêu cầu, đổi ý nghĩa giá trị | API: thêm đường dẫn `/v2/...`, giữ `/v1` tối thiểu 2 tuần sau khi mọi bên sử dụng đã chuyển. Sự kiện: phát song song kiểu sự kiện mới hậu tố `.v2` |

Bên nhận phải bỏ qua trường không biết trong phản hồi và trong sự kiện. Giá trị enum không biết được xử lý như lỗi
có thể bỏ qua và được ghi log cảnh báo.

### 2.3. Làm việc độc lập với module khác

| Nhu cầu | Cách đáp ứng |
|---|---|
| Gọi API của module chưa hoàn thành | Mock sinh từ OpenAPI bằng Prism: `npx @stoplight/prism-cli@5.14.2 mock contracts/runtime-api.yaml` |
| Có môi trường thật khi M1 chưa chạy trên Kubernetes | M1 cung cấp chế độ `docker` của `runtime-api`: cùng API, chạy phiên bằng container Docker trên máy phát triển, không có gVisor (chỉ dùng khi phát triển) |
| Nhận sự kiện của module khác | Tệp sự kiện mẫu trong kiểm thử của bên nhận, được kiểm tra theo lược đồ trong `events.asyncapi.yaml` |
| Bảo đảm bên cung cấp không làm lệch hợp đồng | Kiểm thử của bên cung cấp chạy bộ kiểm thử sinh từ OpenAPI (Schemathesis) trên dịch vụ đang chạy |

## 3. Quy ước chung cho API

### 3.1. Định danh và dữ liệu

- Mã đối tượng là UUID. Mã người dùng là `sub` của token Keycloak.
- Thời điểm theo RFC 3339 có múi giờ, lưu và truyền ở UTC.
- Mã tham chiếu mờ có dạng `<loại>:<mã>` (ví dụ `class:42`, `exam:7`). Bên nhận chỉ lưu và trả lại, không phân tích.
- Danh sách phân trang bằng `page` (từ 0) và `size` (tối đa 100).

### 3.2. Xác thực và scope

Giao diện gọi API bằng access token của người dùng. Dịch vụ gọi dịch vụ bằng token client credentials của tài khoản
dịch vụ trong Keycloak. Scope được cấp cho từng tài khoản dịch vụ như sau:

| Scope | Cấp cho | Cho phép |
|---|---|---|
| `runtime.sessions.write` | `content-service` | Tạo, gia hạn, kết thúc phiên |
| `runtime.sessions.read` | `content-service`, `learning-service` | Đọc phiên bất kỳ |
| `runtime.exec` | `content-service` | Chạy lệnh kiểm tra trong phiên |
| `runtime.view` | `learning-service` | Cấp vé xem chỉ đọc |
| `runtime.recordings.read` | `learning-service` | Lấy bản ghi phiên |
| `runtime.reservations.write` | `learning-service` | Đặt trước và hủy đặt trước tài nguyên |
| `runtime.quota.write` | `learning-service` | Thiết lập hạn mức của lớp trong trần của giảng viên; không sửa được trần |
| `content.attempts.write` | `learning-service` | Tạo lượt làm bài, nộp bài thay học viên khi hết giờ, điều chỉnh điểm |
| `content.attempts.read` | `learning-service` | Đọc lượt làm bài |
| `content.variants.read` | `learning-service` | Đọc giá trị biến thể (phân tích gian lận sau kỳ thi) |

Người dùng không phải chủ đối tượng nhận 404 giống hệt đối tượng không tồn tại.

### 3.3. Lặp lại an toàn

Operation tạo đối tượng yêu cầu header `Idempotency-Key` (UUID). Bên cung cấp lưu khóa cùng mã băm nội dung yêu cầu
và phản hồi trong 24 giờ: gửi lại cùng khóa và cùng nội dung trả về phản hồi đã lưu; cùng khóa nhưng khác nội dung
trả về 422 `IDEMPOTENCY_KEY_REUSED`.

### 3.4. Mã lỗi

Lỗi trả về theo RFC 9457 (`application/problem+json`) với trường `code` ổn định và `traceId`.

| Mã | HTTP | Bên trả về | Ý nghĩa |
|---|---|---|---|
| `VALIDATION_FAILED` | 400 | Mọi module | Yêu cầu sai cú pháp hoặc sai lược đồ |
| `UNAUTHENTICATED` | 401 | Mọi module | Thiếu hoặc sai token |
| `FORBIDDEN` | 403 | Mọi module | Token hợp lệ nhưng thiếu vai trò hoặc scope |
| `NOT_FOUND` | 404 | Mọi module | Không tồn tại hoặc không có quyền xem |
| `RATE_LIMITED` | 429 | Mọi module | Vượt giới hạn tần suất |
| `IDEMPOTENCY_KEY_REUSED` | 422 | Mọi module | Khóa lặp lại an toàn dùng cho nội dung khác |
| `SESSION_LIMIT` | 409 | M1 | Học viên đã có phiên cùng loại đang hoạt động |
| `QUOTA_EXHAUSTED` | 402 | M1 | Tài khoản hạn mức hết giờ máy hoặc hết số phiên đồng thời |
| `QUOTA_CEILING_EXCEEDED` | 422 | M1 | Tổng hạn mức chia cho các lớp vượt trần của giảng viên |
| `RESERVATION_BUDGET_EXHAUSTED` | 402 | M1 | Đặt trước làm vượt ngân sách chỗ × giờ trong tháng của giảng viên |
| `EXTENSION_LIMIT` | 409 | M1 | Gia hạn vượt mức tối đa hoặc phiên thi |
| `SESSION_NOT_ACTIVE` | 409 | M1 | Thao tác yêu cầu phiên ở pha `READY` hoặc `EXPIRING` |
| `EXEC_BUSY` | 409 | M1 | Phiên đang chạy một lệnh kiểm tra khác |
| `CAPACITY_UNAVAILABLE` | 409 | M1, M3 | Không đủ sức chứa cho đặt trước; kèm trường `capacity` gồm số chỗ cần và số chỗ còn |
| `ENVIRONMENT_NOT_ALLOWED` | 422 | M1 | Image ngoài danh sách cho phép hoặc mức tài nguyên không tồn tại |
| `LAB_NOT_PUBLISHED` | 409 | M2 | Bài lab hoặc phiên bản chưa xuất bản |
| `ATTEMPT_STATE_CONFLICT` | 409 | M2 | Thao tác không hợp lệ với trạng thái hiện tại của lượt làm bài |
| `VARIANT_UNSATISFIABLE` | 422 | M2 | Không sinh được biến thể thỏa ràng buộc `distinct` |
| `CHECK_DISABLED` | 409 | M2 | Ngữ cảnh của lượt làm bài không cho học viên kiểm tra |
| `PREREQUISITE_NOT_MET` | 403 | M3 | Chưa đạt bài tiên quyết |
| `ATTEMPT_LIMIT` | 409 | M3 | Hết số lượt làm bài của bài tập |
| `ASSIGNMENT_CLOSED` | 409 | M3 | Bài tập đã quá hạn |
| `EXAM_NOT_OPEN` | 409 | M3 | Kỳ thi chưa mở hoặc đã đóng |
| `EXAM_ENTRY_CLOSED` | 409 | M3 | Đã qua hạn chót được vào thi |
| `REGRADE_NOT_ALLOWED` | 409 | M3 | Quá thời hạn phúc khảo hoặc đã gửi yêu cầu cho bài này |

Lỗi `QUOTA_EXHAUSTED`, `QUOTA_CEILING_EXCEEDED`, `RESERVATION_BUDGET_EXHAUSTED`, `SESSION_LIMIT`, `CAPACITY_UNAVAILABLE`
của M1 được M2, M3 chuyển tiếp nguyên mã tới giao diện.

### 3.5. Tài khoản hạn mức

| Tài khoản | Mã tham chiếu | Tài khoản cha | Người đặt | Dùng cho |
|---|---|---|---|---|
| Trần của giảng viên | `instructor:<mã người dùng>` | Không có | Quản trị | Giới hạn tổng của mọi lớp do giảng viên phụ trách chính; phần chưa chia dùng cho phiên chạy thử bài lab (C13) |
| Lớp | `class:<mã lớp>` | `instructor:<mã giảng viên phụ trách chính>` | Giảng viên qua M3 | Phiên luyện tập, bài tập của sinh viên trong lớp; đặt trước cho kỳ thi của lớp |
| Kiểm thử bài lab | `system:lab-ci` | Không có | Quản trị | Phiên `LAB_CI` |

M3 tạo tài khoản của lớp khi tạo lớp và đổi tài khoản cha khi đổi giảng viên phụ trách chính. M2 tính phiên chạy thử vào
tài khoản `instructor:` của giảng viên đang chạy thử.

## 4. Đặc tả bài lab

### 4.1. Cấu trúc

| Phần | Nội dung |
|---|---|
| `metadata` | Mã định danh (`slug`), tiêu đề, độ khó, thời lượng ước tính, thẻ |
| `params` | Tham số và bộ sinh: `pick`, `range`, `pattern`, `subnet` |
| `constraints.distinct` | Nhóm tham số phải nhận giá trị khác nhau |
| `environment` | Image, mức tài nguyên, thời hạn, mạng, tệp ghi sẵn, script dựng |
| `faults` | Lỗi cài sẵn: số lỗi được chọn và danh sách lỗi |
| `steps` | Các bước: mô tả, điểm, lệnh kiểm tra, gợi ý, kết quả mong đợi ban đầu |
| `solution` | Lời giải mẫu cho từng bước |

Lược đồ đầy đủ tại `contracts/lab-spec/lab-spec.schema.json`. Hai bài mẫu tại `contracts/lab-spec/examples/` thể hiện
bài có tham số và bài có lỗi cài sẵn.

### 4.2. Quy tắc ngữ nghĩa

Ngoài lược đồ, đặc tả phải thỏa các quy tắc sau (kiểm tra bằng `contracts/lab-spec/validate.py`):

1. Mọi tham chiếu `{{tên}}` trong các trường mẫu đều tới tham số đã khai báo; mọi tham số đã khai báo đều được dùng.
2. Mã bước không trùng; mỗi bước có đúng một phần lời giải.
3. `range.min` < `range.max`; tham số trong `constraints.distinct` đã được khai báo.
4. Tổng điểm lớn hơn 0; tổng tỉ lệ trừ điểm của các gợi ý trong một bước không vượt quá 1.
5. Ít nhất một bước không đạt trên môi trường mới dựng.
6. `faults.pick` nhỏ hơn số lỗi trong danh sách; mã lỗi không trùng.

### 4.3. Kết xuất tham số

- Cú pháp `{{tên}}`, không có điều kiện hay vòng lặp. Các trường mẫu: nội dung tệp, `setup`, `inject`, mô tả bước,
  lệnh kiểm tra, thông điệp lỗi, gợi ý, lời giải.
- Giá trị tham số chỉ gồm ký tự an toàn trong shell (chữ, số, `._/:@-`), được bảo đảm bởi bộ sinh và lược đồ, nên được
  chèn trực tiếp vào script mà không cần thoát ký tự.
- M2 kết xuất toàn bộ đặc tả. M1 chỉ nhận phần `environment` đã kết xuất, đã nối script lỗi cài sẵn được chọn vào sau
  `setup`. M1 không nhận tham số, bước hay lời giải.

### 4.4. Sinh biến thể

| Mục | Quy định |
|---|---|
| Seed | `HMAC-SHA256(khóa bí mật của M2, contextRef + "|" + userId + "|" + labId + "|" + số thứ tự lượt làm bài)`. Khóa bí mật không rời M2, nên học viên không suy ra được biến thể của người khác |
| Giá trị mỗi tham số | Duyệt tham số theo thứ tự tên; với mỗi tham số lấy `h = HMAC-SHA256(seed, tên + "#" + lần thử)` và ánh xạ: `pick` lấy phần tử thứ `h mod n`; `range` lấy `min + h mod (max − min + 1)`; `pattern` lấy từng ký tự từ các byte liên tiếp của `h`; `subnet` lấy mạng con thứ `h mod số mạng con` |
| Ràng buộc `distinct` | Khi vi phạm, tăng `lần thử` của tham số đứng sau theo thứ tự tên; quá 100 lần trả lỗi `VARIANT_UNSATISFIABLE` |
| Lỗi cài sẵn | Sắp xếp lỗi theo `HMAC-SHA256(seed, "fault#" + id)` và lấy `pick` lỗi đầu |
| Ngữ cảnh | Luyện tập theo lộ trình: `track:<mã>`; bài tập: `assignment:<mã>`; kỳ thi: `exam:<mã>`. Tiếp tục một lượt làm bài giữ biến thể; lượt làm bài mới (số thứ tự tăng) nhận biến thể mới. Kỳ thi chỉ có một lượt làm bài cho mỗi học viên |

### 4.5. Kiểm thử bài lab khi xuất bản

M2 chỉ xuất bản một phiên bản khi các điều kiện sau đúng với 5 seed khác nhau, mỗi seed chạy trên một phiên ưu tiên `LAB_CI`:

1. Trên môi trường mới dựng: mọi bước `expectInitially: fail` không đạt, mọi bước `expectInitially: pass` đạt; bước
   `expectInitially: any` không được kiểm tra.
2. Sau khi chạy toàn bộ lời giải mẫu theo thứ tự bước: mọi bước đạt.
3. Không phiên nào kết thúc vì vượt giới hạn bộ nhớ (`OUT_OF_MEMORY`).

Sau khi kết thúc các phiên kiểm thử, M2 đọc `resourceUsage` của từng phiên. Nếu bộ nhớ cao nhất của một seed vượt 80%
giới hạn của mức tài nguyên đã chọn, phiên bản vẫn được xuất bản nhưng kết quả kiểm thử kèm cảnh báo và đề xuất mức
tài nguyên lớn hơn. Khi điều kiện 3 không đạt, kết quả đề xuất mức lớn hơn; nếu đã ở mức lớn nhất, đề xuất giảm nhu cầu
bộ nhớ của bài lab. Số đo dựa trên lời giải mẫu nên chỉ là ước tính cho thao tác thật của học viên.

## 5. Sự kiện

### 5.1. Định dạng

Sự kiện theo CloudEvents 1.0 dạng JSON có cấu trúc, gồm thêm hai thuộc tính mở rộng: `sequence` (số thứ tự tăng dần
theo `subject`) và `traceparent` (W3C Trace Context). Danh mục và lược đồ `data` tại `contracts/events.asyncapi.yaml`.

| Kiểu sự kiện | Bên phát | Bên nhận | `subject` |
|---|---|---|---|
| `runtime.session.ready` | M1 | M2 | Mã phiên |
| `runtime.session.expiring` | M1 | M2 | Mã phiên |
| `runtime.session.terminated` | M1 | M2, M3 | Mã phiên |
| `runtime.session.failed` | M1 | M2 | Mã phiên |
| `runtime.session.lost` | M1 | M2, M3 | Mã phiên |
| `runtime.terminal.connected` | M1 | M3 | Mã phiên |
| `runtime.terminal.input-line` | M1 | M3 | Mã phiên |
| `runtime.recording.available` | M1 | M3 | Mã phiên |
| `content.lab.published` | M2 | M1, M3 | Mã bài lab |
| `content.attempt.step-checked` | M2 | M3 | Mã lượt làm bài |
| `content.attempt.resumed` | M2 | M3 | Mã lượt làm bài |
| `content.attempt.graded` | M2 | M3 | Mã lượt làm bài |
| `learning.exam.closed` | M3 | M1 | Mã kỳ thi |

Sự kiện về phiên mang `labels` do bên tạo phiên gửi kèm (M2 gửi `attemptId`, `contextRef`), cho phép M3 liên kết
phiên với lượt làm bài mà M1 không cần biết khái niệm này.

### 5.2. Quy tắc phát và nhận

| Quy tắc | Nội dung |
|---|---|
| Phát | Sự kiện được ghi vào bảng `outbox` trong cùng transaction với thay đổi nghiệp vụ; tiến trình chuyển tiếp gửi lên exchange `crucible.events` với publisher confirm, đánh dấu đã gửi sau khi broker xác nhận |
| Hàng đợi | Mỗi module khai báo hàng đợi của mình, tên `<module>.<mục đích>` (ví dụ `content.session-events`), loại quorum, gắn với các routing key cần nhận |
| Loại trùng | Bên nhận ghi `id` sự kiện vào bảng `inbox` trong cùng transaction xử lý; sự kiện đã có trong `inbox` được xác nhận và bỏ qua |
| Thứ tự | Bên nhận lưu `sequence` cuối cùng đã xử lý theo `subject`; sự kiện có `sequence` nhỏ hơn hoặc bằng được bỏ qua |
| Lỗi | Thử lại 5 lần với khoảng chờ 1, 5, 30, 120, 600 giây; sau đó chuyển vào hàng đợi `<tên hàng đợi>.dlq` và phát cảnh báo |
| Giữ bản ghi | Dòng `outbox` đã gửi và dòng `inbox` được xóa sau 7 ngày |

## 6. Giao thức terminal

| Mục | Quy định |
|---|---|
| Địa chỉ | `wss://<máy chủ>/api/runtime/v1/terminal?ticket=<vé>` |
| Subprotocol | `crucible.terminal.v1` |
| Vé | Lấy từ `createTerminalTicket` (tương tác) hoặc `createViewTicket` (xem); dùng một lần, hết hạn sau 30 giây |
| Khung nhị phân | Dữ liệu terminal thô, hai chiều với kết nối tương tác; chỉ chiều máy chủ → trình duyệt với kết nối xem |
| Khung văn bản từ trình duyệt | JSON: `{"type":"resize","cols":120,"rows":32}` |
| Khung văn bản từ máy chủ | JSON: `{"type":"status","phase":"READY","expiresAt":"..."}`, `{"type":"expiring","expiresAt":"..."}`, `{"type":"closing","reason":"EXPIRED"}` |
| Giữ kết nối | Máy chủ gửi ping mỗi 20 giây; không có pong sau 60 giây thì đóng kết nối |
| Một kết nối tương tác | Mỗi phiên có tối đa một kết nối tương tác; kết nối mới đóng kết nối cũ với mã 4009 |

| Mã đóng | Ý nghĩa | Hành vi của giao diện |
|---|---|---|
| 1000 | Phiên kết thúc bình thường | Hiển thị kết quả |
| 4001 | Vé không hợp lệ hoặc đã dùng | Lấy vé mới, kết nối lại một lần |
| 4004 | Phiên không tồn tại hoặc đã kết thúc | Hiển thị trạng thái phiên |
| 4009 | Bị thay bởi kết nối mới hơn | Thông báo phiên đang mở ở nơi khác, không tự kết nối lại |
| 4010 | Máy chủ đang dừng để cập nhật | Lấy vé mới và kết nối lại sau 1–3 giây ngẫu nhiên |

## 7. Lệnh kiểm tra

| Mục | Quy định |
|---|---|
| Trình thông dịch | `/opt/crucible/bin/sh` (BusyBox liên kết tĩnh) |
| Công cụ | BusyBox, `curl`, `jq` liên kết tĩnh trong `/opt/crucible/bin`, gắn vào container từ volume chỉ đọc |
| Môi trường | Biến môi trường rỗng ngoài `PATH=/opt/crucible/bin`, `LANG=C.UTF-8`; không nạp tệp cấu hình shell; chạy bằng `root` |
| Kết quả | Mã thoát 0 là đạt; khác 0 là không đạt; hết thời gian là không đạt và `timedOut = true` |
| Giới hạn | Thời gian theo `timeoutSeconds`: tối đa 15 giây với `purpose = CHECK`, 60 giây với `purpose = RESTORE`; `stdout`, `stderr` cắt ở 4 KB; mỗi phiên chạy tối đa một lệnh tại một thời điểm |
| Khôi phục | Script khôi phục là lời giải mẫu của các bước đã đạt, nối theo thứ tự bước, chạy với `purpose = RESTORE` bằng `/bin/sh` của môi trường với `PATH` chuẩn của hệ điều hành, vì mục đích là tạo trạng thái chứ không phải đánh giá. Kiểm thử bài lab (mục 4.5) chạy lời giải theo cùng cách |
| Phạm vi bảo đảm | Kết quả không bị ảnh hưởng khi học viên sửa hoặc thay công cụ của hệ thống, biến môi trường, tệp cấu hình shell, `/etc/ld.so.preload`. Lệnh kiểm tra không nên gọi chương trình nằm ngoài `/opt/crucible/bin` (ADR-0011) |

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 17/10/2026 | Ban hành lần đầu |
