# Hợp đồng giữa các module

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-DOC-03 |
| Phiên bản | 1.1 |
| Ngày ban hành | 27/10/2026 |
| Tài liệu liên quan | CR-DOC-02 Kiến trúc hệ thống, `contracts/`, ADR-0001, ADR-0005, ADR-0006, ADR-0010, ADR-0011, ADR-0016 |

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
| Có môi trường thật trên máy phát triển | `runtime-api` chạy trên máy phát triển với cùng API, chạy phiên bằng container Docker với runtime `runc`, không có gVisor. Chạy bằng `docker compose -f dev/compose.yaml --profile runtime up -d --build --wait` (CONTRIBUTING mục 5.2) |
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
| `METHOD_NOT_ALLOWED` | 405 | Mọi module | Phương thức HTTP không có ở đường dẫn |
| `RATE_LIMITED` | 429 | Mọi module | Vượt giới hạn tần suất |
| `IDEMPOTENCY_KEY_REUSED` | 422 | Mọi module | Khóa lặp lại an toàn dùng cho nội dung khác |
| `NOT_IMPLEMENTED` | 501 | M1 | Operation chưa được hiện thực |
| `SESSION_LIMIT` | 409 | M1 | Học viên đã có phiên cùng loại đang hoạt động |
| `QUOTA_EXHAUSTED` | 402 | M1 | Tài khoản hạn mức hết giờ máy hoặc hết số phiên đồng thời |
| `QUOTA_CEILING_EXCEEDED` | 422 | M1 | Tổng hạn mức chia cho các lớp vượt trần của giảng viên |
| `RESERVATION_BUDGET_EXHAUSTED` | 402 | M1 | Đặt trước làm vượt ngân sách chỗ × giờ trong tháng của giảng viên |
| `EXTENSION_LIMIT` | 409 | M1 | Gia hạn vượt mức tối đa hoặc phiên thi |
| `SESSION_NOT_ACTIVE` | 409 | M1 | Thao tác yêu cầu phiên ở pha `READY` hoặc `EXPIRING` |
| `EXEC_BUSY` | 409 | M1 | Phiên đang chạy một lệnh kiểm tra khác |
| `RESERVATION_SLOTS_IN_USE` | 409 | M1 | Giảm số chỗ đặt trước xuống dưới số chỗ đang dùng |
| `RESERVATION_CLOSED` | 409 | M1 | Đặt trước đã hủy hoặc đã qua giờ kết thúc |
| `CAPACITY_UNAVAILABLE` | 409 | M1, M3 | Không đủ sức chứa cho đặt trước; kèm trường `capacity` gồm số chỗ cần và số chỗ còn |
| `ENVIRONMENT_NOT_ALLOWED` | 422 | M1 | Image ngoài danh sách cho phép, mức tài nguyên không tồn tại hoặc capability không được cấp trong bản cài đặt |
| `LAB_NOT_PUBLISHED` | 409 | M2 | Bài lab hoặc phiên bản chưa xuất bản |
| `ATTEMPT_STATE_CONFLICT` | 409 | M2 | Thao tác không hợp lệ với trạng thái hiện tại của lượt làm bài |
| `VARIANT_UNSATISFIABLE` | 422 | M2 | Không sinh được biến thể thỏa ràng buộc `distinct` |
| `CHECK_DISABLED` | 409 | M2 | Ngữ cảnh của lượt làm bài không cho học viên kiểm tra |
| `ANSWER_TRIES_EXHAUSTED` | 409 | M2 | Bước trả lời đã dùng hết số lần trả lời |
| `CHECK_LOCKED` | 409 | M2 | Bài lab yêu cầu trả lời bước chẩn đoán trước khi kiểm tra |
| `PREREQUISITE_NOT_MET` | 403 | M3 | Chưa đạt bài tiên quyết |
| `ATTEMPT_LIMIT` | 409 | M3 | Hết số lượt làm bài của bài tập |
| `ASSIGNMENT_CLOSED` | 409 | M3 | Bài tập đã quá hạn |
| `EXAM_NOT_OPEN` | 409 | M3 | Kỳ thi chưa mở hoặc đã đóng |
| `EXAM_ENTRY_CLOSED` | 409 | M3 | Đã qua hạn chót được vào thi |
| `EXTRA_TIME_LIMIT` | 409 | M3 | Cộng giờ làm tổng thời gian cộng của sinh viên vượt trần `exam.extraTimeMaxMinutes` đã chốt của kỳ thi, hoặc sinh viên đã nộp |
| `REGRADE_NOT_ALLOWED` | 409 | M3 | Quá thời hạn phúc khảo hoặc đã gửi yêu cầu cho bài này |
| `SETTING_OUT_OF_RANGE` | 422 | M1, M2, M3 | Giá trị tham số hệ thống ngoài khoảng hợp lệ |

