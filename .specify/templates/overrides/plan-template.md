# Kế hoạch: [TÊN TÍNH NĂNG]

| Thuộc tính | Giá trị |
|---|---|
| Đặc tả | `specs/[###-ten-ngan]/spec.md` |
| Module | [M1, M2, M3 hoặc Nền tảng] |
| Ngày lập | [NGÀY] |

<!--
  Kế hoạch trả lời "làm thế nào" trong phạm vi một module. Nghiên cứu, mô hình dữ liệu và cách kiểm chứng ghi ngay
  trong tệp này; không tạo research.md, data-model.md, quickstart.md hay thư mục contracts/ riêng.
-->

## 1. Tóm tắt

[Yêu cầu chính của đặc tả và hướng kỹ thuật được chọn, hai, ba câu]

## 2. Bối cảnh kỹ thuật

| Mục | Giá trị |
|---|---|
| Ngôn ngữ, phiên bản | [Go 1.27 (M1); Java 25, Spring Boot 4 (M2, M3); TypeScript, React 19 (giao diện)] |
| Thư viện chính | [Thư viện và phiên bản; bản ổn định đã phát hành ít nhất hai tuần] |
| Lưu trữ | [Bảng PostgreSQL, bucket SeaweedFS, hoặc không có] |
| Kiểm thử | [Mức đơn vị, tích hợp (Testcontainers), hợp đồng, trình duyệt] |
| Ràng buộc | [Thời gian phản hồi, bộ nhớ, giới hạn từ CR-DOC-04, CR-DOC-05] |

## 3. Kiểm tra theo hiến pháp

*Điều kiện phải đạt trước khi thiết kế và kiểm tra lại sau thiết kế (`.specify/memory/constitution.md`).*

| Nguyên tắc | Đạt | Ghi chú |
|---|---|---|
| I. Yêu cầu dẫn từ tài liệu đã duyệt; không còn điểm chưa làm rõ | [Có/Không] | |
| II. Hợp đồng: không đổi, hoặc thay đổi được tách thành pull request riêng | [Có/Không] | |
| III. Mỗi `AC-*`, `BR-*` được hiện thực có kiểm thử mang mã | [Có/Không] | |
| IV. Chỉ sửa mã trong module; phụ thuộc module khác được ghi ở mục 7 | [Có/Không] | |
| V. Bí mật, log, migration, metric theo quy định | [Có/Không] | |

## 4. Cấu trúc mã

```text
[Cây thư mục thực tế của module, chỉ gồm tệp thêm hoặc sửa, mỗi dòng kèm vai trò]
```

## 5. Lựa chọn kỹ thuật

| Vấn đề | Lựa chọn | Lý do | Phương án khác đã xét |
|---|---|---|---|
| [Vấn đề cần quyết định] | [Lựa chọn] | [Lý do] | [Phương án bị loại và vì sao] |

## 6. Dữ liệu

<!-- Bảng, cột chính, chỉ mục, ràng buộc; máy trạng thái nếu có. Xóa mục nếu tính năng không có dữ liệu. -->

[Mô hình dữ liệu, migration cần thêm]

## 7. Hợp đồng và phụ thuộc

| Hạng mục | Nội dung |
|---|---|
| API, sự kiện dùng tới | [Operation trong `contracts/*.yaml`, sự kiện trong `contracts/events.asyncapi.yaml`] |
| Thay đổi hợp đồng | [Không có, hoặc mô tả thay đổi và pull request riêng] |
| Phụ thuộc module khác | [Issue, trạng thái, cách làm khi chưa có (mock)] |

## 8. Kiểm thử và kiểm chứng

| Mức | Nội dung | Mã tiêu chí |
|---|---|---|
| Đơn vị | [Nội dung] | [AC-*, BR-*] |
| Tích hợp | [Nội dung] | [AC-*] |
| Hợp đồng | [Kiểm thử bên cung cấp, kiểm tra sự kiện theo lược đồ] | |
| Thủ công trên môi trường chuẩn | [Các bước, kết quả mong đợi] | [SC-*] |

## 9. Sai khác cần giải trình

<!-- Chỉ điền khi mục 3 có nguyên tắc không đạt. -->

| Sai khác | Vì sao cần | Phương án đơn giản hơn bị loại vì |
|---|---|---|
