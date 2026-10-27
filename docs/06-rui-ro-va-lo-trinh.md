# Rủi ro, thử nghiệm sớm và lộ trình

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-DOC-06 |
| Phiên bản | 1.2 |
| Ngày ban hành | 26/10/2026 |
| Tài liệu liên quan | CR-DOC-01 đến CR-DOC-05, CR-DOC-07 Quy trình phát triển |

## 1. Phạm vi

Tài liệu liệt kê các rủi ro kỹ thuật và tổ chức, các thử nghiệm sớm để xác nhận những giả định kỹ thuật chưa được kiểm
chứng, lộ trình theo mốc, phân công và tiêu chí hoàn thành của từng module.

## 2. Rủi ro

| Mã | Rủi ro | Khả năng | Ảnh hưởng | Biện pháp | Phương án dự phòng |
|---|---|---|---|---|---|
| RK1 | gVisor không chạy ổn định với Docker trong máy ảo KVM, hoặc exec, terminal, giới hạn tài nguyên không hoạt động đúng | Thấp | Cao | Thử nghiệm S1 trong tuần đầu | Dùng runc kèm `userns-remap`, seccomp và bỏ capability, không cấp `NET_ADMIN` cho bài lab (ADR-0016); ghi nhận mức cách ly thấp hơn tại CR-DOC-05 |
| RK2 | Thời gian dựng môi trường mới (khởi động sandbox, ghi tệp, setup) vượt 5 giây | Trung bình | Trung bình | Thử nghiệm S2; kéo sẵn image (ADR-0007) | Bổ sung nhóm môi trường ấm cố định cho image của kỳ thi; rút gọn script setup của bài lab; nới mục tiêu Q1 có căn cứ đo |
| RK3 | Terminal qua Docker exec trên kết nối SSH có độ trễ cao hoặc mất kết nối | Thấp | Cao | Thử nghiệm S3 | Mở Docker Engine qua TCP có TLS chỉ cho `cr-server` thay cho SSH |
| RK4 | Phần mềm trong bài lab (dịch vụ không có systemd, công cụ mạng) hoạt động khác thường trong gVisor | Trung bình | Trung bình | Thử nghiệm S4 với image nền và bốn bài mẫu | Điều chỉnh image nền, loại bỏ dạng bài không tương thích khỏi phạm vi |
| RK5 | Máy chủ vật lý không đủ bộ nhớ cho toàn bộ thành phần | Thấp | Trung bình | Thử nghiệm S5 | Giảm thời gian lưu giữ giám sát, giảm bộ nhớ của `cr-lab` xuống 12 GB (CR-DOC-04 mục 10) |
| RK6 | Thành viên chậm tiến độ làm module khác bị chặn | Trung bình | Cao | Hợp đồng chốt trước khi code; mock và `runtime-api` chạy trên máy phát triển; mốc có tiêu chí cụ thể (mục 4) | Module bị chậm được rút gọn về chức năng MVP; module khác tiếp tục làm việc với mock |
| RK7 | Thời gian học công nghệ mới (Go, Docker, gVisor) của người phụ trách M1 | Trung bình | Trung bình | Kiến trúc một dịch vụ, một Docker Engine (ADR-0002); thử nghiệm sớm S1–S5 do chính người phụ trách thực hiện | Giảm phạm vi M1 về chức năng MVP |
| RK8 | Mã do công cụ AI sinh ra sai nghiệp vụ mà không bị phát hiện | Cao | Cao | Tiêu chí chấp nhận có mã và được kiểm thử; review bắt buộc; kiểm thử hợp đồng trong CI | |
| RK9 | Chạy Docker bên trong môi trường học viên không khả thi trong gVisor | Trung bình | Thấp | Đưa bài lab Docker ra ngoài MVP; thử nghiệm S6 khi còn thời gian | Không hỗ trợ bài lab Docker |
| RK10 | gVisor không cho tạo namespace mạng, `veth` trong user namespace con, hoặc cần `SYS_ADMIN` cho container | Trung bình | Thấp | Bài định tuyến nhiều mạng con là phạm vi mở rộng; thử nghiệm S7 | Không cấp `SYS_ADMIN`; giữ bài mạng trong một mạng con |