Giao diện hiển thị lỗi theo `code` bằng ngôn ngữ của người dùng; `detail` dành cho nhà phát triển và không được hiển thị
cho người dùng (CR-DOC-11 mục 4.1).

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

### 3.6. Tham số hệ thống

Mỗi module giữ tham số nghiệp vụ của mình (CR-DOC-01 mục 5.6) trong CSDL riêng và cung cấp cùng một bộ API, chỉ cho
vai trò `ADMIN`. M1, M2 mô tả trong `runtime-api.yaml`, `content-api.yaml`; M3 hiện thực đúng khuôn dạng này cho API của
mình. Giao diện SCR-42 gọi API của ba module và hiển thị chung một màn hình.

| Operation | Nội dung |
|---|---|
| `GET /v1/settings` | Danh sách tham số của module: khóa, giá trị, mặc định, khoảng hợp lệ, đơn vị, mô tả, lần sửa gần nhất |
| `PUT /v1/settings/{key}` | Sửa giá trị; ngoài khoảng hợp lệ trả 422 `SETTING_OUT_OF_RANGE` |
| `GET /v1/settings/{key}/history` | Lịch sử sửa: giá trị cũ, giá trị mới, người sửa, thời điểm |

Mỗi bản dịch vụ đọc tham số qua bộ nhớ đệm làm mới không quá 60 giây. Giá trị cần cố định cho một đối tượng được chép
vào đối tượng đó khi tạo: kỳ thi chép các tham số `exam.*`, phiên bản bài lab chép `maxTries` mặc định khi xuất bản.
Module khác không đọc tham số của nhau; giá trị cần chuyển qua ranh giới được truyền trong lời gọi, ví dụ M3 truyền
`sessionUntil` (giờ kết thúc chung cộng hai trần đã chốt của kỳ thi) khi tạo lượt thi để M2 đặt thời hạn phiên.

## 4. Đặc tả bài lab

### 4.1. Cấu trúc

| Phần | Nội dung |
|---|---|
| `metadata` | Mã định danh (`slug`), tiêu đề, ngôn ngữ nội dung (`language`: `vi`, `en`), độ khó, thời lượng ước tính, thẻ |
| `params` | Tham số và bộ sinh: `pick`, `range`, `pattern`, `subnet` |
| `constraints.distinct` | Nhóm tham số phải nhận giá trị khác nhau |
| `environment` | Image, mức tài nguyên, thời hạn, mạng, khai báo dùng CPU cao, tệp ghi sẵn, script dựng |
| `faults` | Lỗi cài sẵn: số lỗi được chọn và danh sách lỗi |
| `steps` | Các bước: loại (`check` – kiểm tra trạng thái môi trường, `answer` – trả lời câu hỏi), mô tả, điểm, lệnh kiểm tra hoặc đáp án, gợi ý, kết quả mong đợi ban đầu |
| `solution` | Lời giải mẫu cho từng bước: script với bước `check`; câu trả lời viết sẵn hoặc script in ra câu trả lời với bước `answer` |

Lược đồ đầy đủ tại `contracts/lab-spec/lab-spec.schema.json`. Bốn bài mẫu tại `contracts/lab-spec/examples/` thể hiện
các dạng bài tại CR-DOC-01 mục 5.4: cấu hình theo yêu cầu (`linux-user-nginx`), sửa hệ thống hỏng (`fix-broken-web`),
bảo mật và gia cố (`harden-server`), điều tra và trả lời (`investigate-web-log`). Bài `fix-broken-web` có bước chẩn đoán
trắc nghiệm lấy đáp án theo lỗi cài sẵn.

Trường `machine` của bước và lỗi cài sẵn hiện chỉ nhận `main`; môi trường nhiều máy mở rộng giá trị này (ADR-0015).

