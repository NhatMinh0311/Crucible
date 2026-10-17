# Bài toán, phạm vi và yêu cầu

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-DOC-01 |
| Phiên bản | 1.2 |
| Ngày ban hành | 17/10/2026 |
| Tài liệu liên quan | CR-DOC-02 Kiến trúc hệ thống, CR-DOC-03 Hợp đồng giữa các module, CR-DOC-04 Vận hành, CR-DOC-05 Cách ly và bảo mật |

## 1. Bài toán

Crucible là nền tảng thực hành kỹ năng quản trị hệ thống Linux và dịch vụ mạng cơ bản. Mỗi học viên được cấp một môi trường
Linux riêng, thao tác qua terminal trên trình duyệt và được chấm điểm tự động dựa trên **trạng thái của hệ thống** sau khi
thao tác: tệp, tiến trình, dịch vụ, cấu hình, mạng.

Hệ thống phải giải quyết đồng thời ba bài toán:

1. **Cấp phát môi trường có trạng thái theo yêu cầu.** Mỗi phiên là một môi trường sống từ 30 đến 120 phút, cách ly
   với học viên khác và với hạ tầng của nền tảng. Môi trường phải sẵn sàng trong vài giây kể cả khi cả lớp bắt đầu
   thi cùng lúc, và phải được thu hồi đúng lúc để không chiếm tài nguyên.
2. **Chấm điểm đáng tin cậy trên môi trường do học viên kiểm soát.** Học viên có quyền quản trị trong môi trường của
   mình. Lệnh kiểm tra phải cho kết quả đúng dù học viên có thể sửa cấu hình, và phải đúng với mọi biến thể của đề.
3. **Tổ chức học và thi công bằng.** Mỗi học viên nhận một biến thể riêng của đề. Kỳ thi bắt đầu và kết thúc đồng loạt,
   có bù giờ khi sự cố hạ tầng làm gián đoạn phiên, có tín hiệu phát hiện gian lận để giảng viên xem xét.

### 1.1. Đối tượng và nhu cầu

Đối tượng chính là các học phần thực hành tại trường đại học: Hệ điều hành Linux, Mạng máy tính, Quản trị hệ thống.
Một lớp có từ 30 đến 60 sinh viên; giảng viên giao bài thực hành trong kỳ và tổ chức thi thực hành.

Tại trường chưa có nền tảng cho phép thực hành và thi quản trị hệ thống trên môi trường riêng có chấm tự động; hệ thống
hiện có phục vụ chấm bài lập trình. Khảo sát các giải pháp tương tự và khoảng trống Crucible hướng tới tại CR-DOC-09.
Các nhu cầu Crucible đáp ứng:

| Mã | Nhu cầu |
|---|---|
| N1 | Giảng viên mất nhiều thời gian chấm bài thực hành bằng cách xem trực tiếp hoặc đọc ảnh chụp màn hình |
| N2 | Bài thi thực hành chung một đề dễ bị chép kết quả giữa các sinh viên ngồi gần nhau |
| N3 | Sinh viên thiếu môi trường Linux để luyện tập, hoặc làm hỏng máy cá nhân khi thực hành quản trị hệ thống |
| N4 | Phòng máy không đủ hoặc không đồng nhất cấu hình cho thi thực hành đồng loạt |

### 1.2. Mô hình cung cấp

Crucible được phát hành dưới dạng mã nguồn mở; mỗi trường hoặc khoa tự cài đặt một bản trên hạ tầng của mình và một bản
cài đặt phục vụ một đơn vị đào tạo. Nhóm phát hành các phiên bản có đánh số, kèm Helm chart, mã dựng hạ tầng tham chiếu,
hướng dẫn cài đặt và nâng cấp. Cụm chuẩn của nhóm (CR-DOC-04) là một bản cài đặt mẫu.

Bản dịch vụ do nhóm vận hành cho nhiều đơn vị (nhiều tổ chức trên một bản cài đặt, gói theo số học viên hoặc giờ máy)
là hướng phát triển, không thực hiện trong đồ án. Thiết kế giữ khả năng mở rộng theo hướng này bằng các mã tham chiếu
mờ giữa module và hạn mức theo tài khoản hạn mức (ADR-0013).

## 2. Tác nhân

| Tác nhân | Vai trò | Mô tả |
|---|---|---|
| Học viên | `STUDENT` | Làm bài luyện tập và bài thi trong lớp đã ghi danh |
| Giảng viên | `INSTRUCTOR` | Soạn bài lab; với lớp được phân công là giảng viên phụ trách: quản lý lớp và lộ trình, tổ chức và giám sát kỳ thi, chấm phúc khảo, lập biên bản gian lận |
| Trợ giảng | Vai trò trong lớp | Người dùng được giảng viên phụ trách thêm vào lớp với vai trò trợ giảng: giám sát kỳ thi, xem bản ghi, đề xuất biên bản gian lận; không sửa điểm, không sửa cấu hình bài tập và kỳ thi |
| Quản trị nền tảng | `ADMIN` | Quản lý tài khoản, trần hạn mức tài nguyên của từng giảng viên, giám sát cụm |
| Người vận hành của đơn vị | Ngoài hệ thống | Cài đặt, nâng cấp, sao lưu, giám sát bản cài đặt trên hạ tầng của đơn vị theo tài liệu phát hành |
| Hệ thống chấm bài lab | Nội bộ | Chạy lời giải mẫu để kiểm thử bài lab trước khi xuất bản |

