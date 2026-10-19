# Nhật ký sử dụng công cụ AI

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-DOC-08 |
| Phiên bản | 1.11 |
| Ngày ban hành | 18/10/2026 |

## 1. Mục đích

Tài liệu ghi nhận phạm vi sử dụng công cụ AI trong dự án Crucible, phân định phần việc do thành viên thực hiện và phần
việc do công cụ thực hiện, phục vụ minh bạch học thuật và báo cáo học phần. Nhật ký được bổ sung sau mỗi giai đoạn.

## 2. Công cụ và quy ước

| Mục | Nội dung |
|---|---|
| Công cụ | Claude Code (Anthropic), mô hình Claude Opus 5.5 |
| Hình thức | Trao đổi trực tiếp với công cụ; công cụ đọc, ghi tệp và chạy lệnh trên máy của thành viên khi được cho phép |
| Ghi nhận phần việc của công cụ | Ghi tại nhật ký này theo từng vòng làm việc |

| Ký hiệu | Ý nghĩa |
|---|---|
| **TV-QĐ** | Thành viên quyết định hoặc đặt yêu cầu |
| **TV-TH** | Thành viên trực tiếp thực hiện |
| **AI-ĐX** | Công cụ đề xuất, thành viên chọn hoặc chấp thuận |
| **AI-TH** | Công cụ trực tiếp thực hiện (soạn tài liệu, viết mã, chạy kiểm thử) |

### 2.1. Quy trình ghi nhật ký

Mỗi giai đoạn là một mục cấp 2 (`## n. Giai đoạn k — <tên>`), gồm các phần cố định theo thứ tự: thông tin chung, phần
việc của thành viên, phần việc của công cụ, sản phẩm, quá trình làm việc, các lần điều chỉnh hướng, ước lượng tỷ lệ
đóng góp.

| Thời điểm | Ghi vào | Nội dung |
|---|---|---|
| Thành viên đặt bài toán hoặc yêu cầu mới | Bảng "Quá trình làm việc": thêm một vòng | Bài toán, việc công cụ thực hiện, kết quả đánh giá của thành viên |
| Thành viên trả lời câu hỏi quyết định của công cụ | Bảng "Quyết định của thành viên" | Câu hỏi, các phương án, phương án được chọn |
| Đề xuất hoặc sản phẩm của công cụ không đạt và được điều chỉnh | Bảng "Các lần điều chỉnh hướng" | Đề xuất ban đầu, điều chỉnh của thành viên |
| Công cụ phát hiện lỗi qua kiểm thử | Bảng "Phần việc của công cụ" | Lỗi và cách khắc phục |
| Kết thúc giai đoạn | Các phần còn lại; mục "Lịch sử thay đổi" | Tổng hợp, ước lượng tỷ lệ; tăng phiên bản tài liệu |

Quy định:
- Mỗi vòng ghi kết quả đánh giá thực tế, kể cả khi kết quả không đạt.
- Nhật ký được cập nhật trong cùng pull request với sản phẩm của vòng đó.
- Mỗi thành viên sử dụng công cụ AI ghi phần của mình; cột "Người sử dụng công cụ" xác định người ghi.
- Quyết định nghiệp vụ và trải nghiệm người dùng do thành viên đưa ra; công cụ trình bày phương án kèm đề xuất và lý do.

## 3. Giai đoạn 1 — Chọn đề tài và thiết kế hệ thống

| Mục | Nội dung |
|---|---|
| Thời gian | 16/10/2026 – 17/10/2026 |
| Người sử dụng công cụ | Trưởng nhóm (`NhatMinh0311`) |
| Phạm vi | Chọn đề tài; đặc tả yêu cầu; kiến trúc; hợp đồng giữa module; vận hành; bảo mật; rủi ro và lộ trình; quy trình phát triển |

### 3.1. Phần việc của thành viên

| Hạng mục | Nội dung | Loại |
|---|---|---|
| Bối cảnh đề tài | Mục tiêu đề tài: nghiệp vụ phong phú, có nhu cầu thực tế, hệ thống đủ phức tạp để thực hành DevOps; chia việc theo module, thành viên làm fullstack | TV-QĐ |
| Hướng chọn đề tài | Tiếp cận theo thứ tự domain, nghiệp vụ, bài toán; nhóm quan tâm ngân hàng, học tập, y tế | TV-QĐ |
| Chọn đề tài | Chọn nền tảng thực hành có môi trường riêng và chấm tự động (domain học tập, sandbox) thay cho đề xuất ngân hàng của công cụ | TV-QĐ |
| Yêu cầu thiết kế | Thực hiện thiết kế hệ thống đầy đủ các bước theo khung của dự án thực tế | TV-QĐ |
| Hạ tầng | Dùng máy cá nhân 38 GB bộ nhớ, tạo máy ảo làm cụm; hỏi khả năng dùng máy của thành viên khác (macOS, Windows) | TV-QĐ |
| Công nghệ, tên | Go cho M1, Spring Boot cho M2 và M3; tên tiếng Anh cho dự án | TV-QĐ |
| Duyệt và chỉnh sửa | Loại bỏ phần so sánh với đề tài chấm bài lập trình; yêu cầu được tham gia các quyết định nghiệp vụ và trải nghiệm người dùng; hỏi mức đầy đủ của thiết kế so với khung thực tế | TV-QĐ |
| Mô hình cung cấp | Đặt câu hỏi giữa bán gói dịch vụ và mã nguồn mở tự cài; chọn phương án kết hợp | TV-QĐ |
| Sức chứa và lạm dụng tài nguyên | Nêu trường hợp chưa có: ước lượng tài nguyên của bài lab và báo giảng viên khi kỳ thi không đủ chỗ; giới hạn tài nguyên theo giảng viên; phiên bị dùng để chạy tiến trình nặng như đào tiền ảo | TV-QĐ |

