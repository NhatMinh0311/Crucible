# Thiết kế module <Mx> — <tên module>

| Thuộc tính | Giá trị |
|---|---|
| Mã tài liệu | CR-MOD-<Mx> |
| Phiên bản | 0.1 |
| Ngày ban hành | <dd/mm/yyyy> |
| Người phụ trách | <họ tên, tài khoản GitHub> |
| Tài liệu liên quan | CR-DOC-01 mục 5, 7 (chức năng, tiêu chí chấp nhận của module); CR-DOC-02; CR-DOC-03 |

<!--
Hướng dẫn sử dụng mẫu (xóa khối này khi viết):
- Tài liệu mô tả thiết kế BÊN TRONG module. Hợp đồng với module khác đã có tại CR-DOC-03 và `contracts/`; không lặp lại,
  chỉ tham chiếu.
- Mỗi quy tắc nghiệp vụ có mã `BR-<Mx>.<số>`, dùng để truy vết từ kiểm thử, giống mã tiêu chí chấp nhận.
- Mỗi sơ đồ dùng Mermaid để xem được trên GitHub.
- Viết trước khi viết mã (hạn: cuối tuần 1, CR-DOC-06 mục 5), cập nhật khi thiết kế thay đổi.
-->

## 1. Phạm vi và trách nhiệm

Liệt kê mã chức năng tại CR-DOC-01 mục 5 thuộc module, phần nào thuộc MVP, phần nào để sau.

## 2. Thành phần triển khai

| Thành phần | Công nghệ | Trách nhiệm | Số bản sao |
|---|---|---|---|

## 3. Mô hình dữ liệu

Sơ đồ thực thể – liên kết và mô tả bảng: khóa, ràng buộc, chỉ mục quan trọng. Ghi rõ bảng `outbox`, `inbox`.

## 4. Quy tắc nghiệp vụ và máy trạng thái

Sơ đồ trạng thái của các đối tượng chính; bảng quy tắc:

| Mã | Quy tắc | Tiêu chí chấp nhận liên quan |
|---|---|---|

## 5. API cho giao diện

Operation chỉ phục vụ giao diện của module (đặc tả đầy đủ tại `services/<module>/api/openapi.yaml`), kèm quyền của từng
operation.

## 6. Sự kiện

| Sự kiện | Phát / nhận | Xử lý |
|---|---|---|

## 7. Giao diện

Danh sách màn hình, vai trò sử dụng, luồng chính; phác thảo nếu cần.

## 8. Tác vụ định kỳ

| Tác vụ | Chu kỳ | Nội dung | Khóa phân tán |
|---|---|---|---|

## 9. Kiểm thử

| Mức | Công cụ | Phạm vi |
|---|---|---|

Bảng truy vết: mã tiêu chí chấp nhận và mã quy tắc → kiểm thử.

## 10. Quan sát

Metric (bao gồm chỉ số tại CR-DOC-04 mục 7.2 do module phát), log quan trọng, dashboard của module.

## 11. Cấu hình

| Biến | Ý nghĩa | Mặc định |
|---|---|---|

## Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
