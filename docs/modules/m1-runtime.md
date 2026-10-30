# Thiết kế module M1 — Phiên thực hành và tài nguyên

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-MOD-M1 |
| Phiên bản | 1.3 |
| Ngày ban hành | 27/10/2026 |
| Người phụ trách | Trưởng nhóm (`NhatMinh0311`) |
| Tài liệu liên quan | CR-DOC-01 mục 5.1, 7.1; CR-DOC-02; CR-DOC-03; CR-DOC-04; CR-DOC-05; ADR-0002, ADR-0003, ADR-0006, ADR-0007, ADR-0011, ADR-0016 |

Các mục còn phụ thuộc kết quả thử nghiệm được đánh dấu (S1)–(S5) (CR-DOC-06 mục 3).

## 1. Phạm vi và trách nhiệm

MVP gồm R01–R13. R14 (môi trường nhiều máy) và R15 (bài lab Docker) để sau.

M1 chạy môi trường của phiên bằng container Docker trên một máy phiên (`cr-lab`) với runtime gVisor (`runsc`). Thành
phần điều khiển (`runtime-api`) chạy trên máy nền tảng (`cr-server`) và gọi Docker Engine của máy phiên qua SSH
(ADR-0002, ADR-0009).

## 2. Thành phần triển khai

Một Go module `services/runtime` với một chương trình `runtime-api`, chạy một bản:

| Phần | Thư viện chính | Trách nhiệm |
|---|---|---|
| API | `net/http`, oapi-codegen (sinh từ `contracts/runtime-api.yaml`), pgx, golang-migrate | API tại `contracts/runtime-api.yaml`; kiểm tra hạn mức; ghi phiên; cấp vé; chạy lệnh kiểm tra |
| Vòng phiên | — | Xếp hàng và cấp chỗ; hết hạn, không hoạt động; lấy mẫu tài nguyên; đối soát với máy phiên; dọn tài nguyên mồ côi |
| Trình điều khiển `docker` | Docker Engine API (`github.com/moby/moby/client`) | Dựng, exec, terminal, xóa container và mạng của phiên; kéo sẵn image; đọc số đo |

Luồng exec của Docker Engine khác nhau theo chế độ: với `Tty: false` (lệnh kiểm tra, script khôi phục) hai luồng
`stdout` và `stderr` được ghép kênh trên một kết nối, mỗi khung có tiêu đề 8 byte, nên phải tách bằng `stdcopy.StdCopy`;
với `Tty: true` (terminal) luồng là byte thô. Nhật ký của tiến trình chính lấy qua `ContainerLogs`, đi đường khác với đầu
ra của exec.
| Terminal | `coder/websocket` | Terminal WebSocket, người xem trực tiếp, ghi asciicast, tái dựng dòng lệnh |
| Sự kiện | `amqp091-go` | Chuyển tiếp outbox; nhận sự kiện của M2, M3 |
| Lưu bản ghi | AWS SDK for Go v2 (S3) | Tải bản ghi phiên lên SeaweedFS |

Tầng phiên làm việc với môi trường qua giao diện `driver.Driver`, không biết môi trường chạy ở đâu. Trên máy phát triển,
cùng chương trình gọi Docker Engine cục bộ với runtime `runc` (CR-DOC-03 mục 2.3).

## 3. Môi trường của phiên

### 3.1. Container phiên

| Thuộc tính | Giá trị |
|---|---|
| Tên | `crucible-session-<sessionId>` |
| Nhãn | `crucible.dev/session=<sessionId>`, `crucible.dev/managed-by=runtime-api` |
| Runtime | `runsc` trên môi trường chuẩn (`SESSION_RUNTIME`), `runc` trên máy phát triển |
| Lệnh | Theo image (`tini -- sleep infinity`); PID 1 là tiến trình khởi tạo do M1 chèn (`Init`, BR-M1.31); terminal và lệnh kiểm tra đều đi qua exec |
| Giới hạn | `Memory` bằng `MemorySwap` = bộ nhớ của mức tài nguyên + `SANDBOX_OVERHEAD_MEMORY` (giới hạn áp dụng cho cả sandbox gVisor, CR-DOC-04 mục 10); `NanoCPUs` theo mức tài nguyên; `PidsLimit` 512; lớp ghi của hệ thống tệp gốc tối đa 1 GiB do cấu hình runtime `runsc` của máy phiên (CR-DOC-04 mục 3.4) |
| Bảo mật | Bộ capability mặc định của Docker, `CapAdd` theo BR-M1.26; `no-new-privileges`; không gắn socket Docker, không gắn thư mục của máy |
| Bộ công cụ kiểm tra | Volume `crucible-tools-<phiên bản>` tạo một lần từ image `crucible-tools`, gắn `/opt/crucible/bin` chỉ đọc (ADR-0011) |
| Mạng | `none`, hoặc mạng riêng của phiên với `network: mirror` (mục 3.2) |
| Khởi động lại | Không (`RestartPolicy: no`): vượt giới hạn bộ nhớ hoặc số tiến trình làm dừng cả sandbox; nếu khởi động lại, container thành môi trường mới trống trong khi phiên vẫn tiếp tục (BR-M1.17) |
| Tên máy | `lab` |