### 3.2. Quyết định của thành viên

| Câu hỏi | Phương án được chọn | Phương án khác |
|---|---|---|
| Ngôn ngữ của các module | M1 Go (Kubernetes operator); M2, M3 Spring Boot | Toàn bộ Spring Boot; toàn bộ Go |
| Đối tượng người dùng chính | Học phần thực hành tại trường đại học | Người tự học; đào tạo nội bộ doanh nghiệp |
| Học viên có được kiểm tra bước trong kỳ thi | Giảng viên chọn cho từng kỳ thi (`FULL`, `RECORDED_ONLY`, `DISABLED`) | Không cho kiểm tra; cho kiểm tra như luyện tập |
| Trạng thái môi trường khi tiếp tục bài luyện tập | Khôi phục tới các bước đã đạt bằng lời giải mẫu | Môi trường trắng giữ điểm; để học viên chọn |
| Giao diện soạn bài lab | Form theo từng phần, đồng bộ hai chiều với YAML, có chạy thử | Trình soạn YAML; bài lab quản lý bằng Git |
| Cấu trúc kỳ thi | Một bài lab nhiều bước | Nhiều bài làm lần lượt; nhiều bài song song |
| Biến thể khi làm lại bài luyện tập | Biến thể mới mỗi lượt | Giữ nguyên biến thể |
| Chấm ngầm trong kỳ thi | Mỗi 2 phút | Mỗi 5 phút; không chấm ngầm |
| Xử lý tín hiệu gian lận | Biên bản với kết luận bác bỏ, nhắc nhở, hủy kết quả; học viên được thông báo và phúc khảo | Chỉ gắn cờ; cảnh báo và khóa phiên trực tiếp |
| Mô hình cung cấp | Mã nguồn mở tự cài trong đồ án; bản dịch vụ nhiều tổ chức là hướng phát triển | Chỉ tự cài; chỉ dịch vụ thu phí |
| Giấy phép | AGPL-3.0 | Apache-2.0 |
| Bù giờ khi sự cố hạ tầng | Tự động, tối đa 15 phút mỗi học viên | Không giới hạn; giảng viên quyết định |
| Vai trò trong thiết kế giao diện | Thành viên đóng góp ở mức trải nghiệm người dùng, không trực tiếp thiết kế hình ảnh | |
| Căn cứ nhu cầu | Hiện trạng tại trường và khảo sát dự án tương tự trên GitHub, không khảo sát người dùng | Phỏng vấn giảng viên, khảo sát sinh viên |
| Phân vai giữa thành viên và công cụ | Công cụ triển khai (viết mã, soạn tài liệu); thành viên cố vấn, kiểm định giải pháp, nêu trường hợp còn mù mờ | |
| Tài khoản và vào lớp | Đăng nhập bằng tài khoản trường; giảng viên nhập danh sách CSV | Tự đăng ký kèm mã lớp; cả hai |
| Vào thi muộn | Kết thúc chung, có hạn chót được vào thi | Thời lượng riêng; giảng viên chọn |
| Điểm bài tập nhiều lượt | Lượt cao nhất | Lượt cuối; giảng viên chọn |
| Trợ giảng | Có vai trò trợ giảng: giám sát, xem bản ghi, đề xuất biên bản | Nhiều giảng viên ngang quyền; một giảng viên |
| Bài tập quá hạn | Giảng viên chọn không nhận hoặc nhận muộn có trừ điểm; gia hạn riêng | Không nhận |
| Thư viện bài lab | Riêng tư hoặc chia sẻ; người khác chỉ nhân bản | Đồng tác giả; công khai toàn bộ |
| Quy trình AI viết mã | GitHub Spec Kit, `AGENTS.md`, skill `frontend-design` | Superpowers; chỉ `AGENTS.md` |
| Sinh viên xem gì sau công bố điểm thi | Điểm từng bước và bản ghi của mình | Không bản ghi; chỉ điểm tổng |
| Phúc khảo | Trong 7 ngày sau công bố, giảng viên phụ trách quyết định | Không giới hạn thời gian |
| Lộ trình | Không tính điểm học phần | Tính như một đầu điểm |
| Nộp bài thi sớm | Nộp là kết thúc | Được mở lại khi còn giờ |
| Bố cục trang làm bài | Chia đôi hướng dẫn và terminal, kéo được | Terminal toàn màn hình; hai tab |
| Kỳ thi không đủ sức chứa | Báo số chỗ thiếu và gợi ý khung giờ gần nhất còn đủ chỗ | Thêm đề xuất chia ca; chỉ báo lỗi |
| Tài nguyên của một phiên bài lab | Giảng viên chọn mức; kiểm thử khi xuất bản đo bộ nhớ, cảnh báo khi vượt 80%, chặn khi vượt giới hạn | Hệ thống tự chọn mức; không đo |
| Giới hạn tài nguyên theo giảng viên | Trần theo giảng viên do quản trị đặt (giờ máy, phiên đồng thời, chỗ × giờ đặt trước mỗi tháng); giảng viên chia cho các lớp | Quản trị đặt từng lớp; trần theo bộ môn |
| Bài tập cho lớp đông | Cảnh báo mềm khi tải ước tính gần hạn nộp vượt sức chứa | Không cảnh báo |
| Phiên luyện tập dùng CPU bất thường | Cảnh báo, tiếp diễn 5 phút thì kết thúc phiên | Chỉ ghi nhận; hạ mức CPU |
| Phiên thi dùng CPU bất thường | Không kết thúc, báo giám thị | Xử lý như luyện tập |
| Vi phạm nhiều lần | Chỉ ghi nhận, quản trị tự xử lý | Tự động tạm khóa tạo phiên luyện tập |
| Dạng bài bổ sung | Thêm điều tra và trả lời, bảo mật và gia cố vào MVP; viết script để sau | Chỉ hai dạng hiện có; thêm cả viết script |
| Bài mạng nhiều máy | Chuẩn bị hợp đồng (khái niệm máy), hiện thực sau MVP (ADR-0015) | Đưa vào MVP tối đa 2 máy; giữ nguyên |
| Trả lời sai ở bước trả lời | Tối đa N lần do giảng viên đặt, mặc định 3 | Trừ điểm mỗi lần sai; không giới hạn |
| So khớp câu trả lời | Chuẩn hóa khoảng trắng, hoa thường; nhiều đáp án; tùy chọn biểu thức chính quy | Khớp chính xác; thêm kiểu số có sai số |
| Mẫu kịch bản bổ sung | Sự cố tái phát, truy dấu xâm nhập, chuỗi manh mối, cứu dữ liệu, tự động hóa chấm theo hệ quả, chẩn đoán nguyên nhân | Ít quyền nhất (để sau) |
| Số lần chọn khi chẩn đoán trắc nghiệm | Mặc định 1, giảng viên đổi được | Mặc định 3 như tự luận |
| Chấm chẩn đoán nhiều lỗi | Điểm theo từng phương án, chọn sai bị trừ, không âm | Đúng cả tập mới đạt; mỗi lỗi một câu |
| Bài tường lửa cần `NET_ADMIN` | Cấp theo bài lab khai báo, chỉ khi chạy gVisor (ADR-0016) | Không cấp, bỏ bài tường lửa; cấp cho mọi phiên |
| Bài định tuyến | Mô phỏng nhiều mạng con trong một pod, phạm vi mở rộng, kiểm chứng bằng thử nghiệm S7 | Chỉ sửa tuyến trong một mạng con; router là pod riêng với mạng tầng 2 theo phiên |
| Thứ tự chẩn đoán và sửa | Khóa kiểm tra các bước sửa tới khi đã nộp chẩn đoán | Không ràng buộc; chẩn đoán trước khi có terminal |
| Kênh thông báo | Trong ứng dụng và email cho sự kiện quan trọng | Chỉ trong ứng dụng |
| Thiết bị | Làm bài trên máy tính; trang khác dùng được trên điện thoại | Chỉ máy tính |
| Trước giờ thi | Phòng chờ có kiểm tra kết nối | Phòng chờ đơn giản; không có phòng chờ |
| Lý do kết thúc khi môi trường dừng không do hết bộ nhớ | Lý do chung `ENVIRONMENT_FAILED`, giữ `OUT_OF_MEMORY` riêng | Dùng `RESOURCE_ABUSE`; dùng chung `OUT_OF_MEMORY` |
| Phiên thi bị dừng vì hết bộ nhớ hoặc lỗi môi trường | Mở môi trường mới trong thời gian còn lại, không bù giờ | Chỉ quy định không bù giờ; giảng viên quyết định |
| Cộng giờ thi thủ công | Cộng trong lúc thi cho từng sinh viên | Thêm thời gian ưu tiên khai báo trước; chỉ khai báo trước |
| Người được cộng giờ | Giảng viên phụ trách; trợ giảng đề xuất | Cả trợ giảng; thêm quản trị viên |
| Giới hạn cộng giờ | Tối đa 30 phút mỗi sinh viên mỗi kỳ thi, riêng với bù tự động; lý do bắt buộc, ghi nhật ký, thông báo sinh viên | Chung trần 15 phút; không giới hạn |
| Tham số nghiệp vụ cố định trong thiết kế | Quản trị viên chỉnh qua chức năng tham số hệ thống, không cố định trong mã | Giữ cố định; biến môi trường khi triển khai |
| Phạm vi tham số trên giao diện | Tham số nghiệp vụ; tham số kỹ thuật, bảo mật giữ ở cấu hình triển khai | Chỉ nhóm kỳ thi và hạn mức; mọi tham số |
| Thời điểm có hiệu lực | Kỳ thi chốt tham số khi lên lịch; tham số khác có hiệu lực ngay | Mọi thay đổi có hiệu lực ngay |
| Mức đặt tham số | Một giá trị cho toàn bản cài đặt | Giảng viên chỉnh trong giới hạn của quản trị viên |
| Vai trò quản trị cấp cao cho tham số hệ thống | Giữ một vai trò `ADMIN`; tách quản trị nền tảng và quản trị đơn vị ở bản dịch vụ nhiều đơn vị | Quyền riêng cho sửa tham số; thêm `SUPER_ADMIN` |
| Thêm sinh viên sau khi kỳ thi đã lên lịch | Tự tăng số chỗ đặt trước theo sĩ số; không tăng được thì vẫn thêm sinh viên và cảnh báo giảng viên kèm số chỗ thiếu | Không tự tăng; chặn thêm sinh viên khi không đủ chỗ |
| Phạm vi lần giao đầu của chế độ `docker` | Phiên, exec, terminal, sự kiện; đặt trước, hạn mức, tham số trả 501 | Thêm hạn mức và đặt trước |
| Xác thực ở chế độ `docker` | Token Keycloak phát triển | Tắt xác thực khi phát triển |
| Mạng của bài `network: mirror` ở chế độ `docker` | Ra Internet qua mạng `bridge` | Không có mạng |
| Ngôn ngữ giao diện | Tiếng Việt và tiếng Anh cho giao diện, trang đăng nhập, email, thông báo | Chỉ tiếng Việt |
| Ngôn ngữ nội dung bài lab | Một ngôn ngữ do giảng viên chọn, không dịch | Soạn song ngữ |
| Lưu ngôn ngữ người dùng | Thuộc tính `locale` của tài khoản Keycloak | Chỉ lưu trên trình duyệt |
| Ngôn ngữ email, thông báo | Theo ngôn ngữ người nhận | Song ngữ trong một email |