### 2.1. Ràng buộc dự án

| Loại | Ràng buộc |
|---|---|
| Nhân lực | 3 thành viên, mỗi người phụ trách một module; thành viên dùng trợ lý AI để viết mã |
| Thời gian | Khoảng 11 tuần từ khi đề tài được duyệt (CR-DOC-06 mục 4) |
| Hạ tầng | Một máy vật lý 16 luồng CPU, 38 GB bộ nhớ chạy toàn bộ cụm chuẩn; không có ngân sách cloud |
| Chi phí | Chỉ dùng phần mềm mã nguồn mở và dịch vụ miễn phí (GitHub, GHCR, Tailscale gói miễn phí) |
| Công nghệ | Go cho M1; Spring Boot cho M2, M3; React cho giao diện; Kubernetes (k3s) |
| Ngôn ngữ | Giao diện và tài liệu tiếng Việt |
| Giấy phép | Mã nguồn phát hành theo AGPL-3.0; phụ thuộc phải có giấy phép tương thích |
| Học phần | Đánh giá theo module của từng thành viên; tài liệu và nhật ký sử dụng AI phục vụ báo cáo |

## 3. Thuật ngữ

| Thuật ngữ | Định nghĩa |
|---|---|
| Bài lab | Đơn vị nội dung thực hành gồm mô tả môi trường, các bước, lệnh kiểm tra, gợi ý và điểm |
| Phiên bản bài lab | Ảnh chụp bất biến của bài lab tại thời điểm xuất bản; mọi lượt làm bài gắn với một phiên bản |
| Đặc tả bài lab | Tệp YAML mô tả một phiên bản bài lab theo lược đồ tại CR-DOC-03 |
| Tham số | Giá trị trong đặc tả được sinh riêng cho mỗi học viên, ví dụ tên người dùng, cổng dịch vụ |
| Biến thể | Bộ giá trị tham số cụ thể, sinh tất định từ một seed |
| Lỗi cài sẵn | Sự cố được chủ động đưa vào môi trường khi dựng, dùng cho bài "sửa hệ thống hỏng" |
| Phiên (session) | Một lần cấp môi trường; là khái niệm hạ tầng, thuộc M1 |
| Lượt làm bài (attempt) | Một lần học viên làm một phiên bản bài lab với một biến thể; là khái niệm nghiệp vụ, thuộc M2. Một lượt làm bài có thể dùng nhiều phiên liên tiếp |
| Lệnh kiểm tra | Lệnh chạy trong môi trường để xác định một bước đã đạt hay chưa |
| Nhóm môi trường ấm (warm pool) | Các môi trường dựng sẵn chưa gán cho ai, dùng để cấp phiên nhanh |
| Đặt trước tài nguyên | Phần sức chứa của cụm được giữ cho một kỳ thi trong khung giờ xác định |
| Hạn mức giờ máy | Tổng thời lượng phiên một tài khoản hạn mức được dùng trong một chu kỳ |
| Trần hạn mức của giảng viên | Giới hạn do quản trị đặt cho một giảng viên: tổng giờ máy mỗi tuần chia cho các lớp, số phiên đồng thời, số chỗ × giờ đặt trước cho kỳ thi mỗi tháng |

## 4. Phân chia module

Hệ thống được chia theo nghiệp vụ thành ba module. Mỗi module do một thành viên phụ trách toàn bộ từ giao diện,
dịch vụ, dữ liệu đến triển khai. Lý do và ranh giới được trình bày tại ADR-0001 và CR-DOC-02.

| Module | Tên | Trách nhiệm nghiệp vụ |
|---|---|---|
| M1 | Phiên thực hành và tài nguyên (`runtime`) | Vòng đời phiên, xếp hàng và ưu tiên, hạn mức giờ máy, nhóm môi trường ấm, đặt trước tài nguyên cho kỳ thi, terminal trên trình duyệt, ghi lại phiên, xử lý mất node |
| M2 | Bài lab và chấm điểm (`content`) | Soạn và xuất bản bài lab, quản lý phiên bản, sinh biến thể, kiểm thử bài lab, lượt làm bài, chấm theo bước, gợi ý và trừ điểm |
| M3 | Lớp học, lộ trình và thi (`learning`) | Lớp và ghi danh, lộ trình có điều kiện tiên quyết, bài tập có hạn nộp, kỳ thi, bù giờ, giám sát trực tiếp, xem lại phiên, tín hiệu gian lận, bảng điểm, phúc khảo |

## 5. Chức năng

Mã chức năng có tiền tố theo module. Chức năng đánh dấu (MR) thuộc phạm vi mở rộng sau khi hoàn thành MVP.

### 5.1. M1 — Phiên thực hành và tài nguyên

