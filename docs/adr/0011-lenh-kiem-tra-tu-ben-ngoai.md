# ADR-0011: Lệnh kiểm tra chạy bằng exec với bộ công cụ liên kết tĩnh trên volume chỉ đọc

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Học viên có quyền root trong môi trường nên có thể thay `curl`, `grep`, `id` bằng script luôn trả kết quả thành công,
đặt biến môi trường, sửa tệp cấu hình shell hoặc `/etc/ld.so.preload` để làm sai kết quả chấm (Q7, AC-C2.4).

## Quyết định
- M1 gắn vào mỗi container phiên một volume chỉ đọc tại `/opt/crucible/bin` chứa BusyBox, `curl`, `jq` liên kết tĩnh.
- Lệnh kiểm tra được M2 gửi qua `runtime-api`, chạy bằng exec với `/opt/crucible/bin/sh`, `PATH=/opt/crucible/bin`, biến
  môi trường rỗng, không nạp tệp cấu hình shell.
- Học viên không có `CAP_SYS_ADMIN` nên không gắn đè hoặc gắn lại volume; công cụ liên kết tĩnh không bị ảnh hưởng bởi
  `ld.so.preload`.
- Lệnh kiểm tra không lưu lại trong môi trường.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Dùng công cụ có sẵn trong image | Học viên thay thế được |
| Container kiểm tra riêng chạy cạnh container phiên | Không đọc được hệ thống tệp của container học viên nếu không chia sẻ volume gốc; phức tạp hơn với gVisor |
| Agent chấm điểm chạy liên tục trong môi trường | Học viên dừng hoặc sửa được tiến trình agent |

## Hệ quả
- (+) Kết quả chấm không bị ảnh hưởng bởi thay đổi công cụ và cấu hình shell trong môi trường.
- (−) Lệnh kiểm tra chỉ dùng được công cụ trong bộ công cụ; kiểm tra phụ thuộc chương trình của bài lab (ví dụ `nginx -t`)
  phải thay bằng quan sát trạng thái (cổng, tiến trình, tệp).
- (−) Học viên vẫn có thể tạo trạng thái giả mà lệnh kiểm tra quan sát được (CR-DOC-05 mục 7).