### 3.3. Phần việc của công cụ

| Hạng mục | Nội dung | Loại |
|---|---|---|
| Đề xuất | Tiêu chí chọn đề tài; đề tài cụ thể cho từng domain; đánh giá lại theo cách chia module; cách chia module theo nghiệp vụ cho nền tảng thực hành; đề bài có tham số và lỗi cài sẵn; phương án hạ tầng máy ảo; tên dự án | AI-ĐX |
| Thiết kế kỹ thuật | Kiến trúc ba module và nền tảng; operator và tài nguyên tùy biến; cách ly bằng gVisor; Keycloak; RabbitMQ với outbox; terminal qua gateway; nhóm môi trường ấm và đặt trước; GitOps; hạ tầng máy ảo; sinh biến thể bằng HMAC; lệnh kiểm tra bằng bộ công cụ chỉ đọc; quan sát; phát hành | AI-ĐX |
| Soạn thảo | CR-DOC-01 đến CR-DOC-12, ADR-0001 đến ADR-0013, tài liệu module M1 và mẫu tài liệu module, `README.md`, `CONTRIBUTING.md` | AI-TH |
| Viết mã | Lược đồ và bộ kiểm tra đặc tả bài lab, hai bài lab mẫu, `runtime-api.yaml`, `content-api.yaml`, `events.asyncapi.yaml`, workflow CI, CODEOWNERS, mẫu pull request | AI-TH |
| Kiểm chứng | Lint OpenAPI, kiểm tra AsyncAPI, kiểm tra đặc tả bài lab với tệp hợp lệ và tệp lỗi cố ý; chạy lệnh của CI trên máy; render toàn bộ sơ đồ Mermaid và xem lại hình; kiểm tra tham chiếu chéo giữa tài liệu | AI-TH |
| Khảo sát | Tra cứu dự án tương tự, giấy phép và mức hoạt động qua GitHub API; tra cứu framework và skill cho trợ lý AI | AI-TH |
| Phát hiện lỗi thiết kế qua kiểm tra | Cấu trúc binding AMQP sai trong AsyncAPI; lệnh kiểm tra của bài mẫu gọi chương trình của môi trường, trái nguyên tắc toàn vẹn chấm điểm; sức chứa 2 node phiên không đủ cho kịch bản thi 30 học viên và không còn chỗ cấp lại khi mất node, tăng lên 3 node; lời giải mẫu không chạy được với `PATH` của bộ công cụ kiểm tra, tách chế độ `RESTORE` | AI-TH |