| Mã | Chức năng |
|---|---|
| R01 | Tạo phiên từ đặc tả môi trường đã kết xuất tham số; trả về mã phiên ngay, trạng thái cập nhật bất đồng bộ |
| R02 | Vòng đời phiên: chờ cấp, đang cấp, sẵn sàng, sắp hết hạn, kết thúc, lỗi, mất kết nối node |
| R03 | Gia hạn phiên trong giới hạn; kết thúc phiên khi hết thời lượng hoặc không hoạt động quá ngưỡng |
| R04 | Hàng đợi có ưu tiên `EXAM` > `PRACTICE` > `LAB_CI` khi cụm không đủ sức chứa; báo vị trí trong hàng đợi |
| R05 | Hạn mức theo tài khoản hạn mức: số phiên đồng thời và giờ máy mỗi tuần; sổ ghi nhận mức sử dụng. Tài khoản phân cấp: trần của giảng viên do quản trị đặt, hạn mức của lớp chia từ trần; ngân sách chỗ × giờ đặt trước cho kỳ thi mỗi tháng |
| R06 | Nhóm môi trường ấm theo cấu hình môi trường; tự bổ sung sau khi được cấp |
| R07 | Đặt trước tài nguyên cho kỳ thi; nhóm môi trường ấm được nâng trước giờ thi |
| R08 | Terminal trên trình duyệt qua WebSocket, truy cập bằng vé dùng một lần |
| R09 | Chế độ xem chỉ đọc terminal của học viên cho giảng viên |
| R10 | Ghi lại toàn bộ luồng terminal; phát dòng lệnh đã nhập cho phân tích gian lận (chỉ với phiên thi) |
| R11 | Chạy lệnh kiểm tra trong phiên theo yêu cầu của M2, tách biệt với môi trường thao tác của học viên |
| R12 | Phát hiện node mất kết nối; đánh dấu phiên bị ảnh hưởng và phát sự kiện |
| R13 | Trang quản trị: sức chứa cụm, hàng đợi, phiên đang chạy, hạn mức, lần phát hiện dùng CPU bất thường |
| R16 | Phát hiện phiên dùng CPU bất thường kéo dài: phiên luyện tập bị cảnh báo rồi kết thúc; phiên thi được báo cho giám thị, không bị kết thúc; mọi lần phát hiện được ghi nhận cho quản trị |
| R14 (MR) | Môi trường nhiều máy: 2–3 container nối mạng riêng trong một phiên; hợp đồng đã có khái niệm máy, hiện chỉ một máy `main` (ADR-0015) |
| R15 (MR) | Bài lab Docker: chạy Docker Engine bên trong môi trường học viên (phụ thuộc kết quả thử nghiệm S6 tại CR-DOC-06) |

### 5.2. M2 — Bài lab và chấm điểm

| Mã | Chức năng |
|---|---|
| C01 | Soạn bài lab bằng form theo từng phần (thông tin chung, tham số, môi trường, lỗi cài sẵn, bước, gợi ý, lời giải), đồng bộ hai chiều với tab YAML; kiểm tra lược đồ và quy tắc ngữ nghĩa khi nhập; xem trước mô tả các bước với một biến thể |
| C13 | Chạy thử bài lab nháp: giảng viên chọn seed, nhận một phiên có terminal, chạy từng lệnh kiểm tra, chạy lời giải mẫu |
| C02 | Tham số với bộ sinh `pick`, `range`, `pattern`, `subnet`; ràng buộc giá trị khác nhau |
| C03 | Lỗi cài sẵn chọn ngẫu nhiên theo seed từ danh sách do giảng viên khai báo |
| C04 | Xuất bản phiên bản bất biến; lượt làm bài đang chạy giữ phiên bản cũ |
| C05 | Kiểm thử bài lab khi xuất bản: lời giải mẫu đạt mọi bước với N seed; môi trường mới dựng không đạt bước nào; đo bộ nhớ cao nhất để cảnh báo khi mức tài nguyên đã chọn không đủ |
| C06 | Sinh biến thể tất định cho mỗi học viên theo ngữ cảnh (bài tập, kỳ thi) |
| C07 | Lượt làm bài: tạo, tiếp tục bằng phiên mới, nộp bài |
| C14 | Khôi phục tiến độ khi tiếp tục: chạy lời giải mẫu của các bước đã đạt trên phiên mới trước khi giao cho học viên |
| C15 | Chấm ngầm mỗi 2 phút các lượt làm bài thi đang diễn ra để ghi nhận tiến độ, phục vụ khôi phục và giám sát; kết quả không hiển thị cho học viên ngoài chế độ `FULL` |
| C08 | Kiểm tra từng bước theo yêu cầu học viên; điểm thành phần |
| C17 | Bước trả lời: học viên nhập câu trả lời cho câu hỏi về môi trường; so với đáp án theo biến thể, theo biểu thức chính quy hoặc theo trạng thái môi trường lúc trả lời; số lần trả lời tối đa do giảng viên đặt |
| C09 | Gợi ý theo bước, trừ điểm theo tỉ lệ khai báo; kỳ thi có thể tắt gợi ý |
| C16 | Chế độ phản hồi khi kiểm tra theo ngữ cảnh: `FULL` (hiện đạt hay chưa đạt), `RECORDED_ONLY` (chỉ báo đã ghi nhận), `DISABLED` (không cho kiểm tra) |
| C10 | Chấm cuối khi nộp bài: chạy lại toàn bộ lệnh kiểm tra trên trạng thái môi trường lúc nộp |
| C11 | Thư viện bài lab: tác giả đặt bài lab ở chế độ riêng tư hoặc chia sẻ trong bản cài đặt; giảng viên khác giao được bài chia sẻ hoặc nhân bản thành bài của mình để sửa; chỉ tác giả sửa được bài gốc |
| C12 (MR) | Bài lab có bước phụ thuộc: bước sau chỉ mở khi bước trước đạt |