## 3. Thử nghiệm sớm

Mỗi thử nghiệm có thời hạn cố định và tiêu chí kết thúc. Kết quả được ghi tại `docs/spikes/<mã>.md` gồm cách làm,
số đo, kết luận và quyết định tiếp theo. Mã thử nghiệm được bỏ đi sau khi có kết luận; phần dùng lại được viết lại
theo quy trình bình thường.

| Mã | Câu hỏi cần trả lời | Thời hạn | Tiêu chí kết thúc |
|---|---|---|---|
| S1 | gVisor chạy được với Docker trên `cr-lab` trong máy ảo KVM; `runtime-api` điều khiển qua `DOCKER_HOST=ssh://`; exec có TTY, giới hạn bộ nhớ, CPU, số tiến trình, mạng `internal` riêng hoạt động với `runsc`; cách gVisor lưu lớp ghi của hệ thống tệp; tệp ghi qua exec `tar` hiện ra trong môi trường | 2 ngày | Container `runsc` chạy bash, exec và terminal thành công qua SSH; các kiểm thử I1–I9, I13 tại CR-DOC-05 cho kết quả mong đợi; ghi rõ hành vi khi vượt bộ nhớ và số tiến trình |
| S2 | Thời gian từ khi được cấp chỗ tới khi sẵn sàng (tạo container, ghi tệp, chạy setup) khi image đã có trên máy phiên | 1 ngày | Số đo p50, p95 trên 50 lần với bài `linux-user-nginx`; kết luận có cần nhóm môi trường ấm (RK2) |
| S3 | Terminal của `runtime-api` qua Docker exec trên kết nối SSH và ghi asciicast | 2 ngày | Gõ lệnh tương tác qua trình duyệt với xterm.js; độ trễ khứ hồi p95 đo được với 20 kết nối đồng thời; bản ghi phát lại được |
| S4 | Image nền và bốn bài mẫu chạy đúng trong gVisor; các cơ chế mẫu kịch bản cần: `cron`, tiến trình nền tách khỏi phiên, `/proc/<PID>/fd`, `netstat -ltnp` hiện PID; lớp bọc bộ mô tả tệp của `setup` (BR-M1.29) với dịch vụ tự chuyển nền và dịch vụ chạy ở tiền cảnh; đầu ra của lệnh kiểm tra bị cắt khi tiến trình nền giữ luồng; với `NET_ADMIN` (ADR-0016): `iptables`, `nftables`, `ip addr`, `ip route`; `ping` khi không có `NET_RAW`; số tiến trình của dịch vụ đa luồng so với giới hạn 512 của container (CR-DOC-04 mục 3.4) | 3 ngày | Lời giải mẫu đạt và môi trường mới dựng không đạt với 5 seed, chạy thủ công; ghi rõ cơ chế nào không dùng được |
| S5 | Toàn bộ thành phần nền tảng vừa bộ nhớ của `cr-server`; chi phí bộ nhớ của sandbox gVisor | 1 ngày | Số đo bộ nhớ thực tế của từng thành phần sau khi triển khai; giá trị `SANDBOX_OVERHEAD_MEMORY` |
| S6 | Docker Engine chạy bên trong container gVisor | 1 ngày | Kết luận khả thi hay không |
| S7 | Mô phỏng nhiều mạng con và router trong một container gVisor: không cấp `SYS_ADMIN`, `setup` tạo user namespace con chứa các namespace mạng có tên, nối bằng `veth` và `bridge`; tuyến tĩnh, `ip_forward` riêng từng namespace; bảng `filter`, `nat` của `iptables` trong namespace router; học viên và lệnh kiểm tra vào topo bằng `nsenter` | 1 ngày | Topo `client – r1 – r2 – server` thông khi tuyến đúng, không thông khi tuyến sai; `/opt/crucible/bin` không bị gắn đè từ trong topo; I1–I13 vẫn đạt; ghi rõ lệnh nào không dùng được |