### 3.4. Sản phẩm

| Sản phẩm | Soạn bởi | Yêu cầu và duyệt bởi |
|---|---|---|
| Tài liệu dự án (`docs/`, 12 tài liệu, 13 ADR, tài liệu module) | Công cụ | Thành viên |
| Quy ước cho trợ lý AI (`AGENTS.md`, `CLAUDE.md`) | Công cụ | Thành viên |
| Hợp đồng (`contracts/`) | Công cụ | Thành viên |
| Quy trình và CI (`.github/`, `CONTRIBUTING.md`) | Công cụ | Thành viên |

### 3.5. Quá trình làm việc

| Vòng | Thành viên đặt ra | Công cụ thực hiện | Kết quả đánh giá của thành viên |
|---|---|---|---|
| 1 | Chọn domain trong ngân hàng, học tập, y tế; cần hệ thống đủ phức tạp để học DevOps | Tiêu chí chọn, một đề tài cụ thể cho mỗi domain, khuyến nghị đề tài ngân hàng | Bổ sung ràng buộc: chia theo module, thành viên làm fullstack |
| 2 | Chia theo module | Đánh giá độ cân bằng module của từng đề tài | Chọn nền tảng thực hành dù công cụ khuyến nghị ngân hàng |
| 3 | Nền tảng thực hành, lo ngại cách chia module | Chia lại ba module theo nghiệp vụ; vai trò của trưởng nhóm ở M1 | Đạt; hỏi đề bài do giáo viên hay hệ thống tạo |
| 4 | Nguồn đề bài | Giáo viên soạn mẫu có tham số, hệ thống sinh biến thể; bài sửa hệ thống hỏng | Đạt; yêu cầu thiết kế hệ thống |
| 5 | Thiết kế hệ thống | Danh mục tài liệu thiết kế; câu hỏi về hạ tầng, công nghệ, tên | Hỏi khả năng chạy cụm trên ba máy cá nhân |
| 6 | Hạ tầng trên ba máy cá nhân | Phân tích hạn chế, đề xuất máy ảo trên một máy; kiểm tra cấu hình máy | Chọn máy ảo trên máy 38 GB; chọn công nghệ; yêu cầu tên tiếng Anh |
| 7 | Soạn bộ thiết kế | CR-DOC-01, CR-DOC-02 | Không đạt một phần: phần so sánh với đề tài chấm bài lập trình không cần thiết; hỏi khả năng tích hợp máy macOS, Windows |
| 8 | Tiếp tục | Hợp đồng, CR-DOC-03 đến CR-DOC-07, ADR, tài liệu M1, CI; mục máy thành viên tại CR-DOC-04 | Không đạt ở cách làm: công cụ tự chốt nhiều quyết định nghiệp vụ; yêu cầu được hỏi để tham gia; hỏi mức đầy đủ so với khung thực tế |
| 9 | Tham gia quyết định | Đánh giá phần còn thiếu so với khung thực tế; hỏi đợt 1 (đối tượng, phản hồi khi thi, khôi phục, soạn bài) | Chọn phương án cho 4 câu hỏi |
| 10 | Tiếp tục hỏi | Cập nhật tài liệu, hợp đồng theo đợt 1; hỏi đợt 2 (cấu trúc thi, biến thể, chấm ngầm, gian lận) | Chọn phương án cho 4 câu hỏi; hỏi mô hình cung cấp |
| 11 | Mô hình cung cấp | So sánh tự cài, dịch vụ, kết hợp; cập nhật phát hành và cài đặt; hỏi đợt 3 (giấy phép, bù giờ, giao diện) | Chọn kết hợp, AGPL-3.0, bù giờ tối đa 15 phút; đóng góp giao diện ở mức trải nghiệm người dùng |
| 12 | Công cụ AI thế hệ mới cho làm giao diện và viết mã có lập trình viên tham gia; bỏ khảo sát người dùng, thay bằng khảo sát dự án tương tự; phân vai; commit theo hạng mục | Tra cứu framework và skill; khảo sát dự án trên GitHub (CR-DOC-09); tách lịch sử thành commit theo hạng mục | Đạt |
| 13 | Làm các mục còn thiếu, hỏi khi luồng nghiệp vụ chưa rõ | Hỏi đợt 4–6 (tài khoản, thi muộn, điểm bài tập, trợ giảng, quá hạn, thư viện, quy trình AI, công bố, phúc khảo, lộ trình, nộp sớm); ràng buộc dự án; `AGENTS.md`; CR-DOC-10, CR-DOC-11, CR-DOC-12 | Chọn phương án cho các câu hỏi |
| 14 | Hỏi trải nghiệm người dùng của giao diện; hỏi khó khăn khi thành viên dùng Windows | Hỏi đợt 7 (bố cục trang làm bài, thông báo, thiết bị, phòng chờ); hướng dẫn phát triển trên Windows | Chọn phương án cho 4 câu hỏi |
| 15 | Hỏi hệ thống đã ước lượng được sức chứa, báo giảng viên khi kỳ thi không đủ chỗ, giới hạn tài nguyên theo giảng viên chưa | Đối chiếu thiết kế, chỉ ra 4 thiếu sót (phản hồi khi thiếu chỗ, tài nguyên bài lab không được đo, chưa có giới hạn theo giảng viên, bài tập không có cảnh báo tải); hỏi đợt 8 | Chọn phương án cho 4 câu hỏi |
| 16 | Hỏi cơ chế giám sát phiên bị dùng chạy tiến trình nặng | Đối chiếu biện pháp hiện có (chặn mạng, giới hạn tài nguyên, thời hạn); chỉ ra thiếu sót về CPU kéo dài và máy chủ gói có thể bị dùng làm proxy; hỏi đợt 9; cập nhật hợp đồng, tài liệu, backlog | Chọn phương án cho 3 câu hỏi |
| 17 | Hỏi các dạng bài hiện có; hỏi bài quản trị mạng và số pod mỗi phiên | Liệt kê hai dạng hiện có và ba dạng phổ biến chưa hỗ trợ trọn vẹn; phân tích ba mức bài mạng; hỏi đợt 10; thêm bước trả lời, hai bài mẫu, danh mục dạng bài, ADR-0015 | Chọn dạng điều tra và gia cố, để sau dạng viết script; chọn phương án cho 3 câu hỏi |
| 18 | Hỏi thêm dạng bài sáng tạo ít thay đổi | Đề xuất 5 mẫu kịch bản dùng cơ chế sẵn có và 2 dạng cần thay đổi nhỏ; hỏi đợt 11 về chẩn đoán; thêm trắc nghiệm, biến `{{faults}}`, `answerFirst`, bước chẩn đoán trong bài mẫu | Chọn 6 dạng; chọn phương án cho 3 câu hỏi |
| 19 | Hỏi khối lượng thay đổi khi cần nhiều máy cho một phiên; hỏi lại quy trình ra đề và chấm | Phân tích phần hạ tầng giữ nguyên và phần mã cần sửa; chỉ ra mâu thuẫn giữa bài tường lửa và quy định không thêm capability; mô tả đặc tả bài lab, kiểm thử khi xuất bản, chế độ phản hồi; thêm hướng phát triển vào README, ADR-0016 | Giữ một máy mỗi phiên cho MVP, đưa vào hướng phát triển; đồng ý cấp `NET_ADMIN` theo bài lab |
| 20 | Hỏi độ phức tạp của bài thay đổi bảng định tuyến; yêu cầu tra cứu kỹ tài liệu gVisor | Phân tích ba mức (trong một mạng con, nhiều mạng con trong một pod, router là máy riêng); đọc bảng tương thích syscall và mã nguồn netstack của gVisor để xác định cơ chế hỗ trợ; đề xuất user namespace con để không cấp `SYS_ADMIN`; thêm thử nghiệm S7 | Loại mức router là máy riêng; chọn mức nhiều mạng con trong một pod làm hướng phát triển |