### 5.3. M3 — Lớp học, lộ trình và thi

| Mã | Chức năng |
|---|---|
| L01 | Lớp học: tạo lớp; nhập danh sách sinh viên từ tệp CSV (mã số sinh viên, họ tên, email trường); sinh viên đăng nhập lần đầu bằng email trùng được gắn vào lớp; thêm giảng viên phụ trách và trợ giảng |
| L02 | Lộ trình: chuỗi bài lab luyện tập có điều kiện tiên quyết; bài sau mở khi bài trước đạt ngưỡng điểm; không giới hạn số lượt; không tính vào điểm học phần, giảng viên theo dõi tiến độ của lớp |
| L03 | Bài tập: giao bài lab cho lớp, có thời hạn và số lượt làm tối đa; điểm bài tập là điểm của lượt cao nhất; chính sách quá hạn do giảng viên chọn: không nhận (mặc định) hoặc nhận muộn trong hạn muộn tối đa với tỉ lệ trừ điểm mỗi ngày; gia hạn riêng cho từng sinh viên kèm lý do; khi giao hoặc đổi hạn nộp, hệ thống ước tính tải cao nhất gần hạn nộp và cảnh báo nếu vượt sức chứa (không chặn) |
| L04 | Kỳ thi: một bài lab (cố định phiên bản khi lên lịch), giờ bắt đầu, giờ kết thúc chung cho mọi sinh viên, hạn chót được vào thi, chế độ phản hồi khi kiểm tra, bật hoặc tắt gợi ý, bắt đầu đồng loạt, tự nộp khi hết giờ. Kỳ thi nhiều bài lab thuộc phạm vi mở rộng |
| L05 | Đặt trước tài nguyên tự động khi lên lịch kỳ thi |
| L06 | Bù giờ khi phiên thi bị gián đoạn do sự cố hạ tầng, có giới hạn và ghi nhật ký |
| L07 | Giám sát kỳ thi: trạng thái từng học viên, xem trực tiếp terminal, cảnh báo phiên dùng CPU bất thường |
| L08 | Xem lại bản ghi phiên: giảng viên phụ trách và trợ giảng xem mọi bản ghi của lớp; sinh viên xem bản ghi phiên thi của chính mình sau khi điểm được công bố |
| L09 | Tín hiệu gian lận: lộ giá trị biến thể, chuỗi lệnh tương đồng, kết nối từ hai địa chỉ, dán khối lớn |
| L13 | Xử lý tín hiệu gian lận: giảng viên lập biên bản với một trong ba kết luận `DISMISSED` (bác bỏ), `WARNED` (nhắc nhở), `VOIDED` (hủy kết quả, điểm 0) kèm lý do; học viên được thông báo và được phúc khảo |
| L10 | Bảng điểm lớp; xuất CSV |
| L11 | Phúc khảo: trong thời hạn (mặc định 7 ngày sau công bố, giảng viên đổi được), sinh viên gửi một yêu cầu cho mỗi bài, nêu bước và lý do; trợ giảng góp ý; giảng viên phụ trách giữ hoặc sửa điểm kèm trả lời |
| L14 | Công bố điểm kỳ thi: sinh viên xem điểm từng bước và bản ghi phiên của mình; không xem lời giải mẫu |
| L15 | Thông báo trong ứng dụng cho mọi người dùng; email cho sự kiện quan trọng (CR-DOC-11 mục 8.2), gửi qua outbox |
| L16 | Phòng chờ thi mở trước 15 phút: quy chế, đếm ngược, kiểm tra kết nối bằng một môi trường mẫu ngắn |
| L12 | Hạn mức giờ máy theo lớp, giảng viên chia từ trần hạn mức của mình |

### 5.4. Dạng bài lab

Các dạng bài được ghép từ cơ chế chung của đặc tả bài lab (tham số, lỗi cài sẵn, bước kiểm tra, bước trả lời). Một bài
lab có thể kết hợp nhiều dạng.

| Dạng | Học viên làm gì | Cơ chế | Phạm vi |
|---|---|---|---|
| Cấu hình theo yêu cầu | Dựng, cấu hình dịch vụ trên môi trường sạch theo đề | Tham số, bước kiểm tra | MVP |
| Sửa hệ thống hỏng | Tìm nguyên nhân và khắc phục sự cố được cài sẵn | Lỗi cài sẵn chọn theo seed, bước kiểm tra | MVP |
| Bảo mật và gia cố | Sửa cấu hình, quyền, dịch vụ không an toàn | Lỗi cài sẵn hoặc setup, bước kiểm tra | MVP |
| Điều tra và trả lời | Đọc nhật ký, trạng thái hệ thống để trả lời câu hỏi | Tham số, bước trả lời (C17) | MVP |
| Mạng trong một máy | Cổng, địa chỉ lắng nghe, phân giải tên, tường lửa cục bộ | Như các dạng trên; tường lửa phụ thuộc kết quả thử nghiệm S4 | MVP |
| Mạng nhiều máy | Dịch vụ giữa các máy: DNS, cơ sở dữ liệu từ xa, reverse proxy, tường lửa giữa máy | Môi trường nhiều máy (R14) | Mở rộng |
| Viết script | Viết script đạt yêu cầu, hệ thống chạy script để chấm | Cần chế độ chấm chạy mã của học viên trong điều kiện kiểm soát, khác nguyên tắc tại ADR-0011 | Mở rộng |

