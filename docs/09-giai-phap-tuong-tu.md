# Khảo sát giải pháp tương tự

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-DOC-09 |
| Phiên bản | 1.0 |
| Ngày ban hành | 26/10/2026 |
| Tài liệu liên quan | CR-DOC-01 mục 1.1 (đối tượng và nhu cầu) |

## 1. Mục đích và phương pháp

Tài liệu đánh giá các nền tảng hiện có cho thực hành quản trị hệ thống có môi trường riêng và chấm tự động, nhằm xác
định khoảng trống mà Crucible hướng tới. Thông tin về giấy phép, mức độ hoạt động lấy từ GitHub API và tài liệu chính
thức của từng dự án tại thời điểm 17/10/2026.

Tiêu chí so sánh xuất phát từ bài toán tại CR-DOC-01 mục 1:

| Mã | Tiêu chí |
|---|---|
| T1 | Mã nguồn mở, tự cài đặt được |
| T2 | Môi trường Linux riêng cho mỗi học viên, terminal trên trình duyệt |
| T3 | Chấm tự động theo trạng thái hệ thống |
| T4 | Biến thể đề riêng cho mỗi học viên |
| T5 | Tổ chức lớp, bài tập, kỳ thi đồng loạt, bảng điểm |
| T6 | Cấp phát tài nguyên cho kỳ thi: đặt trước, kéo sẵn image, xử lý mất máy chạy phiên |
| T7 | Cách ly môi trường bằng sandbox nhân riêng hoặc máy ảo |

## 2. Hiện trạng tại Đại học Bách khoa Hà Nội

Trong quá trình học tại trường, nhóm chưa thấy nền tảng nào của trường cho phép sinh viên thực hành và thi quản trị hệ
thống trên môi trường riêng có chấm tự động. Hệ thống trực tuyến hiện có phục vụ chấm bài lập trình theo đầu vào và
đầu ra; bài thực hành hệ thống được làm trên máy cá nhân hoặc phòng máy và được kiểm tra thủ công.

## 3. Giải pháp mã nguồn mở

| Giải pháp | Giấy phép | Hoạt động | Mô tả |
|---|---|---|---|
| Educates (`educates/educates-training-platform`) | Apache-2.0 | Có commit 10/2026, 161 sao | Nền tảng workshop tương tác trên Kubernetes do VMware Tanzu mở mã: mỗi người dùng một phiên cách ly, hướng dẫn từng bước, terminal và VS Code trên trình duyệt; có "examiner" chạy script kiểm tra theo mã thoát |
| Labtainers (`mfthomps/Labtainers`) | Phạm vi công cộng (sản phẩm của Chính phủ Mỹ) | Có commit 09/2026, 316 sao | Khung bài lab an ninh mạng dựa trên Docker của Naval Postgraduate School; tham số hóa bài lab theo từng học viên để hạn chế chia sẻ; chấm tự động dựa trên dấu vết thao tác |
| PrairieLearn (`PrairieLearn/PrairieLearn`) | Không xác định được qua GitHub API | Có commit 10/2026, 519 sao | Hệ thống học và đánh giá theo câu hỏi; workspace là container bền vững mở bằng VS Code, JupyterLab trên trình duyệt, chấm bằng bộ chấm ngoài trên tệp của workspace |
| KYPO Cyber Range Platform (Đại học Masaryk) | MIT | Ngừng phát triển; kế thừa bởi CyberRangeCZ | Cyber range trên OpenStack: máy ảo cho mỗi học viên, bài huấn luyện an ninh mạng; nhóm phát triển đã công bố nghiên cứu về sinh đề cá nhân hóa và phát hiện gian lận |
| Play with Docker (`play-with-docker/play-with-docker`) | MIT | Có commit 08/2026, 3.238 sao | Sân chơi Docker trên trình duyệt, phiên tạm thời; không có bài tập, chấm điểm hay lớp học |
| JupyterHub và KubeSpawner | BSD-3-Clause | Có commit 10/2026 | Cấp môi trường riêng cho mỗi người dùng trên Kubernetes; không có nghiệp vụ bài lab hay chấm điểm |
| CTFd (`CTFd/CTFd`) | Apache-2.0 | Có commit 10/2026, 6.878 sao | Nền tảng thi CTF; chấm bằng cờ (flag) do người chơi nộp; môi trường riêng cho mỗi người chơi phải dùng plugin bên ngoài |
| Autolab (`autolab/Autolab`) | Apache-2.0 | Có commit 10/2026, 808 sao | Quản lý môn học và chấm bài lập trình tự động; chấm tệp nộp, không cấp môi trường tương tác |

## 4. Giải pháp thương mại