### 3.6. Các lần điều chỉnh hướng

| Đề xuất ban đầu của công cụ | Điều chỉnh của thành viên |
|---|---|
| Đề tài chuyển tiền với sổ cái kép | Nền tảng thực hành có môi trường riêng và chấm tự động |
| Tài liệu có mục so sánh với hệ thống chấm bài lập trình | Bỏ phần so sánh; đề tài chấm bài lập trình chỉ là ví dụ về domain nhóm quan tâm |
| Công cụ tự chốt quy tắc nghiệp vụ khi soạn tài liệu | Quyết định nghiệp vụ và trải nghiệm người dùng được hỏi thành viên theo từng đợt (mục 3.2) |

### 3.7. Ước lượng tỷ lệ đóng góp

Ước lượng định tính dựa trên nguồn gốc của từng quyết định và sản phẩm trong lịch sử trao đổi.

| Hạng mục | Thành viên | Công cụ |
|---|---|---|
| Định hướng đề tài, đối tượng, mô hình cung cấp | 70% | 30% |
| Quyết định nghiệp vụ và trải nghiệm người dùng | 55% | 45% |
| Quyết định kiến trúc và công nghệ | 20% | 80% |
| Soạn tài liệu, hợp đồng, CI | 5% | 95% |
| Phát hiện thiếu sót, điều chỉnh hướng | 60% | 40% |
| **Tổng thể giai đoạn** | **khoảng 30%** | **khoảng 70%** |