### 5.5. Chức năng dùng chung

| Mã | Chức năng | Phụ trách |
|---|---|---|
| P01 | Đăng nhập qua Keycloak; liên kết nhà cung cấp định danh của trường (ví dụ Microsoft Entra ID của tài khoản email trường) khi đơn vị cấu hình, tài khoản cục bộ khi không có; vai trò `STUDENT`, `INSTRUCTOR`, `ADMIN` | Nền tảng (trưởng nhóm) |
| P02 | Khung giao diện chung: điều hướng, đăng nhập, thông báo | Nền tảng |
| P03 | Phát hành phiên bản: image, Helm chart tổng, ghi chú phát hành; hướng dẫn cài đặt và nâng cấp cho người vận hành | Nền tảng |

## 6. Ngoài phạm vi

- Nhiều tổ chức dùng chung một bản cài đặt (multi-tenant), gói dịch vụ và thanh toán. Đây là hướng phát triển của bản dịch vụ (mục 1.2).
- Môi trường dạng máy ảo, môi trường Windows, môi trường có giao diện đồ họa.
- Sinh đề hoặc chấm điểm bằng mô hình AI.
- Người tự học không thuộc lớp nào; khóa học công khai.
- Chấm bài lập trình theo đầu vào và đầu ra.
- Bài mạng tầng liên kết dữ liệu và định tuyến (VLAN, DHCP, định tuyến giữa các mạng con): mạng của cụm chỉ kết nối các pod ở tầng IP (ADR-0015).
- Ứng dụng di động, đa ngôn ngữ.

## 7. User story và tiêu chí chấp nhận

Mã tiêu chí có dạng `AC-<module><số>.<số>`, dùng để truy vết từ kiểm thử.

### 7.1. M1

**US-R1 (R01–R03, R06)** — Là học viên, tôi muốn có môi trường làm bài ngay khi bấm bắt đầu.
- AC-R1.1: Khi nhóm môi trường ấm còn chỗ, thời gian từ lúc tạo phiên đến lúc terminal nhận được lệnh đầu tiên có p95 ≤ 5 giây.
- AC-R1.2: Khi không có môi trường ấm, thời gian trên có p95 ≤ 30 giây.
- AC-R1.3: Mỗi học viên có tối đa một phiên luyện tập và một phiên thi đang hoạt động; yêu cầu vượt bị từ chối (409 `SESSION_LIMIT`).
- AC-R1.4: Phiên không có thao tác trong thời gian `idleMinutes` bị kết thúc với lý do `IDLE`; học viên nhận cảnh báo trước 2 phút.
- AC-R1.5: Gia hạn không vượt quá thời lượng tối đa của cấu hình môi trường (409 `EXTENSION_LIMIT`); phiên thi không gia hạn được qua M1 mà theo quy tắc bù giờ của M3.

**US-R2 (R04, R05, R07)** — Là quản trị nền tảng, tôi muốn tài nguyên được phân phối công bằng và kỳ thi không thiếu chỗ.
- AC-R2.1: Khi cụm hết sức chứa, yêu cầu mới vào hàng đợi theo ưu tiên rồi theo thời điểm yêu cầu; phiên đang chạy không bị dừng để nhường chỗ.
- AC-R2.2: Trong khung giờ đặt trước, phần sức chứa đặt trước chỉ cấp cho phiên `EXAM` của kỳ thi tương ứng.
- AC-R2.3: Yêu cầu chờ quá 10 phút bị hủy với lý do `QUEUE_TIMEOUT`.
- AC-R2.4: Tài khoản hạn mức hết giờ máy trong tuần không tạo được phiên `PRACTICE` (402 `QUOTA_EXHAUSTED`); phiên `EXAM` không bị chặn bởi giờ máy.
- AC-R2.5: Mức sử dụng được ghi theo phút thực tế của phiên, kể cả phiên kết thúc bất thường.
- AC-R2.6: Tổng giờ máy mỗi tuần chia cho các lớp không vượt trần của giảng viên (422 `QUOTA_CEILING_EXCEEDED`); tổng phiên đang chạy của các lớp không vượt số phiên đồng thời của trần (402 `QUOTA_EXHAUSTED`).
- AC-R2.7: Đặt trước làm tổng chỗ × giờ trong tháng của giảng viên vượt ngân sách bị từ chối (402 `RESERVATION_BUDGET_EXHAUSTED`); hủy đặt trước trước giờ bắt đầu hoàn lại ngân sách.
- AC-R2.8: Chỉ quản trị đặt được trần hạn mức của giảng viên; dịch vụ M3 chỉ ghi được hạn mức của lớp.

**US-R3 (R08–R10)** — Là học viên, tôi muốn thao tác terminal ổn định; là giảng viên, tôi muốn xem được học viên đang làm gì trong kỳ thi.
- AC-R3.1: Vé terminal chỉ cấp cho chủ phiên, dùng một lần, hết hạn sau 30 giây nếu không được dùng.
- AC-R3.2: Mất kết nối WebSocket rồi kết nối lại bằng vé mới không làm mất trạng thái môi trường.
- AC-R3.3: Kết nối xem trực tiếp không gửi được dữ liệu vào terminal.
- AC-R3.4: Bản ghi phiên phát lại được đúng trình tự và nhịp thời gian của luồng terminal.