S1–S5 do người phụ trách M1 thực hiện trong tuần 1–5, mỗi thử nghiệm xong trước việc dùng kết luận của nó: S1, S4
trước khi đưa gVisor vào `runtime-api`; S3 trước ghi phiên; S2, S5 trước đặt trước tài nguyên. Trong thời gian đó, hai
module còn lại làm việc với `runtime-api` chạy phiên bằng Docker trên máy phát triển. S6, S7 thực hiện khi còn thời
gian sau mốc MC.

## 4. Lộ trình

Tuần 1 là tuần đề tài được duyệt (26/10–01/11/2026); tuần 11 kết thúc ngày 10/01/2027. Mỗi tuần tính từ thứ Hai
đến Chủ nhật.

| Hạng mục | T1 | T2 | T3 | T4 | T5 | T6 | T7 | T8 | T9 | T10 | T11 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| Chung: nhận quyền truy cập, môi trường phát triển, Spec Kit | ■ |  |  |  |  |  |  |  |  |  |  |
| M1: thử nghiệm S1–S5 | ■ | ■ | ■ | ■ | ■ |  |  |  |  |  |  |
| M1: gVisor, mạng của phiên, máy chủ gói nội bộ |  |  | ■ | ■ |  |  |  |  |  |  |  |
| M1: hạ tầng hai máy ảo, script triển khai |  |  | ■ | ■ |  |  |  |  |  |  |  |
| M1: ghi phiên, giám sát, sao lưu |  |  |  |  | ■ |  |  |  |  |  |  |
| M1: kiểm thử cách ly; hàng đợi, hạn mức, sức chứa, đặt trước |  |  |  |  |  | ■ |  |  |  |  |  |
| M1: tham số hệ thống, trần hạn mức, gợi ý khung giờ |  |  |  |  |  |  | ■ |  |  |  |  |
| M1: lấy mẫu tài nguyên, mất máy phiên, đối soát |  |  |  |  |  |  |  | ■ |  |  |  |
| M1: trang quản trị, tham số hệ thống của ba module |  |  |  |  |  |  |  |  | ■ |  |  |
| M1: đo, thí nghiệm sự cố, khôi phục, phát hành |  |  |  |  |  |  |  |  |  | ■ | ■ |
| M2: tài liệu module | ■ | ■ |  |  |  |  |  |  |  |  |  |
| M2: lưu trữ bài lab, phiên bản, biến thể |  | ■ | ■ |  |  |  |  |  |  |  |  |
| M2: lượt làm bài, kiểm tra bước, gợi ý |  |  | ■ |  |  |  |  |  |  |  |  |
| M2: bước trả lời, form soạn bài |  |  |  | ■ | ■ |  |  |  |  |  |  |
| M2: ngôn ngữ nội dung bài lab |  |  |  |  |  | ■ |  |  |  |  |  |
| M2: kiểm thử bài lab, chấm cuối, khôi phục, chấm ngầm |  |  |  |  |  |  | ■ | ■ |  |  |  |
| M2: thư viện bài lab, tham số hệ thống |  |  |  |  |  |  |  | ■ | ■ |  |  |
| M2: bài lab mẫu, hoàn thiện |  |  |  |  |  |  |  |  |  | ■ | ■ |
| M3: tài liệu module | ■ | ■ |  |  |  |  |  |  |  |  |  |
| M3: lớp, nhập danh sách, lộ trình |  | ■ | ■ |  |  |  |  |  |  |  |  |
| M3: bài tập, thông báo |  |  |  | ■ | ■ | ■ |  |  |  |  |  |
| M3: kỳ thi, tham số hệ thống |  |  |  |  |  |  | ■ | ■ |  |  |  |
| M3: giám sát kỳ thi, gian lận, phúc khảo, bảng điểm |  |  |  |  |  |  |  | ■ | ■ |  |  |
| M3: hoàn thiện, kiểm thử |  |  |  |  |  |  |  |  |  | ■ | ■ |