## 4. Giai đoạn 2 — Thử nghiệm sớm

| Mục | Nội dung |
|---|---|
| Thời gian | Từ 18/10/2026 |
| Người sử dụng công cụ | Trưởng nhóm (`NhatMinh0311`) |
| Phạm vi | Các thử nghiệm S1–S7 tại CR-DOC-06 mục 3; cập nhật thiết kế theo kết quả |

Ở giai đoạn này, công cụ soạn script và phân tích kết quả; thành viên tự chạy lệnh trên máy ảo, đọc kết quả và thực hiện
các phép thử tương tác. Công cụ chỉ chạy trực tiếp lệnh đọc trạng thái và phép đo không ảnh hưởng tới máy chủ của
thành viên.

### 4.1. Phần việc của thành viên

| Hạng mục | Nội dung | Loại |
|---|---|---|
| Thứ tự thử nghiệm | Chọn làm S1 trước vì kết quả quyết định ADR-0003 | TV-QĐ |
| Thực hiện | Dựng máy ảo, cài k3s, gVisor; chạy từng script kiểm thử và gửi kết quả | TV-TH |
| Kiểm thử tương tác | Exec có TTY (nhập liệu, Ctrl-C, đổi kích thước), attach, tách và attach lại | TV-TH |

### 4.2. Phần việc của công cụ