**US-R4 (R12)** — Là học viên, tôi muốn được cấp lại môi trường khi hạ tầng gặp sự cố.
- AC-R4.1: Phiên trên node mất kết nối được chuyển sang `LOST` trong ≤ 60 giây kể từ khi node mất kết nối.
- AC-R4.2: Không còn pod hoặc tài nguyên nào thuộc phiên đã kết thúc sau 2 phút (không có môi trường mồ côi).

**US-R5 (R16)** — Là học viên, tôi muốn terminal không bị chậm vì người khác chạy tiến trình nặng; là quản trị, tôi muốn biết ai lạm dụng tài nguyên.
- AC-R5.1: Phiên luyện tập dùng từ 90% giới hạn CPU trở lên, tính trung bình trong 10 phút, nhận cảnh báo trong terminal; nếu mức dùng trung bình trong 5 phút tiếp theo vẫn từ 90% trở lên, phiên bị kết thúc với lý do `RESOURCE_ABUSE`, lượt làm bài tạm dừng và giữ điểm các bước đã đạt.
- AC-R5.2: Phiên thi thỏa điều kiện trên không bị kết thúc; giám thị thấy cảnh báo trên màn hình giám sát.
- AC-R5.3: Mỗi lần phát hiện được ghi nhận với người dùng, phiên, thời điểm, mức dùng và cách xử lý; quản trị xem được theo người dùng. Hệ thống không tự động hạn chế người dùng vi phạm nhiều lần.
- AC-R5.4: Phiên của bài lab khai báo `cpuIntensive` và phiên `LAB_CI` không bị giám sát CPU bất thường.

### 7.2. M2

**US-C1 (C01–C05, C11, C13)** — Là giảng viên, tôi muốn soạn bài lab và chắc chắn bài lab giải được trước khi giao.
- AC-C1.1: Đặc tả không hợp lệ theo lược đồ bị từ chối kèm vị trí lỗi; trên form, lỗi hiển thị tại trường tương ứng.
- AC-C1.2: Phiên bản chỉ được xuất bản khi lời giải mẫu đạt toàn bộ bước với 5 seed khác nhau và môi trường mới dựng không đạt bước nào (trừ bước khai báo `expectInitially: pass`).
- AC-C1.3: Phiên bản đã xuất bản không sửa được; sửa bài lab tạo phiên bản nháp mới.
- AC-C1.4: Lượt làm bài đang dùng phiên bản cũ tiếp tục với phiên bản cũ sau khi phiên bản mới được xuất bản.
- AC-C1.5: Sửa trên form rồi chuyển sang tab YAML và ngược lại không làm mất hoặc thay đổi nội dung ngoài phần đã sửa.
- AC-C1.6: Phiên chạy thử của giảng viên dùng biến thể theo seed được chọn và không tạo lượt làm bài.
- AC-C1.7: Bài lab riêng tư không hiển thị với giảng viên khác; nhân bản tạo bài lab mới thuộc người nhân bản, ghi nguồn gốc.
- AC-C1.8: Kiểm thử bài lab không đạt khi một seed bị dừng vì vượt giới hạn bộ nhớ; kết quả kiểm thử cảnh báo khi bộ nhớ cao nhất vượt 80% giới hạn của mức tài nguyên đã chọn, kèm đề xuất mức tài nguyên.

**US-C2 (C06–C10, C14–C17)** — Là học viên, tôi muốn biết bước nào đã làm đúng và nhận điểm công bằng.
- AC-C2.1: Một lượt làm bài luôn dùng cùng một biến thể, kể cả khi tiếp tục bằng phiên mới; lượt làm bài mới của cùng bài luyện tập nhận biến thể mới; học viên khác nhau trong cùng kỳ thi nhận biến thể khác nhau khi miền tham số cho phép.
- AC-C2.2: Học viên không suy ra được seed hoặc biến thể của người khác từ thông tin của mình.
- AC-C2.3: Kết quả kiểm tra một bước trả về trong p95 ≤ 3 giây.
- AC-C2.4: Thay thế hoặc sửa các công cụ hệ thống trong môi trường (`curl`, `id`, `grep`...) không làm thay đổi kết quả lệnh kiểm tra.
- AC-C2.5: Với luyện tập, bước đã đạt được giữ điểm kể cả khi phiên kết thúc. Với kỳ thi, điểm lấy theo lần chấm cuối trên trạng thái môi trường lúc nộp.
- AC-C2.6: Xem gợi ý trừ điểm của bước theo tỉ lệ khai báo; gợi ý bị tắt trong kỳ thi có cấu hình tắt gợi ý.
- AC-C2.7: Lượt làm bài có phiên bị `LOST` hoặc đã hết hạn tiếp tục được bằng phiên mới; các bước đã đạt được giữ điểm.
- AC-C2.8: Khi tiếp tục, môi trường mới được khôi phục bằng lời giải mẫu của các bước đã đạt trước khi học viên nhận terminal; học viên không thấy nội dung lời giải trong terminal. Bước khôi phục không thành công được báo cho học viên và không bị trừ điểm.
- AC-C2.9: Trong ngữ cảnh `RECORDED_ONLY`, học viên chỉ nhận thông báo đã ghi nhận; trong ngữ cảnh `DISABLED`, nút kiểm tra không hiển thị và API từ chối (409 `CHECK_DISABLED`).
- AC-C2.10: Ở chế độ `FULL`, bước trả lời nhận tối đa `maxTries` lần; hết lượt mà chưa đúng thì bước không đạt và API từ chối (409 `ANSWER_TRIES_EXHAUSTED`). Câu trả lời khác đáp án chỉ ở khoảng trắng đầu cuối hoặc hoa thường (khi không bật phân biệt) được chấp nhận.
- AC-C2.11: Ở chế độ `RECORDED_ONLY`, `DISABLED`, học viên chỉ nhận thông báo đã ghi nhận câu trả lời; khi chấm cuối, bước trả lời được chấm theo câu trả lời cuối cùng.