| Mốc | Hạn | Tiêu chí |
|---|---|---|
| MA | Tuần 3 (15/11/2026) | Trên môi trường phát triển: học viên bắt đầu bài lab từ lộ trình, có terminal (`runtime-api` chạy phiên bằng Docker trên máy phát triển), kiểm tra bước và thấy điểm. Ba module gọi nhau bằng API thật |
| MB | Tuần 6 (06/12/2026) | Trên môi trường chuẩn qua script triển khai: cùng luồng như MA với gVisor, ghi phiên; hạn mức và đặt trước; dashboard phiên và sức chứa hoạt động; sao lưu chạy theo lịch; bộ kiểm thử cách ly đạt; phát hành `v0.1.0` |
| MC | Tuần 9 (27/12/2026) | Kỳ thi đầy đủ trên môi trường chuẩn: lên lịch có đặt trước, bắt đầu đồng loạt, mất máy phiên và bù giờ, nộp tự động, tín hiệu gian lận, phúc khảo |
| MD | Tuần 11 (10/01/2027) | Báo cáo đo toàn bộ Q1–Q14; biên bản thí nghiệm sự cố CH1–CH7; biên bản diễn tập khôi phục; biên bản cài đặt thử do người ngoài nhóm thực hiện; phát hành `v1.0.0` kèm hướng dẫn cài đặt và nâng cấp |

## 5. Phân công

| Thành viên | Phụ trách | Sản phẩm chính |
|---|---|---|
| Trưởng nhóm | M1; nền tảng; hợp đồng | Mã M1 (`services/runtime/`, `web/src/modules/runtime/`); hạ tầng (`infra/`, `deploy/`); giám sát; bộ kiểm thử cách ly, chương trình đo, thí nghiệm sự cố; tài liệu CR-DOC-01 đến CR-DOC-07 |
| Thành viên phụ trách M2 | M2 | Mã M2 (`services/content/`, `web/src/modules/content/`); bài lab mẫu (`labs/`, tối thiểu 10 bài thuộc đủ các dạng); tài liệu `docs/modules/m2-content.md` |
| Thành viên phụ trách M3 | M3 | Mã M3 (`services/learning/`, `web/src/modules/learning/`); tài liệu `docs/modules/m3-learning.md` |

Khung giao diện chung (`web/src/shell/`: điều hướng, đăng nhập, thông báo, hai ngôn ngữ) do trưởng nhóm dựng trước tuần 1.

Mỗi thành viên viết tài liệu thiết kế module của mình theo mẫu `docs/modules/_mau.md` trước khi bắt đầu viết mã
(hạn cuối tuần 1). Tài liệu module mô tả mô hình dữ liệu, máy trạng thái, API nội bộ cho giao diện, quy tắc nghiệp vụ
và cách kiểm thử các tiêu chí chấp nhận của module.

## 6. Tiêu chí hoàn thành

Một chức năng được coi là hoàn thành khi:

1. Mọi tiêu chí chấp nhận liên quan có kiểm thử tự động mang mã tiêu chí trong tên hoặc chú thích.
2. Kiểm thử hợp đồng của bên cung cấp đạt; sự kiện phát ra hợp lệ theo lược đồ.
3. Chức năng chạy trên môi trường chuẩn qua script triển khai, không chỉ trên máy phát triển.
4. Có metric và log đủ để chẩn đoán lỗi; chỉ số nghiệp vụ liên quan tại CR-DOC-04 mục 7.2 đã được phát.
5. Tài liệu module được cập nhật nếu thiết kế thay đổi.

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 26/10/2026 | Ban hành lần đầu |
| 1.1 | 26/10/2026 | Lộ trình theo ngày duyệt đề tài: thử nghiệm S1–S5 trải trong tuần 1–5, thứ tự việc theo issue và milestone |
| 1.2 | 27/10/2026 | Thử nghiệm S4 bổ sung lớp bọc bộ mô tả tệp của `setup` và đầu ra lệnh kiểm tra bị cắt |