| Hạng mục | Nội dung | Loại |
|---|---|---|
| Chuẩn bị | Cấu hình cloud-init, script tạo máy ảo bằng `virt-install`, script cài gVisor và cấu hình containerd cho k3s | AI-TH |
| Kiểm thử | Script kiểm thử cách ly I1–I8 đo trước và sau NetworkPolicy; thử hai cách gắn bộ công cụ chỉ đọc; đo số tiến trình trên node theo số tiến trình trong sandbox | AI-TH |
| Phân tích | Đọc nhật ký kubelet và nhân của node để xác định sandbox dừng do bộ điều khiển `pids`; suy ra yêu cầu `restartPolicy: Never` | AI-TH |
| Soạn thảo | `docs/spikes/S1.md`; cập nhật ADR-0003, CR-MOD-M1, CR-DOC-04, CR-DOC-05, CR-DOC-06 | AI-TH |
| Lỗi trong script của công cụ | Tham số `kubectl` bị shell phía xa tách lại khi gọi qua `ssh`; dùng `ps` không có trong image; ghi cấu hình kubelet vào thư mục chưa tồn tại làm lần thử giới hạn số tiến trình đầu tiên không có hiệu lực; lệnh tạo pod để thử attach đặt sai vị trí tham số. Các lỗi được phát hiện khi đọc kết quả và sửa trước khi kết luận | AI-TH |

### 4.3. Quá trình làm việc