### 3.2. Mạng của phiên

| `network` | Cách dựng |
|---|---|
| `none` | `NetworkMode: none`; container chỉ có giao diện loopback |
| `mirror` | Mạng Docker `internal` riêng của phiên `crucible-net-<sessionId>` (không có đường ra ngoài); container máy chủ gói nội bộ `crucible-mirror` được nối vào mạng này với tên `mirror`. Mạng bị xóa khi phiên kết thúc |

Mỗi phiên một mạng riêng nên hai phiên không thấy nhau. Dải địa chỉ của các mạng phiên lấy từ `default-address-pools`
của Docker trên máy phiên, mỗi mạng một khối `/28` (CR-DOC-04 mục 3.3).

### 3.3. Nguồn sự thật và đối soát

Bảng `sessions` trong CSDL `runtime` là nguồn sự thật về trạng thái mong muốn và pha của phiên. Docker trên máy phiên là
trạng thái thực tế. `runtime-api` đối soát hai bên khi khởi động và mỗi 60 giây:

| Tình huống | Xử lý |
|---|---|
| Container, mạng mang nhãn của Crucible nhưng không thuộc phiên chưa kết thúc | Xóa (BR-M1.12) |
| Phiên ở `READY`, `EXPIRING` mà container không còn | Phiên chuyển `TERMINATED` (BR-M1.12) |
| Phiên ở `PROVISIONING` khi `runtime-api` khởi động lại | Xóa container dở dang, dựng lại từ đầu trong thời hạn `PROVISION_TIMEOUT` |
| Container của phiên đã dừng | Xử lý theo BR-M1.17 |

## 4. Quy tắc nghiệp vụ