Trường `environment.capabilities` hiện chỉ nhận `NET_ADMIN`, dành cho bài tường lửa, cấu hình địa chỉ và tuyến
(ADR-0016). M2 chuyển trường này cho M1 trong `RenderedEnvironment` và trả trong `LabSummary`; M3 truyền giá trị đó vào
`warmCapabilities` khi đặt trước tài nguyên cho kỳ thi để M1 kiểm tra trước capability của kỳ thi.

### 4.2. Quy tắc ngữ nghĩa

Ngoài lược đồ, đặc tả phải thỏa các quy tắc sau (kiểm tra bằng `contracts/lab-spec/validate.py`):

1. Mọi tham chiếu `{{tên}}` trong các trường mẫu đều tới tham số đã khai báo; mọi tham số đã khai báo đều được dùng.
2. Mã bước không trùng; mỗi bước có đúng một phần lời giải.
3. `range.min` < `range.max`; tham số trong `constraints.distinct` đã được khai báo.
4. Tổng điểm lớn hơn 0; tổng tỉ lệ trừ điểm của các gợi ý trong một bước không vượt quá 1.
5. Ít nhất một bước không đạt trên môi trường mới dựng.
6. `faults.pick` nhỏ hơn số lỗi trong danh sách; mã lỗi không trùng.
7. Bước `check` có lời giải `run`; bước `answer` có lời giải `answer` hoặc `run`.
8. Biểu thức `answer.regex` hợp lệ.
9. Trắc nghiệm: mã phương án không trùng; `correct` chỉ gồm mã có trong `choices`; nhiều phương án đúng cần `multiple: true`;
   với `correctFromFaults`, mọi lỗi cài sẵn có phương án cùng mã, và bài chọn nhiều lỗi (`faults.pick` > 1) cần `multiple: true`.

### 4.3. Kết xuất tham số

- Cú pháp `{{tên}}`, không có điều kiện hay vòng lặp. Các trường mẫu: nội dung tệp, `setup`, `inject`, mô tả bước,
  lệnh kiểm tra, thông điệp lỗi, gợi ý, lời giải.
- Giá trị tham số chỉ gồm ký tự an toàn trong shell (chữ, số, `._/:@-`), được bảo đảm bởi bộ sinh và lược đồ, nên được
  chèn trực tiếp vào script mà không cần thoát ký tự.
- Biến có sẵn `{{faults}}`: mã các lỗi cài sẵn được chọn cho biến thể, sắp theo tên, nối bằng dấu phẩy (ví dụ
  `portConflict,serviceStopped`). Chỉ dùng trong bài có lỗi cài sẵn; `faults` không được dùng làm tên tham số.
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
   `expectInitially: any` và bước `answer` không được kiểm tra.
2. Sau khi chạy toàn bộ lời giải mẫu theo thứ tự bước: mọi bước `check` đạt; câu trả lời trong lời giải của mỗi bước
   `answer` (viết sẵn, hoặc stdout của script lời giải) được chấp nhận theo mục 4.6.
3. Không phiên nào kết thúc vì vượt giới hạn bộ nhớ (`OUT_OF_MEMORY`).
4. Không lần chạy nào trả `invalidCommand = true` (mục 7): lệnh kiểm tra hoặc lời giải gọi chương trình không có
   trong môi trường hoặc không chạy được là lỗi đặc tả, không phải kết quả của học viên.
5. Không phiên nào kết thúc vì môi trường dừng (`ENVIRONMENT_FAILED`): `setup` kết thúc bằng dịch vụ chạy ở tiền
   cảnh, hoặc script tạo tiến trình vượt giới hạn 512 tác vụ, đều biểu hiện ở đây (mục 4.7).

Sau khi kết thúc các phiên kiểm thử, M2 đọc `resourceUsage` của từng phiên. Nếu bộ nhớ cao nhất của một seed vượt 80%
giới hạn của mức tài nguyên đã chọn, phiên bản vẫn được xuất bản nhưng kết quả kiểm thử kèm cảnh báo và đề xuất mức
tài nguyên lớn hơn. Khi điều kiện 3 không đạt, kết quả đề xuất mức lớn hơn; nếu đã ở mức lớn nhất, đề xuất giảm nhu cầu
bộ nhớ của bài lab. Số đo dựa trên lời giải mẫu nên chỉ là ước tính cho thao tác thật của học viên.

