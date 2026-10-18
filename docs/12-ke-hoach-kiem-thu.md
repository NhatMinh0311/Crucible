# Kế hoạch kiểm thử

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-DOC-12 |
| Phiên bản | 1.3 |
| Ngày ban hành | 18/10/2026 |
| Tài liệu liên quan | CR-DOC-01 mục 7, 8; CR-DOC-04 mục 9, 11; CR-DOC-05 mục 6; CR-DOC-06 mục 6; tài liệu module |

## 1. Mục tiêu và phạm vi

Kế hoạch kiểm thử bảo đảm mọi tiêu chí chấp nhận (CR-DOC-01 mục 7) và mọi yêu cầu chất lượng (CR-DOC-01 mục 8) có bằng
chứng kiểm thử truy vết được. Trưởng nhóm phụ trách kế hoạch, kiểm thử hệ thống, phi chức năng và bảo mật; chủ module
phụ trách kiểm thử trong module.

## 2. Các mức kiểm thử

| Mức | Mục đích | Công cụ | Môi trường | Người thực hiện | Thời điểm chạy |
|---|---|---|---|---|---|
| Đơn vị | Quy tắc nghiệp vụ `BR-*`, tính toán (điểm, biến thể, sức chứa, giờ máy) | `go test`; JUnit 5; Vitest | Máy phát triển, CI | Chủ module | Mỗi pull request |
| Tích hợp | Truy cập CSDL, outbox và inbox, broker, Kubernetes API (envtest) | Testcontainers, envtest | CI | Chủ module | Mỗi pull request |
| Hợp đồng | Bên cung cấp khớp OpenAPI, sự kiện khớp AsyncAPI, đặc tả bài lab hợp lệ | Schemathesis, kiểm tra lược đồ sự kiện, `validate.py` | CI | Chủ module; trưởng nhóm duy trì job | Mỗi pull request |
| Thành phần giao diện | Hành vi màn hình, trạng thái lỗi | Vitest, Testing Library | CI | Chủ module | Mỗi pull request |
| Đầu cuối | Luồng nghiệp vụ qua ba module trên cụm thật | Playwright | k3d trong CI; cụm chuẩn | Trưởng nhóm | Mỗi pull request từ mốc MA; hằng đêm trên cụm chuẩn |
| Phi chức năng | Q1–Q5, Q8–Q14 | k6, Chaos Mesh, `virsh`, kịch bản cài đặt và nâng cấp | Cụm chuẩn | Trưởng nhóm | Theo mốc MB, MC, MD |
| Bảo mật | Q6, Q7, CR-DOC-05 mục 6 | Bộ kiểm thử cách ly, ma trận phân quyền, OWASP ZAP, Trivy, gitleaks | CI; cụm chuẩn | Trưởng nhóm | Mỗi pull request (quét tĩnh); sau mỗi lần thay đổi cụm (cách ly) |

## 3. Bất biến hệ thống

Các bất biến dưới đây phải đúng ở mọi thời điểm, kể cả khi có sự cố. Mỗi bất biến được kiểm thử ở mức thấp nhất có thể
và được kiểm tra lại sau mỗi thí nghiệm chaos bằng script đối chiếu dữ liệu `tests/invariants/`.