### 7.3. M3

**US-L1 (L01–L03)** — Là giảng viên, tôi muốn tổ chức lớp học theo lộ trình.
- AC-L1.1: Học viên chỉ bắt đầu được bài lab khi thuộc lớp và đã đạt ngưỡng điểm của các bài tiên quyết (403 `PREREQUISITE_NOT_MET`).
- AC-L1.2: Bài tập quá hạn (đã tính gia hạn riêng) không tạo được lượt làm bài mới (409 `ASSIGNMENT_CLOSED`), trừ khi cho nhận muộn và còn trong hạn muộn tối đa; lượt đang làm được nộp tự động tại thời điểm hết hạn cuối cùng.
- AC-L1.3: Số lượt làm bài không vượt quá số lượt tối đa của bài tập, kể cả khi có yêu cầu đồng thời.
- AC-L1.4: Điểm bài tập bằng điểm cao nhất trong các lượt đã chấm, sau khi trừ điểm nộp muộn của từng lượt.
- AC-L1.5: Nhập danh sách CSV báo lỗi theo từng dòng (thiếu cột, email sai định dạng, trùng mã số); sinh viên có email trong danh sách đăng nhập lần đầu được gắn vào lớp, không cần thao tác thêm.
- AC-L1.6: Khi giao bài tập hoặc đổi hạn nộp, nếu số phiên ước tính chạy cùng lúc trong 3 giờ trước hạn nộp vượt số chỗ còn trống của khung giờ đó hoặc vượt số phiên đồng thời của lớp, giảng viên nhận cảnh báo kèm số liệu; bài tập vẫn được lưu nếu giảng viên xác nhận.

**US-L2 (L04–L07)** — Là giảng viên, tôi muốn tổ chức kỳ thi đồng loạt cho cả lớp.
- AC-L2.1: Lên lịch kỳ thi tạo yêu cầu đặt trước tài nguyên tương ứng số học viên; cụm không đủ sức chứa cho khung giờ đó thì kỳ thi không được lên lịch (409 `CAPACITY_UNAVAILABLE`). Giảng viên được báo số chỗ cần, số chỗ còn và tối đa 3 khung giờ gần nhất cùng độ dài còn đủ chỗ trong khoảng 07:00–21:00.
- AC-L2.2: 30 học viên bắt đầu trong cùng một phút đều được cấp phiên, p95 thời gian cấp ≤ 10 giây.
- AC-L2.3: Khi hết giờ, mọi lượt làm bài của kỳ thi được nộp tự động và chấm cuối; phiên được thu hồi.
- AC-L2.4: Phiên thi bị `LOST` được bù đúng thời gian gián đoạn, tối đa 15 phút cho mỗi học viên mỗi kỳ thi; mỗi lần bù được ghi nhật ký.
- AC-L2.5: Chỉ giảng viên phụ trách và trợ giảng của lớp xem được trực tiếp terminal và bản ghi của học viên trong kỳ thi.
- AC-L2.7: Sinh viên bắt đầu sau hạn chót được vào thi bị từ chối (409 `EXAM_ENTRY_CLOSED`); mọi sinh viên được nộp tự động tại giờ kết thúc chung cộng thời gian bù giờ của riêng mình.
- AC-L2.8: Nộp bài thi trước giờ kết thúc cần xác nhận hai lần; sau khi nộp, lượt làm bài được chấm cuối, phiên được thu hồi và sinh viên không vào lại được.
- AC-L2.9: Trước khi công bố, sinh viên không xem được điểm và bản ghi của kỳ thi; sau khi công bố, sinh viên xem được điểm từng bước và bản ghi của chính mình.
- AC-L2.6: Màn hình giám sát hiển thị tiến độ từng học viên theo kết quả chấm ngầm gần nhất, kể cả khi học viên không được xem kết quả kiểm tra.
- AC-L2.10: Sự kiện `runtime.session.resource-alert` của phiên thi hiển thị thành cảnh báo trên dòng của học viên tại màn hình giám sát trong ≤ 30 giây; cảnh báo không phải tín hiệu gian lận và không ảnh hưởng điểm.