| Vòng | Thành viên đặt ra | Công cụ thực hiện | Kết quả đánh giá của thành viên |
|---|---|---|---|
| 21 | Thực hiện S1 | Script dựng cụm thử, cài gVisor, kiểm thử I1–I8, gắn bộ công cụ, giới hạn số tiến trình, attach; phân tích kết quả từng bước; tổng hợp kết luận và cập nhật thiết kế | S1 đạt; xác nhận exec và attach tương tác bằng tay |
| 22 | Chọn lý do kết thúc và cách xử lý phiên thi khi môi trường bị dừng; hỏi cơ chế cộng giờ thủ công | Đề xuất phương án; rà thiết kế, chỉ ra thiếu cơ chế cộng giờ và thời hạn phiên thi không vượt hạn nộp ban đầu nên bù giờ không có hiệu lực; cập nhật hợp đồng, yêu cầu, ca sử dụng, giao diện, kiểm thử | Chọn `ENVIRONMENT_FAILED`; mở môi trường mới không bù giờ; giảng viên cộng giờ tối đa 30 phút, trợ giảng đề xuất |
| 23 | Yêu cầu tham số như trần cộng giờ do quản trị viên chỉnh, không cố định | Kiểm kê 18 tham số nghiệp vụ của ba module; tách tham số kỹ thuật, bảo mật; đề xuất mỗi module giữ tham số của mình với cùng khuôn API; kỳ thi chép tham số khi lên lịch; thay mức cộng cố định của thời hạn phiên thi bằng `sessionUntil` | Chọn phương án đề xuất cho 3 câu hỏi; hỏi có nên thêm vai trò quản trị cấp cao |
| 24 | Yêu cầu xóa dòng đồng tác giả của công cụ khỏi lịch sử commit; hỏi lại vai trò quản trị cấp cao | Viết lại lịch sử, kiểm tra nội dung từng commit không đổi, nêu ảnh hưởng của việc ghi đè nhánh trên GitHub trước khi đẩy; đổi quy ước commit; phân tích chi phí và lợi ích của vai trò quản trị cấp cao | Đồng ý ghi đè nhánh; giữ một vai trò `ADMIN` |
| 25 | Chọn cách xử lý sĩ số tăng sau khi lên lịch thi; yêu cầu hỏi bằng hộp chọn thay cho văn bản | Thêm thao tác đổi số chỗ đặt trước vào hợp đồng M1, cập nhật yêu cầu, ca sử dụng, giao diện, thiết kế M1 | Chọn tự tăng số chỗ, cảnh báo khi thiếu |
| 26 | Hỏi M2, M3 kiểm thử thế nào trên máy Windows 16 GB; yêu cầu làm chế độ `docker`; yêu cầu giao diện hai ngôn ngữ | Phân tích ba giai đoạn kiểm thử (mock, chế độ `docker`, cụm qua Tailscale); đặc tả, kế hoạch, danh sách việc `specs/001-runtime-docker-mode`; image bộ công cụ và image nền bài lab; cập nhật yêu cầu, hợp đồng, thiết kế giao diện cho hai ngôn ngữ | Duyệt đặc tả; chọn phương án cho 6 câu hỏi |
| 27 | Triển khai chế độ `docker` theo đặc tả đã duyệt | Viết `services/runtime` (API sinh từ hợp đồng, CSDL, outbox, xác thực, trình điều khiển Docker, vòng đời phiên, terminal WebSocket), kiểm thử tích hợp với PostgreSQL, RabbitMQ, Docker thật; job CI; profile `runtime` cho thành viên không cài Go; chạy Schemathesis, phát hiện và sửa 4 chỗ lệch giữa hợp đồng và dịch vụ | Chờ thành viên đọc lại mã |
| 28 | Hỏi chế độ `docker` đã hoàn thành bao nhiêu phần M1; hỏi cơ chế ghi lại lệnh và đầu ra của thí sinh; hỏi khả năng thay container bằng máy ảo; hỏi cách quản lý backlog | Đối chiếu mã với thiết kế M1: phần dùng lại, phần chỉ dùng cho môi trường phát triển, 14 thao tác chưa hiện thực; chỉ ra giới hạn của tái dựng dòng lệnh và đề xuất Runtime Monitoring của gVisor; phân tích điểm mở rộng của giao diện `Driver`; bổ sung Project board và đóng issue đã xong trong script tạo backlog, quy ước quản lý công việc tại CR-DOC-07 | Chọn GitHub Issues kèm Project board; tạm dừng viết mã M1 để tự học |
| 29 | Yêu cầu bật ruleset, xuất cấu hình; yêu cầu làm issue #3, #4 (khung giao diện, hai ngôn ngữ) | Xuất ruleset, phát hiện CODEOWNERS khiến pull request tài liệu của trưởng nhóm không thể được duyệt và giải thích; kiểm chứng Account API của Keycloak trước khi thiết kế; đặc tả, kế hoạch, danh sách việc `specs/002-web-shell`; hiện thực i18n có kiểm tra khóa khi biên dịch, kiểm thử chặn văn bản viết thẳng, menu tài khoản, chế độ màu, trang lỗi; kiểm thử thành phần và trên Chrome thật | Chọn sửa CODEOWNERS; chọn chế độ màu, cách hiển thị lỗi, cách tổ chức pull request; duyệt đặc tả |

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 17/10/2026 | Ghi nhận giai đoạn 1 |
| 1.1 | 17/10/2026 | Vòng 15, 16: sức chứa, hạn mức theo giảng viên, giám sát tài nguyên |
| 1.2 | 18/10/2026 | Vòng 17: dạng bài, bài mạng nhiều máy; vòng 18: mẫu kịch bản, chẩn đoán |
| 1.3 | 18/10/2026 | Vòng 19: hướng phát triển, capability `NET_ADMIN` |
| 1.4 | 18/10/2026 | Vòng 20: bài định tuyến, thử nghiệm S7 |
| 1.5 | 18/10/2026 | Giai đoạn 2, vòng 21: thử nghiệm S1 |
| 1.6 | 18/10/2026 | Vòng 22: môi trường bị dừng, cộng giờ thi |
| 1.7 | 18/10/2026 | Vòng 23: tham số hệ thống |
| 1.8 | 18/10/2026 | Vòng 24: lịch sử commit, vai trò quản trị |
| 1.9 | 18/10/2026 | Vòng 25: sĩ số lớp và đặt trước |
| 1.10 | 18/10/2026 | Vòng 26: chế độ docker, hai ngôn ngữ |
| 1.11 | 18/10/2026 | Vòng 27: chế độ docker của runtime-api |