### 4.6. Bước trả lời

| Mục | Quy định |
|---|---|
| Chuẩn hóa | Bỏ khoảng trắng đầu và cuối; không phân biệt hoa thường trừ khi `caseSensitive: true`. Áp dụng cho câu trả lời và đáp án |
| Đúng khi (tự luận) | Khớp một giá trị trong `accepted` (đã kết xuất tham số), hoặc khớp toàn bộ `regex` (RE2), hoặc bằng stdout của `from.run` |
| Trắc nghiệm | Phương án hiển thị theo thứ tự xáo trộn bằng `HMAC-SHA256(seed, "choice#" + id)`. Tập đúng là `correct`, hoặc các lỗi cài sẵn được chọn khi `correctFromFaults`. Chọn một: đúng thì đạt. Chọn nhiều: điểm = điểm bước × max(0, (số phương án đúng đã chọn − số phương án sai đã chọn) / số phương án đúng); đạt khi bằng điểm bước. Câu trả lời gửi lên là danh sách mã phương án |
| `answerFirst` | Trước khi bước có `answerFirst` được trả lời (đúng, sai hoặc hết lượt), yêu cầu kiểm tra các bước `check` bị từ chối (409 `CHECK_LOCKED`). Không ảnh hưởng chấm ngầm và chấm cuối |
| `from.run` | Chạy như lệnh kiểm tra (mục 7) tại thời điểm học viên trả lời; kết quả so khớp được lưu cùng câu trả lời |
| Chế độ `FULL` | Học viên biết ngay đúng hay sai; tối đa `maxTries` lần (mặc định 3 với tự luận, 1 với trắc nghiệm) cho mỗi bước trong một lượt làm bài; điểm của bước là điểm của lần trả lời cao nhất; hết lượt mà chưa đúng thì bước không đạt (409 `ANSWER_TRIES_EXHAUSTED`); trả lời đúng thì bước khóa lại |
| Chế độ `RECORDED_ONLY`, `DISABLED` | Học viên vẫn nhập được câu trả lời, chỉ nhận thông báo đã ghi nhận; không giới hạn số lần; câu trả lời cuối cùng được dùng khi chấm cuối |
| Gợi ý | Như bước `check` |
| Khôi phục khi tiếp tục | Bước `answer` không cần khôi phục môi trường; kết quả đã ghi được giữ |
| Chấm cuối | Bước `answer` lấy kết quả của câu trả lời đúng hoặc câu trả lời cuối cùng đã ghi; không chạy lại `from.run` |

### 4.7. Môi trường khi script của bài lab chạy

Áp dụng cho `environment.setup`, script lỗi cài sẵn (`faults[].inject`) và lời giải mẫu chạy với `purpose = RESTORE`.