| Mã | Bất biến | Kiểm thử |
|---|---|---|
| INV1 | Mỗi sinh viên có tối đa một phiên luyện tập và một phiên thi chưa kết thúc | Tích hợp M1 (yêu cầu đồng thời); script đối chiếu |
| INV2 | Không tồn tại pod hoặc ConfigMap phiên không thuộc `LabSession` chưa kết thúc quá 2 phút | envtest M1; chỉ số `crucible_orphan_resources`; chaos CH2, CH3 |
| INV3 | Mỗi phiên được ghi giờ máy đúng một lần, bằng thời gian sẵn sàng thực tế | Tích hợp M1; script đối chiếu sổ hạn mức với lịch sử phiên |
| INV4 | Một lượt làm bài giữ một biến thể suốt vòng đời; hai lượt khác nhau của cùng sinh viên trong luyện tập có biến thể khác nhau khi miền tham số cho phép | Đơn vị M2 |
| INV5 | Điểm các bước đã đạt của lượt luyện tập không giảm khi tiếp tục hoặc khi phiên bị mất | Tích hợp M2; chaos CH1 |
| INV6 | Điểm thi bằng kết quả chấm cuối; mọi thay đổi sau đó có bản ghi điều chỉnh với điểm cũ, người sửa, lý do | Tích hợp M2, M3 |
| INV7 | Mọi sự kiện đã ghi vào outbox được giao tới hàng đợi của bên nhận; mỗi sự kiện được xử lý đúng một lần | Tích hợp mỗi module; chaos CH4 |
| INV8 | Từ pod phiên không kết nối được tới thành phần nền tảng, pod khác, Kubernetes API, Internet (kể cả qua máy chủ gói) | Bộ kiểm thử cách ly I1–I12 |
| INV9 | Kết quả lệnh kiểm tra không phụ thuộc công cụ, biến môi trường, cấu hình shell trong môi trường | Tích hợp M2 với M1 (AC-C2.4) |
| INV10 | Sinh viên chỉ đọc được dữ liệu của mình; điểm và bản ghi thi chỉ hiển thị sau khi công bố | Ma trận phân quyền API |
| INV11 | Tổng chỗ đặt trước chồng lấn thời gian không vượt 80% sức chứa | Đơn vị M1 |
| INV12 | Tổng thời gian bù giờ của một sinh viên trong một kỳ thi không vượt 15 phút | Đơn vị M3 |
| INV13 | Tổng hạn mức của các lớp không vượt trần của giảng viên; tổng chỗ × giờ đặt trước trong tháng không vượt ngân sách | Tích hợp M1 (yêu cầu đồng thời) |

## 4. Truy vết tiêu chí chấp nhận

| Nhóm tiêu chí | Đơn vị, tích hợp | Hợp đồng | Đầu cuối | Phi chức năng, bảo mật |
|---|---|---|---|---|
| AC-R1.x Cấp phiên, vòng đời | M1 | `runtime-api` | E2E-01 | Q1, Q2 (k6) |
| AC-R2.x Hàng đợi, hạn mức, đặt trước | M1 | `runtime-api` | E2E-02 | Q2, Q3 |
| AC-R3.x Terminal, xem trực tiếp, bản ghi | M1 | Giao thức terminal | E2E-01, E2E-02 | Kiểm thử vé (CR-DOC-05 mục 6.2) |
| AC-R4.x Mất node, tài nguyên mồ côi | M1 (envtest) | | | CH1–CH3 |
| AC-C1.x Soạn, chạy thử, xuất bản | M2 | `lab-spec` | E2E-03 | |
| AC-C2.x Biến thể, chấm điểm, khôi phục | M2 | `content-api`, sự kiện | E2E-01, E2E-02 | INV9 |
| AC-L1.x Lớp, lộ trình, bài tập | M3 | | E2E-01 | |
| AC-L2.x Kỳ thi | M3 | Sự kiện | E2E-02 | Q2, CH1 |
| AC-L3.x Gian lận, phúc khảo | M3 | `content-api` | E2E-04 | |

Mỗi kiểm thử mang mã tiêu chí trong tên hoặc chú thích (AGENTS.md mục 2). Trước mỗi mốc, script `tests/trace/report`
quét mã kiểm thử và liệt kê tiêu chí chưa có kiểm thử.

## 5. Kịch bản đầu cuối

| Mã | Kịch bản | Bước chính |
|---|---|---|
| E2E-01 | Luyện tập theo lộ trình | Giảng viên tạo lớp, nhập CSV, tạo lộ trình 2 bài; sinh viên đăng nhập lần đầu được gắn vào lớp; làm bài 1 qua terminal, kiểm tra bước, mở gợi ý, nộp; bài 2 được mở; bỏ dở bài 2, phiên hết hạn, tiếp tục và thấy tiến độ được khôi phục |
| E2E-02 | Kỳ thi | Lên lịch thi với chế độ `RECORDED_ONLY`; 5 sinh viên vào thi, 1 sinh viên vào sau hạn chót bị từ chối; giảng viên xem trực tiếp; một phiên bị xóa giả lập sự cố, sinh viên tiếp tục và được bù giờ; một sinh viên nộp sớm; hết giờ các sinh viên còn lại được nộp tự động; công bố điểm; sinh viên xem điểm từng bước và bản ghi |
| E2E-03 | Soạn và xuất bản bài lab | Giảng viên soạn bằng form, chuyển sang YAML và ngược lại; chạy thử; xuất bản thất bại do lời giải sai; sửa và xuất bản thành công; giảng viên khác nhân bản bài chia sẻ |
| E2E-04 | Gian lận và phúc khảo | Hai sinh viên trong kỳ thi, sinh viên A nhập giá trị biến thể của B; tín hiệu `VARIANT_LEAK` xuất hiện; trợ giảng đề xuất, giảng viên lập biên bản hủy kết quả; sinh viên A gửi phúc khảo, giảng viên trả lời; gửi lần hai bị từ chối |