| Mã | Quy tắc | Liên quan |
|---|---|---|
| BR-M1.1 | Mỗi học viên có tối đa một phiên `PRACTICE` và một phiên `EXAM` ở pha chưa kết thúc. Bảo đảm bằng chỉ mục duy nhất có điều kiện trên `sessions(owner_user_id, priority)` | AC-R1.3 |
| BR-M1.2 | Phiên `PRACTICE` bị từ chối khi tài khoản hạn mức đã dùng hết `weeklyMinutes` của tuần hoặc đạt `maxConcurrentSessions`. Phiên `EXAM` không bị kiểm tra giờ máy. Phiên `LAB_CI` tính vào tài khoản `system:lab-ci` | AC-R2.4 |
| BR-M1.3 | Giờ máy của một phiên = số phút làm tròn lên từ `readyAt` tới `endedAt` (hoặc `lostAt`), bằng 0 nếu phiên chưa từng sẵn sàng; ghi một lần vào sổ khi phiên kết thúc. Tuần bắt đầu lúc 00:00 thứ Hai giờ Việt Nam | AC-R2.5 |
| BR-M1.4 | Phiên chờ được cấp theo thứ tự ưu tiên `EXAM` > `PRACTICE` > `LAB_CI`, cùng ưu tiên thì theo thời điểm tạo | AC-R2.1 |
| BR-M1.5 | Không dừng phiên đang chạy để lấy chỗ | AC-R2.1 |
| BR-M1.6 | Chỗ đặt trước được giữ từ `startsAt` − `reservation.leadMinutes` tới `endsAt`; trong khoảng này, phần chỗ đặt trước chưa dùng chỉ cấp cho phiên `EXAM` có cùng `reservationRef` | AC-R2.2 |
| BR-M1.7 | Đặt trước mới bị từ chối nếu làm tổng chỗ đặt trước của các đặt trước chồng lấn thời gian vượt `reservation.maxCapacityRatio` × sức chứa của mức tài nguyên; phần còn lại dành cho luyện tập. Đặt trước đã tạo được giữ khi tham số giảm | AC-L2.1 |
| BR-M1.8 | Phiên ở `PENDING` quá `queue.timeoutMinutes` chuyển `FAILED` với lý do `QUEUE_TIMEOUT` | AC-R2.3 |
| BR-M1.9 | Phiên `PRACTICE` không có dữ liệu nhập trong `idleMinutes` bị kết thúc với lý do `IDLE`; cảnh báo gửi qua terminal trước 2 phút. Không áp dụng cho `EXAM` | AC-R1.4 |
| BR-M1.10 | Chỉ phiên `PRACTICE` ở pha `READY`, `EXPIRING` được gia hạn; tổng gia hạn không vượt `maxExtensionMinutes` | AC-R1.5 |
| BR-M1.11 | Docker Engine của máy phiên không phản hồi liên tục quá `LAB_LOST_AFTER`: mọi phiên chưa kết thúc chuyển `LOST`, không nhận phiên mới cho tới khi kết nối lại. Khi kết nối lại, container của các phiên `LOST` bị xóa | AC-R4.1 |
| BR-M1.12 | Container, mạng mang nhãn `crucible.dev/managed-by=runtime-api` không thuộc phiên chưa kết thúc bị xóa trong vòng 2 phút; phiên `READY`, `EXPIRING` mà container không còn trên máy phiên được chuyển `TERMINATED` với lý do `ADMIN` và phát sự kiện | AC-R4.2 |
| BR-M1.13 | Vé terminal chỉ cấp cho chủ phiên ở pha `READY`, `EXPIRING`; hết hạn sau 30 giây; dùng một lần; chỉ lưu mã băm SHA-256 | AC-R3.1 |
| BR-M1.14 | Mỗi phiên có tối đa một kết nối tương tác; kết nối mới đóng kết nối cũ với mã 4009 | CR-DOC-03 mục 6 |
| BR-M1.15 | Khi `SESSION_RUNTIME=runsc`, `runtime-api` chỉ tạo container phiên với runtime `runsc` và không khởi động nếu Docker của máy phiên không có runtime này | ADR-0003 |
| BR-M1.16 | Mỗi phiên chạy tối đa một lệnh exec (`CHECK` hoặc `RESTORE`) tại một thời điểm (`EXEC_BUSY`); chỉ ở pha `READY`, `EXPIRING`. `RESTORE` chạy bằng shell và `PATH` của môi trường, `CHECK` chạy bằng bộ công cụ chỉ đọc | CR-DOC-03 mục 7 |
| BR-M1.17 | Container phiên dừng khi phiên chưa được yêu cầu kết thúc làm phiên chuyển `TERMINATED`: lý do `OUT_OF_MEMORY` khi Docker báo `OOMKilled`, `ENVIRONMENT_FAILED` với mọi trường hợp khác (vượt giới hạn số tiến trình, tiến trình chính kết thúc, lỗi sandbox). Container không được khởi động lại | AC-C1.8, AC-L2.11 |
| BR-M1.18 | Tổng `weekly_minutes` của các tài khoản con không vượt `weekly_minutes` của tài khoản cha; phần chưa chia dùng cho phiên tính trực tiếp vào tài khoản cha. Kiểm tra trong transaction có khóa dòng tài khoản cha | AC-R2.6 |
| BR-M1.19 | Phiên `PRACTICE` bị từ chối khi tổng phiên đang chạy của tài khoản cha và mọi tài khoản con đạt `max_concurrent_sessions` của tài khoản cha | AC-R2.6 |
| BR-M1.20 | Chỗ × giờ của đặt trước = `slots` × thời lượng tính theo giờ, làm tròn lên theo nửa giờ; tính vào tháng dương lịch (giờ Việt Nam) chứa `startsAt` của tài khoản gốc. Tổng trong tháng không vượt `monthly_reservation_slot_hours`. Hủy trước `startsAt` hoàn lại; kết thúc sớm do kỳ thi đóng không hoàn lại | AC-R2.7 |
| BR-M1.21 | Scope `runtime.quota.write` chỉ ghi tài khoản có `parent_ref`; tài khoản không có cha chỉ vai trò `ADMIN` ghi. Tài khoản cha chưa tồn tại được tạo với giá trị `quota.default*` hiện hành | AC-R2.8 |
| BR-M1.22 | Phiên `PRACTICE`, `EXAM` có CPU dùng trung bình trong `cpuAbuse.windowMinutes` từ `cpuAbuse.ratio` lần giới hạn CPU trở lên bị coi là dùng CPU bất thường. Không áp dụng cho phiên `LAB_CI` và phiên có `environment.cpuIntensive` | AC-R5.1, AC-R5.4 |
| BR-M1.23 | Phiên `PRACTICE` dùng CPU bất thường nhận khung `warning` mã `CPU_ABUSE` với `terminateAt` = thời điểm phát hiện + `cpuAbuse.graceMinutes` và sự kiện `resource-alert` (`WARNED`). Tới `terminateAt`, nếu CPU trung bình trong `cpuAbuse.graceMinutes` vừa qua vẫn từ ngưỡng trở lên, phiên chuyển `TERMINATED` với lý do `RESOURCE_ABUSE`; nếu không, cảnh báo được gỡ và lần phát hiện ghi `WARNED` | AC-R5.1 |
| BR-M1.24 | Phiên `EXAM` dùng CPU bất thường không bị kết thúc; phát `runtime.session.resource-alert` (`REPORTED`), tối đa một lần mỗi 10 phút cho một phiên | AC-R5.2, AC-L2.10 |
| BR-M1.25 | Mỗi lần phát hiện được ghi vào `resource_violations` với cách xử lý cuối cùng; không có biện pháp tự động theo số lần vi phạm | AC-R5.3 |
| BR-M1.26 | `environment.capabilities` chỉ chứa giá trị có trong `SESSION_ALLOWED_CAPABILITIES`, nếu không từ chối tạo phiên với 422 `ENVIRONMENT_NOT_ALLOWED`. Capability được thêm vào `CapAdd` của container phiên; phiên không khai báo không có capability bổ sung (ADR-0016) | AC-C1.9 |
| BR-M1.27 | Tham số hệ thống của M1 (CR-DOC-01 mục 5.6) đọc từ bảng `settings` qua bộ nhớ đệm làm mới mỗi 30 giây; chỉ `ADMIN` sửa, giá trị ngoài khoảng hợp lệ trả 422 `SETTING_OUT_OF_RANGE`; mỗi lần sửa ghi `settings_history` trong cùng transaction | AC-P1.1–AC-P1.4 |
| BR-M1.28 | Đổi số chỗ của đặt trước chưa hủy và chưa qua `endsAt`: phần tăng được kiểm tra theo BR-M1.7 và BR-M1.20 như đặt trước mới, bị từ chối thì số chỗ giữ nguyên; giảm không xuống dưới số chỗ đang dùng; giảm trước `startsAt` hoàn lại chỗ × giờ của phần giảm | AC-L2.13 |
| BR-M1.29 | `environment.setup` chạy một lần khi dựng môi trường, bằng `/bin/sh` của image, với `stdin` từ `/dev/null` và `stdout`, `stderr` vào `/var/log/crucible/setup.log` trong môi trường. Dịch vụ nền do script khởi động thừa hưởng ba bộ mô tả tệp này nên không giữ luồng của lệnh exec. Script không thoát trong `PROVISION_TIMEOUT` làm phiên chuyển `FAILED` với lý do `PROVISION_TIMEOUT`; nội dung tệp nhật ký được đọc tối đa 4 KiB và ghi vào bản ghi phiên | CR-DOC-03 mục 4.7 |
| BR-M1.30 | Lệnh exec trả mã thoát 126 hoặc 127 được trả về với `invalidCommand = true`. Sau khi mọi lệnh exec trả về hoặc hết thời hạn, M1 kết thúc cả nhóm tiến trình của lệnh trong môi trường; hủy `context` của bên gọi không tự dừng tiến trình trong container | CR-DOC-03 mục 7 |
| BR-M1.31 | Container phiên được tạo với tiến trình khởi tạo của Docker (`HostConfig.Init`) làm PID 1, để tiến trình mồ côi luôn được thu dọn mà không phụ thuộc `ENTRYPOINT` của image bài lab. Tiến trình đã kết thúc chưa được thu dọn vẫn chiếm một suất của `PidsLimit` cho tới khi hủy phiên | CR-DOC-03 mục 4.7, mục 7 |