| Mục | Quy định |
|---|---|
| Hệ thống khởi tạo | Môi trường không có `systemd`. `systemctl` không dùng được; `service` và `/etc/init.d/*` dùng được vì chỉ gọi `start-stop-daemon`. Tiến trình `init` của môi trường là `tini`, có chuyển tiếp tín hiệu và thu dọn tiến trình mồ côi |
| Dịch vụ nền | Script được chạy với `stdin` từ `/dev/null` và `stdout`, `stderr` vào `/var/log/crucible/setup.log` của môi trường, do M1 đặt trước khi chạy (BR-M1.29). Dịch vụ nền thừa hưởng ba bộ mô tả tệp này, nên đặc tả bài lab không cần tự chuyển hướng |
| Dịch vụ chạy nền hay chạy trước | Dịch vụ phải tự chuyển sang chạy nền (`service nginx start`, `busybox httpd -p …`, `cron`) hoặc được đưa ra nền trong script. Script không được kết thúc bằng một dịch vụ chạy ở tiền cảnh: script không thoát thì phiên bị hủy khi hết `PROVISION_TIMEOUT` |
| Nhật ký của dịch vụ | Không có `rsyslog` và không có trình gửi thư trong môi trường. Dịch vụ ghi nhật ký qua syslog hoặc cron gửi đầu ra qua thư đều mất dấu; đặc tả bài lab cần dịch vụ ghi vào tệp thì phải chuyển hướng tường minh trong script |
| `cron` | `cron` không tự chạy. Script phải tự khởi động (`/usr/sbin/cron`). `cron` đặt `PATH=/usr/bin:/bin` và không giữ biến môi trường của script, nên dòng trong `crontab` dùng đường dẫn tuyệt đối và tự chuyển hướng đầu ra vào tệp |
| Múi giờ | Môi trường chạy theo UTC. Bài lab có yếu tố thời gian (mốc giờ trong nhật ký, công việc định kỳ) phải tự đặt `TZ` trong script hoặc `CRON_TZ` trong `crontab`; giờ của môi trường không trùng giờ Việt Nam dùng ở giao diện và sổ giờ máy |
| Số tiến trình | Môi trường giới hạn 512 tác vụ, tính cả luồng và tiến trình đã kết thúc chưa được thu dọn. Script tạo tiến trình trong vòng lặp phải chờ chúng kết thúc |

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
| `runtime.session.resource-alert` | M1 | M3 | Mã phiên |
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
| Khung văn bản từ máy chủ | JSON: `{"type":"status","phase":"READY","expiresAt":"..."}`, `{"type":"expiring","expiresAt":"..."}`, `{"type":"warning","code":"IDLE","terminateAt":"..."}`, `{"type":"warning","code":"CPU_ABUSE","terminateAt":"..."}`, `{"type":"closing","reason":"EXPIRED"}` |
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
| Kết quả | Mã thoát 0 là đạt; khác 0 là không đạt; hết thời gian là không đạt và `timedOut = true`. Mã thoát 126 (tìm thấy nhưng không chạy được) và 127 (không tìm thấy lệnh) nghĩa là script chưa từng chạy: M1 trả `invalidCommand = true`, M2 không ghi bước là không đạt mà ghi lỗi đặc tả bài lab và báo tác giả; bài lab có lỗi này không được xuất bản (mục 4.5) |
| Giới hạn | Thời gian theo `timeoutSeconds`: tối đa 15 giây với `purpose = CHECK`, 60 giây với `purpose = RESTORE`; `stdout`, `stderr` cắt ở 4 KB; mỗi phiên chạy tối đa một lệnh tại một thời điểm |
| Thu nhận đầu ra | Chỉ đầu ra sinh ra trước khi tiến trình chính của script thoát được bảo đảm. Docker Engine đóng luồng sau một khoảng ngắn kể từ lúc tiến trình đó thoát, nên đầu ra của tiến trình nền sinh sau mốc đó bị mất trong khi mã thoát vẫn là 0. Script phải in xong kết luận trước khi thoát và không để lại tiến trình nền giữ `stdout`, `stderr`. Khác với `setup` (mục 4.7), lệnh kiểm tra không được chuyển hướng ra tệp trong môi trường: học viên có quyền `root` nên tệp đó sửa được, và kết quả kiểm tra phải không phụ thuộc trạng thái do học viên tạo ra |
| Dọn tiến trình | Sau khi lệnh trả về, M1 kết thúc cả nhóm tiến trình của lệnh (BR-M1.30), nên lệnh kiểm tra không để lại tiến trình chạy tiếp trong môi trường của học viên |
| Khôi phục | Script khôi phục là lời giải mẫu của các bước đã đạt, nối theo thứ tự bước, chạy với `purpose = RESTORE` bằng `/bin/sh` của môi trường với `PATH` chuẩn của hệ điều hành, vì mục đích là tạo trạng thái chứ không phải đánh giá. Kiểm thử bài lab (mục 4.5) chạy lời giải theo cùng cách |
| Phạm vi bảo đảm | Kết quả không bị ảnh hưởng khi học viên sửa hoặc thay công cụ của hệ thống, biến môi trường, tệp cấu hình shell, `/etc/ld.so.preload`. Lệnh kiểm tra không nên gọi chương trình nằm ngoài `/opt/crucible/bin` (ADR-0011) |

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 26/10/2026 | Ban hành lần đầu |
| 1.1 | 27/10/2026 | Mục 4.7 môi trường khi script bài lab chạy; mục 7 bổ sung mã 126, 127, phạm vi thu nhận đầu ra, dọn nhóm tiến trình |