## 6. Dữ liệu kiểm thử

| Bộ dữ liệu | Nội dung | Vị trí |
|---|---|---|
| Lớp mẫu | 2 lớp, 2 giảng viên, 1 trợ giảng, 60 sinh viên giả (tệp CSV) | `tests/data/` |
| Bài lab mẫu | Tối thiểu 10 bài; mỗi dạng trong phạm vi MVP (CR-DOC-01 mục 5.4) có ít nhất hai bài; có ít nhất một bài theo mỗi mẫu kịch bản sự cố tái phát, truy dấu xâm nhập, chuỗi manh mối; ít nhất một bài có bước chẩn đoán; ít nhất hai bài dùng cho thi | `labs/` |
| Tài khoản kiểm thử tải | 100 tài khoản cục bộ trong Keycloak, tạo bằng script | `tests/load/` |

Dữ liệu kiểm thử không chứa thông tin cá nhân thật.

## 7. Tiêu chí theo mốc

| Mốc | Tiêu chí vào | Tiêu chí ra |
|---|---|---|
| MA | Hợp đồng v1 đã merge | E2E-01 đạt trên môi trường phát triển; mọi kiểm thử mức pull request đạt |
| MB | MA đạt; cụm chuẩn hoạt động | E2E-01 đạt trên cụm chuẩn; bộ kiểm thử cách ly đạt; Q1 đo được |
| MC | MB đạt | E2E-02, E2E-03, E2E-04 đạt; CH1 đạt |
| MD | MC đạt | Báo cáo Q1–Q14; CH1–CH7; không còn lỗi mức nghiêm trọng hoặc cao đang mở; báo cáo truy vết không còn tiêu chí thiếu kiểm thử |

## 8. Quản lý lỗi

| Mức | Định nghĩa | Thời hạn xử lý |
|---|---|---|
| Nghiêm trọng | Mất dữ liệu điểm, vượt ranh giới cách ly, kỳ thi không diễn ra được | Trước khi merge thêm tính năng |
| Cao | Một luồng chính không hoàn thành được; sai điểm | Trong mốc hiện tại |
| Trung bình | Luồng phụ lỗi, có cách khắc phục tạm | Trước mốc MD |
| Thấp | Giao diện, văn bản | Khi có thời gian |

Lỗi được ghi thành GitHub Issue với nhãn `bug`, nhãn mức độ, nhãn module, mã tiêu chí bị vi phạm. Lỗi phát hiện bởi
kiểm thử phi chức năng hoặc chaos kèm biên bản thí nghiệm.

## 9. Báo cáo kiểm thử

| Báo cáo | Nội dung | Vị trí |
|---|---|---|
| Kết quả CI | Kết quả mỗi lần chạy, gắn với pull request | GitHub Actions |
| Báo cáo truy vết | Tiêu chí và kiểm thử tương ứng, tiêu chí còn thiếu | `tests/trace/report.md`, cập nhật mỗi mốc |
| Báo cáo tải, sức chứa | Số đo Q1–Q3, cấu hình, phiên bản | `tests/load/results/` |
| Biên bản chaos | Giả thuyết, cách làm, kết quả, đối chiếu bất biến | `tests/chaos/reports/` |
| Báo cáo bảo mật | Kết quả kiểm thử cách ly, ma trận phân quyền, quét tự động | `tests/security/reports/` |

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| 1.0 | 17/10/2026 | Ban hành lần đầu |
| 1.1 | 17/10/2026 | INV13; INV8 gồm kiểm thử I10, I11 |
| 1.2 | 18/10/2026 | Bộ bài mẫu theo dạng bài và mẫu kịch bản |
| 1.3 | 18/10/2026 | INV8 gồm kiểm thử I12 (`NET_ADMIN`) |