## 5. Cấp phát

### 5.1. Mô hình sức chứa

- Chi phí một phiên = bộ nhớ của mức tài nguyên + chi phí gVisor (CR-DOC-04 mục 10). Bộ nhớ là tài nguyên giới hạn;
  CPU không được dùng để tính chỗ.
- Bộ nhớ dành cho phiên = tổng bộ nhớ của máy phiên (`docker info`) − `LAB_RESERVED_MEMORY` (hệ điều hành, Docker, máy
  chủ gói nội bộ).
- Chỗ trống = (bộ nhớ dành cho phiên − tổng chi phí của các phiên `PROVISIONING`, `READY`, `EXPIRING`) / chi phí phiên.
- Sức chứa của một mức tài nguyên = bộ nhớ dành cho phiên / chi phí phiên của mức đó; bằng 0 khi máy phiên mất kết nối
  (BR-M1.11).
- `GET /v1/capacity` trả sức chứa hiện tại trừ chỗ đặt trước chồng lấn khung giờ hỏi.
- Đặt trước bị từ chối trả `capacity.requiredSlots` và `capacity.availableSlots`, trong đó số chỗ còn là phần có thể đặt
  trước theo BR-M1.7 sau khi trừ các đặt trước chồng lấn.
- `GET /v1/capacity/windows` duyệt các mốc bắt đầu cách nhau 30 phút từ `from`, bỏ mốc làm khung giờ vượt ra ngoài
  `dailyStart`–`dailyEnd`, tính số chỗ còn như trên và trả `limit` khung giờ đầu tiên đủ `slots`; dừng sau 14 ngày.

Kiểm tra sức chứa và ghi đặt trước chạy trong một transaction có khóa dòng (`SELECT ... FOR UPDATE`) trên các đặt trước
chồng lấn, để hai yêu cầu đồng thời không cùng vượt giới hạn.

### 5.2. Vòng xếp hàng

`runtime-api` chạy vòng xếp hàng mỗi khi có phiên mới, phiên kết thúc, và định kỳ mỗi 5 giây:

1. Lấy các phiên ở `PENDING`, sắp theo BR-M1.4.
2. Với mỗi phiên: tính chỗ được phép dùng theo BR-M1.6.
3. Nếu còn chỗ: chuyển `PROVISIONING`, ghi `admittedAt`, dựng môi trường (mục 5.3).
4. Nếu không: giữ `PENDING` và dừng xét các phiên cùng ưu tiên sau nó.

Vị trí trong hàng đợi trả về ở `getSession` được tính từ danh sách `PENDING`.

### 5.3. Dựng môi trường

