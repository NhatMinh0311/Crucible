# Nhật ký sử dụng công cụ AI

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-DOC-08 |
| Phiên bản | 1.0 |
| Ngày ban hành | 17/10/2026 |

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
| Kênh thông báo | Trong ứng dụng và email cho sự kiện quan trọng | Chỉ trong ứng dụng |
| Thiết bị | Làm bài trên máy tính; trang khác dùng được trên điện thoại | Chỉ máy tính |
| Trước giờ thi | Phòng chờ có kiểm tra kết nối | Phòng chờ đơn giản; không có phòng chờ |

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

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 17/10/2026 | Ghi nhận giai đoạn 1 |