**US-L3 (L08–L11, L13, L14)** — Là giảng viên, tôi muốn phát hiện dấu hiệu gian lận và xử lý khiếu nại.
- AC-L3.1: Học viên nhập giá trị tham số thuộc biến thể của học viên khác trong cùng kỳ thi tạo tín hiệu `VARIANT_LEAK`.
- AC-L3.2: Tín hiệu gian lận không tự động thay đổi điểm; trợ giảng đề xuất, giảng viên phụ trách lập biên bản có kết luận và lý do.
- AC-L3.4: Kết luận `VOIDED` đặt điểm lượt làm bài về 0 qua điều chỉnh điểm của M2 (lưu điểm cũ); học viên nhận thông báo kèm lý do và gửi được yêu cầu phúc khảo.
- AC-L3.3: Điều chỉnh điểm khi phúc khảo lưu điểm cũ, điểm mới, người điều chỉnh và lý do.
- AC-L3.5: Yêu cầu phúc khảo gửi sau thời hạn hoặc gửi lần thứ hai cho cùng một bài bị từ chối (409 `REGRADE_NOT_ALLOWED`).

## 8. Yêu cầu chất lượng

Yêu cầu chất lượng được viết dưới dạng kịch bản có số đo. Mỗi kịch bản là căn cứ cho một quyết định kiến trúc
tại CR-DOC-02 và được kiểm chứng bằng phương pháp ghi ở cột cuối. Môi trường đo là cụm chuẩn mô tả tại CR-DOC-04
(1 node nền tảng, 3 node chạy phiên).

| Mã | Thuộc tính | Kịch bản | Số đo | Kiểm chứng |
|---|---|---|---|---|
| Q1 | Hiệu năng | Học viên bắt đầu bài luyện tập khi nhóm ấm còn chỗ | Thời gian sẵn sàng p95 ≤ 5 s | k6, chỉ số `crucible_session_ready_seconds` |
| Q2 | Khả năng chịu tải đột biến | 30 học viên bắt đầu kỳ thi trong 60 giây | 100% được cấp; p95 ≤ 10 s | k6 kịch bản thi |
| Q3 | Sức chứa | Tăng dần số phiên đồng thời có học viên giả lập thao tác | Xác định số phiên tối đa giữ được Q1 và độ trễ terminal p95 ≤ 150 ms | k6, báo cáo sức chứa |
| Q4 | Tính sẵn sàng | Tắt một node chạy phiên | Phiên bị ảnh hưởng chuyển `LOST` ≤ 60 s; học viên được cấp lại ≤ 30 s; không mất điểm đã đạt | Thí nghiệm chaos, biên bản |
| Q5 | Tính đúng đắn tài nguyên | Khởi động lại bộ điều khiển phiên, xóa pod, mất kết nối broker trong lúc tải | Không có môi trường mồ côi sau 2 phút; không có phiên tính giờ sai | Thí nghiệm chaos, chỉ số `crucible_orphan_pods` |
| Q6 | Cách ly | Từ trong môi trường học viên thử truy cập Kubernetes API, metadata, pod khác, dịch vụ nền tảng, thoát sandbox | Mọi thử nghiệm thất bại | Bộ kiểm thử cách ly chạy trên cụm (CR-DOC-05) |
| Q7 | Tính toàn vẹn chấm điểm | Học viên thay công cụ hệ thống trong môi trường | Kết quả kiểm tra không đổi | Kiểm thử tích hợp M2 |
| Q8 | Không mất sự kiện | Dừng broker 2 phút trong lúc có phiên kết thúc và bài được chấm | Mọi sự kiện được giao sau khi broker phục hồi; dữ liệu các module nhất quán | Thí nghiệm chaos |
| Q9 | Khả năng triển khai | Cập nhật phiên bản một module trong lúc có phiên đang chạy | Không phiên nào bị kết thúc; terminal kết nối lại tự động | Thí nghiệm triển khai, GitOps |
| Q10 | Khả năng quan sát | Một yêu cầu tạo phiên bị chậm | Truy vết được từ giao diện qua M3, M2, M1 tới pod bằng một trace ID | Tempo, Grafana |
| Q11 | Khôi phục dữ liệu | Mất cơ sở dữ liệu | RPO ≤ 15 phút, RTO ≤ 30 phút | Biên bản diễn tập khôi phục |
| Q12 | Khả năng tái lập | Dựng lại toàn bộ cụm từ mã nguồn | Một lệnh dựng hạ tầng, một lệnh khởi tạo GitOps; ≤ 30 phút | Biên bản dựng lại |
| Q13 | Khả năng cài đặt | Người chưa tham gia dự án cài đặt bản phát hành lên một cụm Kubernetes có node hỗ trợ gVisor, chỉ theo hướng dẫn cài đặt | Hoàn tất trong ≤ 1 giờ, không cần hỏi nhóm | Biên bản cài đặt thử do người ngoài nhóm thực hiện |
| Q14 | Khả năng nâng cấp | Nâng cấp từ bản phát hành trước lên bản mới khi đã có dữ liệu | Không mất dữ liệu; phiên đang chạy không bị kết thúc | Kiểm thử nâng cấp tự động trong CI |

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 17/10/2026 | Ban hành lần đầu |
| 1.1 | 17/10/2026 | Phản hồi khi thiếu sức chứa (AC-L2.1), đo bộ nhớ bài lab (C05, AC-C1.8), trần hạn mức theo giảng viên (R05, AC-R2.6–AC-R2.8), cảnh báo tải bài tập (AC-L1.6), giám sát CPU (R16, US-R5, AC-L2.10) |
| 1.2 | 18/10/2026 | Dạng bài lab (mục 5.4), bước trả lời (C17, AC-C2.10, AC-C2.11), mạng tầng 2 ngoài phạm vi |