```mermaid
sequenceDiagram
    autonumber
    participant A as runtime-api
    participant D as Docker Engine (cr-lab)
    participant C as Container phiên

    opt network = mirror
        A->>D: Tạo mạng internal crucible-net-<id>, nối crucible-mirror
    end
    A->>D: Tạo container (runsc, giới hạn, nhãn, bộ công cụ chỉ đọc)
    A->>D: Khởi động container
    A->>C: exec: /opt/crucible/bin/tar -x --no-same-owner -C / (luồng tar các tệp của biến thể)
    A->>C: exec: chạy setup (thời hạn SETUP_TIMEOUT)
    alt Thành công
        A->>A: READY, readyAt, expiresAt, sự kiện vào outbox trong cùng transaction
    else Lỗi hoặc quá hạn
        A->>D: Xóa container và mạng
        A->>A: FAILED, PROVISION_FAILED
    end
```

Tệp của biến thể được ghi qua exec bằng `tar` của bộ công cụ, không dùng `CopyToContainer`, `CopyFromContainer` của
Docker: gVisor giữ lớp ghi của hệ thống tệp gốc bên trong sandbox, tệp sao chép từ máy vào không hiện ra bên trong và
tệp tạo bên trong không đọc được từ máy (S1). `tar` giải nén với `--no-same-owner` để tệp thuộc người dùng chạy lệnh.

Image của bài lab được kéo sẵn lên máy phiên khi bài lab được xuất bản và khi đặt trước được tạo (mục 8), nên thời gian
dựng thường chỉ gồm khởi động sandbox và chạy setup (S2).

### 5.4. Đồng bộ trạng thái với CSDL và sự kiện

Mỗi lần đổi pha, `runtime-api` cập nhật bảng `sessions` và ghi sự kiện vào `outbox` trong cùng một transaction, với
`sequence` tăng dần theo phiên; `sequence` bảo đảm bên nhận bỏ qua sự kiện trùng. Thao tác trên Docker làm trước, đổi
pha sau: nếu `runtime-api` dừng giữa hai bước, lần đối soát sau (mục 3.3) xử lý phần còn dở.

### 5.5. Đo tài nguyên của phiên

`runtime-api` lấy mẫu số đo của container phiên mỗi 15 giây qua Docker Engine API (thống kê một lần, không giữ luồng)
và giữ bộ nhớ cao nhất của phiên. Container dừng với `OOMKilled` được ghi thành `oomKilled` (BR-M1.17); trường hợp dừng
khác được ghi mã thoát vào lịch sử phiên. Khi phiên kết thúc, giá trị được ghi vào bảng `sessions`, trả trong
`resourceUsage` của `getSession` và sự kiện `runtime.session.terminated`. Độ chính xác của số đo dưới gVisor (phần chi
phí của sandbox có bị tính vào hay không) được xác định ở thử nghiệm S5.

Cùng lần lấy mẫu, `runtime-api` đọc tổng thời gian CPU của container và tính tỉ lệ CPU dùng so với giới hạn CPU của mức
tài nguyên cho từng khoảng 15 giây. Trung bình trượt 10 phút và 5 phút của tỉ lệ này là đầu vào của BR-M1.22, BR-M1.23.
Phiên mới sẵn sàng chưa đủ 10 phút dữ liệu thì chưa được xét.

## 6. Terminal

| Mục | Thiết kế |
|---|---|
| Kết nối | `runtime-api` tiêu thụ vé (`UPDATE ... SET used_at = now() WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now() RETURNING ...`), mở exec có TTY chạy `script -qc "bash -l" /dev/null` tới container phiên qua Docker Engine API: PTY của exec dưới gVisor không có tên `/dev/pts/N`, `script` tạo PTY bên trong sandbox để `tty`, `gpg`, `screen` hoạt động (S1) |
| Một kết nối tương tác | Giữ trong bộ nhớ theo mã phiên (một bản `runtime-api`); kết nối mới đóng kết nối cũ với mã 4009 (BR-M1.14) |
| Xem trực tiếp | Đầu ra của kết nối tương tác được phát cho các kết nối xem trong cùng tiến trình. Khi mới kết nối, người xem nhận 64 KB đầu ra gần nhất. Không dùng `tmux` trong môi trường vì học viên có thể thoát hoặc dừng tiến trình đó |
| Shell thoát | Học viên gõ `exit`: mở shell mới trong cùng kết nối |
| Ghi phiên | Đầu ra (và đầu vào với phiên `recordInput`) ghi theo asciicast v2; đệm 64 KB hoặc 10 giây rồi tải lên SeaweedFS thành phần `recordings/<sessionId>/<số thứ tự>.part`; khi phiên kết thúc, ghép thành `recordings/<sessionId>.cast` và phát `runtime.recording.available` |
| Tái dựng dòng lệnh | Bộ đệm theo kết nối: ký tự in được thêm vào, `0x7f`/`0x08` xóa ký tự cuối, `Ctrl-U` xóa dòng, `\r` phát dòng; chuỗi escape bị bỏ qua. Khung nhập dài hơn 200 byte đánh dấu `pasted`. Tối đa 5 sự kiện mỗi giây mỗi phiên |
| Hoạt động | Mỗi khung nhập cập nhật thời điểm hoạt động trong bộ nhớ; ghi vào `sessions.last_activity_at` tối đa mỗi 30 giây |
| Dừng có kiểm soát | Nhận SIGTERM: ngừng nhận kết nối mới, gửi `closing` mã 4010 cho mọi kết nối, đẩy phần bản ghi còn đệm, thoát sau tối đa 20 giây. Giao diện lấy vé mới và kết nối lại; môi trường của phiên không bị ảnh hưởng |