| Giải pháp | Mô tả | Hạn chế đối với đối tượng của Crucible |
|---|---|---|
| Killercoda (kế thừa Katacoda) | Kịch bản tương tác Linux, Kubernetes trên trình duyệt, có bước kiểm tra | Dịch vụ đóng mã, không tự cài đặt; không có lớp học, kỳ thi |
| KodeKloud, iximiuz Labs, LabEx | Khóa học và phòng lab trên trình duyệt có kiểm tra | Dịch vụ thu phí cho cá nhân hoặc doanh nghiệp; nội dung do nhà cung cấp soạn |
| Instruqt, Strigo | Phòng lab cho đào tạo sản phẩm của doanh nghiệp | Định hướng doanh nghiệp, chi phí cao |
| SadServers | Bài khắc phục sự cố Linux trên máy chủ thật | Kho bài cố định, không có lớp học |

## 5. Đối chiếu theo tiêu chí

Ký hiệu: ● đáp ứng; ◐ đáp ứng một phần; ○ không đáp ứng hoặc không tìm thấy trong tài liệu.

| Giải pháp | T1 | T2 | T3 | T4 | T5 | T6 | T7 |
|---|---|---|---|---|---|---|---|
| Educates | ● | ● | ◐ | ○ | ○ | ◐ | ◐ |
| Labtainers | ● | ◐ | ● | ● | ○ | ○ | ○ |
| PrairieLearn | ● | ◐ | ● | ◐ | ● | ○ | ○ |
| KYPO CRP | ● | ● | ◐ | ● | ◐ | ○ | ● |
| CTFd | ● | ○ | ○ | ○ | ◐ | ○ | ○ |
| Killercoda | ○ | ● | ● | ○ | ○ | ○ | ◐ |
| Crucible (mục tiêu) | ● | ● | ● | ● | ● | ● | ● |

Ghi chú:
- Educates: examiner chạy script trong chính container của học viên nên học viên có thể can thiệp vào kết quả; tài
  liệu có cơ chế giữ sẵn phiên (reserved sessions) và chọn RuntimeClass, nhưng không có đặt trước theo kỳ thi; không có
  điểm, lớp, kỳ thi.
- Labtainers: chạy bằng Docker trên máy của học viên hoặc máy ảo; container cần đặc quyền; giao diện dòng lệnh, không
  có hệ thống lớp học và thi trực tuyến.
- PrairieLearn: mạnh về quản lý môn học và đánh giá; workspace thiên về soạn mã nguồn, chấm trên tệp; không hướng tới
  quản trị hệ thống với quyền root và kỳ thi đồng loạt có đặt trước tài nguyên.
- KYPO CRP: cách ly bằng máy ảo trên OpenStack, yêu cầu hạ tầng lớn (bản thử nghiệm cần khoảng 48 GB bộ nhớ cho riêng
  OpenStack); đã ngừng phát triển.

## 6. Khoảng trống và định vị của Crucible

Chưa có giải pháp mã nguồn mở nào kết hợp đồng thời:

1. Môi trường Linux có quyền root, cách ly bằng sandbox nhân riêng, chạy bằng Docker trên một hoặc hai máy mà một khoa tự
   vận hành được (khác KYPO dựa trên OpenStack, Educates dựa trên Kubernetes).
2. Chấm theo trạng thái hệ thống bằng bộ công cụ học viên không can thiệp được (khác Educates).
3. Biến thể đề theo từng học viên và tín hiệu gian lận dựa trên biến thể (kế thừa ý tưởng của Labtainers và KYPO,
   tích hợp vào thi trực tuyến).
4. Nghiệp vụ lớp học, lộ trình, kỳ thi đồng loạt với đặt trước tài nguyên, bù giờ khi sự cố hạ tầng.

Crucible kế thừa có chủ đích các điểm sau:

| Kế thừa từ | Nội dung |
|---|---|
| Educates, Killercoda | Hướng dẫn từng bước kèm terminal; kiểm tra bước bằng mã thoát của script |
| Labtainers | Tham số hóa bài lab theo học viên để hạn chế chia sẻ kết quả |
| KYPO CRP | Dùng biến thể làm căn cứ phát hiện gian lận |
| JupyterHub | Cấp môi trường riêng cho mỗi người dùng trên Kubernetes |

## Nguồn tham khảo

- Educates: https://docs.educates.dev/ ; https://github.com/educates/educates-training-platform
- Labtainers: https://github.com/mfthomps/Labtainers ; https://www.caecommunity.org/sites/default/files/Labtainers-CAE-TechTalk.pdf
- PrairieLearn: https://docs.prairielearn.com/workspaces/
- KYPO CRP: https://crp.kypo.muni.cz/ ; https://docs.crp.kypo.muni.cz/installation-guide/openstack-requirements/
- Nghiên cứu về sinh đề cá nhân hóa và phát hiện gian lận (SIGCSE 2022): https://arxiv.org/pdf/2201.01154
- Killercoda: https://alternativeto.net/software/killercoda/about

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 26/10/2026 | Ban hành lần đầu |