## 7. Mô hình dữ liệu (CSDL `runtime`)

```mermaid
erDiagram
    quota_accounts ||--o{ sessions : "quota_account_ref"
    quota_accounts ||--o{ quota_accounts : "parent_ref"
    quota_accounts ||--o{ reservations : "quota_account_ref"
    quota_accounts ||--o{ usage_ledger : ""
    sessions ||--o| usage_ledger : ""
    sessions ||--o{ terminal_tickets : ""
    sessions ||--o{ view_access_log : ""
    sessions ||--o{ resource_violations : ""
    reservations ||--o{ sessions : "reservation_ref"
    settings ||--o{ settings_history : ""

    sessions {
        uuid id PK
        uuid owner_user_id
        text priority
        text profile
        text image
        text network
        text quota_account_ref
        text reservation_ref
        jsonb labels
        text phase
        text end_reason
        bigint last_sequence
        timestamptz created_at
        timestamptz admitted_at
        timestamptz ready_at
        timestamptz expires_at
        timestamptz ended_at
        timestamptz last_activity_at
        int used_minutes
        text container_id
        bigint memory_peak_bytes
        bool oom_killed
    }
    quota_accounts {
        text ref PK
        text parent_ref FK
        int weekly_minutes
        int max_concurrent_sessions
        int monthly_reservation_slot_hours
    }
    usage_ledger {
        bigint id PK
        text quota_account_ref FK
        uuid session_id UK
        date week_start
        int minutes
        text priority
    }
    reservations {
        uuid id PK
        text reservation_ref UK
        text quota_account_ref FK
        numeric slot_hours
        text profile
        int slots
        timestamptz starts_at
        timestamptz ends_at
        text[] prepull_images
        timestamptz cancelled_at
    }
    terminal_tickets {
        bytea token_hash PK
        uuid session_id FK
        text mode
        uuid user_id
        timestamptz expires_at
        timestamptz used_at
    }
    view_access_log {
        bigint id PK
        uuid session_id FK
        uuid viewer_user_id
        timestamptz issued_at
    }
    resource_violations {
        uuid id PK
        uuid session_id FK
        uuid owner_user_id
        text priority
        timestamptz detected_at
        numeric cpu_ratio
        text action
    }
    settings {
        text key PK
        numeric value
        timestamptz updated_at
        uuid updated_by
    }
    settings_history {
        bigint id PK
        text key FK
        numeric old_value
        numeric new_value
        timestamptz changed_at
        uuid changed_by
    }
```

`reservations.prepull_images` lưu trường `warmImages` của hợp đồng. Bảng `settings` chỉ chứa tham số đã được sửa; tham
số chưa sửa dùng mặc định trong mã, cùng với khoảng hợp lệ và mô tả.

Bảng hạ tầng: `idempotency_keys`, `outbox`, `inbox` (CR-DOC-03 mục 3.3, 5.2).

## 8. Sự kiện

| Sự kiện | Phát / nhận | Xử lý |
|---|---|---|
| `runtime.session.*`, `runtime.terminal.*`, `runtime.recording.available` | Phát | Theo mục 5.4 và mục 6; `runtime.session.resource-alert` theo BR-M1.23, BR-M1.24 |
| `content.lab.published` | Nhận | Kéo sẵn image của bài lab lên máy phiên |
| `learning.exam.closed` | Nhận | Kết thúc sớm đặt trước có `reservationRef` tương ứng |

Tạo đặt trước cũng kéo sẵn các image trong `warmImages`.

## 9. Giao diện

| Màn hình, thành phần | Vai trò | Nội dung |
|---|---|---|
| Thành phần `<Terminal sessionId mode>` | Học viên, giảng viên | xterm.js; lấy vé, kết nối, tự kết nối lại theo mã đóng; được M2 nhúng vào trang làm bài, M3 nhúng vào trang giám sát |
| Thanh trạng thái phiên | Học viên | Pha, thời gian còn lại, vị trí hàng đợi, gia hạn, cảnh báo không hoạt động, cảnh báo dùng CPU bất thường |
| Quản trị phiên và sức chứa | `ADMIN` | Sức chứa của máy phiên theo mức tài nguyên, trạng thái kết nối máy phiên, hàng đợi, phiên đang chạy (lọc, kết thúc, sắp theo CPU), đặt trước theo lịch, trần hạn mức của giảng viên, lần phát hiện dùng CPU bất thường theo người dùng |

## 10. Tác vụ định kỳ

Mọi tác vụ chạy trong `runtime-api`.

| Tác vụ | Chu kỳ | Nội dung |
|---|---|---|
| Vòng xếp hàng | 5 giây và theo sự kiện | Mục 5.2 |
| Hết hạn và không hoạt động | 15 giây | Chuyển `EXPIRING`, `TERMINATED`; cảnh báo không hoạt động |
| Lấy mẫu tài nguyên | 15 giây | Bộ nhớ cao nhất, tỉ lệ CPU; xét BR-M1.22–BR-M1.24 (mục 5.5) |
| Theo dõi máy phiên | 10 giây | Gọi `ping` Docker Engine; xét BR-M1.11; nhận sự kiện container dừng (BR-M1.17) |
| Đối soát và dọn tài nguyên mồ côi | 60 giây | Mục 3.3, BR-M1.12; cập nhật `crucible_orphan_resources` |
| Chuyển tiếp outbox | 1 giây | CR-DOC-03 mục 5.2 |
| Dọn vé, khóa lặp lại an toàn, outbox, inbox | 1 giờ | Xóa bản ghi quá hạn |

## 11. Kiểm thử

| Mức | Công cụ | Phạm vi |
|---|---|---|
| Đơn vị | `go test` | Tính chỗ, thứ tự xếp hàng, quy tắc đặt trước, tính giờ máy, tái dựng dòng lệnh |
| Tích hợp | Testcontainers (PostgreSQL, RabbitMQ, SeaweedFS), Docker thật (runc) | Luồng HTTP và WebSocket, outbox, vé, hạn mức và đặt trước đồng thời, đối soát, mất kết nối máy phiên (trình điều khiển giả lập) |
| Hợp đồng | Schemathesis | Khớp `contracts/runtime-api.yaml` |
| Môi trường chuẩn | Bộ kiểm thử cách ly, script đo | I1–I12 với `runsc`; thời gian sẵn sàng, độ trễ terminal (Q1–Q3) |

| Tiêu chí, quy tắc | Kiểm thử |
|---|---|
| AC-R1.1, AC-R1.2 | Script đo thời gian sẵn sàng trên môi trường chuẩn (S2) |
| AC-R1.3, BR-M1.1 | Tích hợp: 10 yêu cầu đồng thời của một học viên, đúng một thành công |
| AC-R2.1, BR-M1.4, BR-M1.5 | Đơn vị: thứ tự cấp, không dừng phiên đang chạy |
| AC-R2.2, BR-M1.6, BR-M1.7 | Đơn vị và tích hợp: đặt trước; hai đặt trước đồng thời không cùng vượt giới hạn |
| AC-L2.13, BR-M1.28 | Tích hợp: tăng chỗ vượt sức chứa hoặc ngân sách bị từ chối và giữ nguyên; giảm dưới số chỗ đang dùng bị từ chối; giảm trước giờ bắt đầu hoàn lại ngân sách |
| AC-R2.4, AC-R2.5, BR-M1.2, BR-M1.3 | Tích hợp: sổ hạn mức |
| AC-R2.6–AC-R2.8, BR-M1.18–BR-M1.21 | Tích hợp: chia hạn mức đồng thời cho hai lớp cùng giảng viên không vượt trần; ngân sách đặt trước theo tháng; quyền ghi theo scope |
| AC-R3.1, BR-M1.13 | Tích hợp: dùng lại vé, vé hết hạn, vé của người khác |
| AC-R4.1, BR-M1.11 | Tích hợp với trình điều khiển giả lập mất kết nối; môi trường chuẩn: dừng Docker Engine trên `cr-lab` (thí nghiệm CH1) |
| AC-R4.2, BR-M1.12 | Tích hợp: container mồ côi bị xóa, container bị xóa ngoài `runtime-api` |
| AC-R5.1–AC-R5.4, BR-M1.22–BR-M1.25 | Đơn vị: trung bình trượt và chuyển trạng thái cảnh báo với chuỗi mẫu giả lập; môi trường chuẩn: kiểm thử cách ly I10 |
| AC-C1.9, BR-M1.26 | Đơn vị: dựng cấu hình container theo capability; tích hợp: capability ngoài danh sách bị từ chối; môi trường chuẩn: `iptables` chạy được khi khai báo, kiểm thử cách ly I12 |
| BR-M1.15 | Tích hợp: `runtime-api` không khởi động khi Docker không có runtime `runsc` mà cấu hình yêu cầu |
| AC-P1.1–AC-P1.4, BR-M1.27 | Tích hợp: sửa tham số bằng vai trò khác `ADMIN` bị từ chối; giá trị ngoài khoảng bị từ chối; sau khi giảm `queue.timeoutMinutes`, phiên chờ mới bị hủy theo giá trị mới trong vòng 60 giây; lịch sử ghi đủ giá trị cũ và mới |
| BR-M1.29 | Tích hợp: `setup` khởi động dịch vụ nền không chuyển hướng vẫn làm lệnh dựng trả về đúng hạn và dịch vụ còn chạy; `setup` kết thúc bằng dịch vụ ở tiền cảnh làm phiên `FAILED` với lý do `PROVISION_TIMEOUT` |
| BR-M1.30 | Tích hợp: lệnh gọi chương trình không có trong bộ công cụ trả `invalidCommand = true` và mã thoát 127; lệnh để lại tiến trình nền không còn tiến trình nào sau khi trả về; đầu ra in trước khi thoát không bị cắt |
| BR-M1.31 | Tích hợp: `/proc/1/comm` trong môi trường là tiến trình khởi tạo của Docker, không phải `ENTRYPOINT` của image; sau khi lệnh exec để lại tiến trình nền kết thúc, môi trường không còn tiến trình chưa được thu dọn |

## 12. Cấu hình

Cấu hình triển khai qua biến môi trường. Tham số nghiệp vụ (`queue.*`, `reservation.*`, `quota.default*`,
`cpuAbuse.*`) là tham số hệ thống do quản trị viên sửa qua `/v1/settings` (BR-M1.27, CR-DOC-01 mục 5.6).

| Biến | Ý nghĩa | Mặc định |
|---|---|---|
| `RUNTIME_DRIVER` | Trình điều khiển môi trường; hiện chỉ có `docker` | `docker` |
| `DOCKER_HOST` | Docker Engine của máy phiên; môi trường chuẩn dùng `ssh://crucible@cr-lab`. Moby client không hỗ trợ `ssh://`; với lược đồ này `runtime-api` tự kết nối bằng `ssh -T -- <đích> docker system dial-stdio` (S1) | Socket cục bộ |
| `SESSION_RUNTIME` | Runtime của container phiên (BR-M1.15) | `runc` |
| `LAB_RESERVED_MEMORY` | Bộ nhớ của máy phiên không dùng cho phiên (mục 5.1) | `2GiB` |
| `SANDBOX_OVERHEAD_MEMORY` | Bộ nhớ cộng thêm vào giới hạn của container cho sandbox gVisor (mục 3.1) | `0` (môi trường chuẩn: `64MiB`) |
| `LAB_LOST_AFTER` | Thời gian Docker Engine không phản hồi trước khi coi máy phiên mất kết nối (BR-M1.11); cộng chu kỳ theo dõi 10 giây không vượt 60 giây của AC-R4.1 | `45s` |
| `MIRROR_CONTAINER` | Tên container máy chủ gói nội bộ trên máy phiên; rỗng thì phiên `mirror` dùng mạng `bridge` (máy phát triển) | Rỗng |
| `DB_URL` | Kết nối CSDL `runtime` | `postgres://runtime:runtime-dev@localhost:5432/runtime` |
| `AMQP_URL` | Kết nối RabbitMQ | `amqp://crucible:rabbitmq-dev@localhost:5672/` |
| `OIDC_ISSUER`, `OIDC_AUDIENCE` | Xác thực token | `http://localhost:8180/realms/crucible`, `runtime-api` |
| `OIDC_JWKS_URL` | Địa chỉ JWKS khi dịch vụ không tới được Keycloak qua địa chỉ issuer (chạy trong container) | `<OIDC_ISSUER>/protocol/openid-connect/certs` |
| `S3_ENDPOINT`, `S3_BUCKET` | Lưu bản ghi | `http://localhost:8333`, `recordings` |
| `SETUP_TIMEOUT` | Thời hạn setup | `60s` |
| `PROVISION_TIMEOUT` | Thời hạn dựng môi trường, gồm kéo image và setup | `5m` |
| `SESSION_ALLOWED_CAPABILITIES` | Capability được cấp cho bài lab khai báo (BR-M1.26); chỉ đặt khi phiên chạy `runsc` | Rỗng (môi trường chuẩn: `NET_ADMIN`) |
| `SESSION_IMAGE_PREFIXES` | Tiền tố image bài lab được phép; image khác bị từ chối với 422 `ENVIRONMENT_NOT_ALLOWED` | `ghcr.io/crucible-project/crucible-labs/` |
| `TOOLS_IMAGE` | Image bộ công cụ kiểm tra | `ghcr.io/crucible-project/crucible/tools:1.0.0` |
| `PORT` | Cổng HTTP | `8081` |
| `EVENTS_EXCHANGE`, `EVENTS_SOURCE` | Exchange và `source` của sự kiện | `crucible.events`, `/crucible/runtime-api` |
| `TERMINAL_PUBLIC_URL` | Địa chỉ WebSocket trả trong vé terminal | `ws://localhost:8081/v1/terminal` |
| `TERMINAL_ALLOWED_ORIGINS` | Origin của giao diện được mở terminal | `localhost:*,127.0.0.1:*` |

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 26/10/2026 | Thiết kế phiên chạy bằng container Docker với gVisor trên máy phiên riêng, một bản `runtime-api` gồm API, vòng phiên và terminal |
| 1.1 | 27/10/2026 | Kết quả thử nghiệm S1: kết nối SSH do `runtime-api` thực hiện, shell của terminal qua `script`, `tar --no-same-owner`, giới hạn lớp ghi |
| 1.2 | 27/10/2026 | BR-M1.29 lớp bọc bộ mô tả tệp cho `setup`; BR-M1.30 mã thoát 126, 127 và dọn nhóm tiến trình của lệnh exec; ghi chú luồng exec ghép kênh |
| 1.3 | 28/10/2026 | BR-M1.31 tiến trình khởi tạo của container do M1 đặt |
