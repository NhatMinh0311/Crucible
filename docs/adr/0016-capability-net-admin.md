# ADR-0016: Cấp capability NET_ADMIN cho bài lab mạng khi chạy trong gVisor

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận (chờ kết quả thử nghiệm S4) |
| Ngày | 18/10/2026 |

## Bối cảnh
Dạng bài mạng trong một máy (CR-DOC-01 mục 5.4) gồm cấu hình tường lửa cục bộ (`iptables`, `nftables`), thêm địa chỉ
và tuyến (`ip addr`, `ip route`). Các lệnh này cần capability `NET_ADMIN`, không có trong bộ capability mặc định của
container runtime. ADR-0003 và CR-DOC-05 (T1) quy định không thêm capability ngoài mặc định, nên các bài này thất bại
với lỗi không đủ quyền.

Trong gVisor, ngăn xếp mạng của môi trường do nhân riêng của sandbox (netstack) xử lý. Capability của tiến trình trong
sandbox chỉ có hiệu lực với nhân này: thay đổi tường lửa hay địa chỉ chỉ tác động lên giao diện mạng ảo của phiên,
không tác động lên bảng `iptables` và giao diện mạng của node. NetworkPolicy được thực thi trên node, bên ngoài sandbox.

## Quyết định
- Đặc tả bài lab có trường `environment.capabilities`, hiện chỉ nhận `NET_ADMIN`, mặc định rỗng. Bài lab không khai
  báo thì pod phiên không có capability bổ sung.
- Operator thêm capability đã khai báo vào `securityContext.capabilities.add` của container phiên.
- Chính sách `session-pod-hardening` chỉ cho phép `NET_ADMIN` trong `capabilities.add` và vẫn bắt buộc
  `runtimeClassName: gvisor`.
- Cấu hình `SESSION_ALLOWED_CAPABILITIES` của `runtime-api` (mặc định `NET_ADMIN`) quyết định capability được cấp trên
  cụm. Yêu cầu tạo phiên có capability ngoài danh sách bị từ chối với 422 `ENVIRONMENT_NOT_ALLOWED`.
- Nếu thử nghiệm S1 thất bại và cụm chuyển sang runc (CR-DOC-06 RK1), `SESSION_ALLOWED_CAPABILITIES` để rỗng và chính
  sách từ chối mọi capability bổ sung; bài lab khai báo `capabilities` không qua được kiểm thử xuất bản.
- Pod ấm khớp phiên theo image, mức tài nguyên và tập capability. Đặt trước cho kỳ thi truyền `warmCapabilities` để pod
  ấm được tạo đúng tập capability.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Không cấp, bỏ bài tường lửa | Mất một nội dung cơ bản của học phần quản trị hệ thống |
| Cấp `NET_ADMIN` cho mọi pod phiên | Trái nguyên tắc quyền tối thiểu; bài không cần vẫn nhận quyền |
| Mô phỏng tường lửa bằng tệp cấu hình, chấm theo nội dung tệp | Không kiểm chứng được hệ quả thực; học viên không quan sát được tác dụng của quy tắc |

## Hệ quả
- (+) Bài tường lửa cục bộ và cấu hình địa chỉ, tuyến làm được trong MVP; cùng cơ chế dùng cho môi trường nhiều máy
  (ADR-0015).
- (+) Capability chỉ cấp theo từng bài lab và bị chính sách giới hạn ở `NET_ADMIN`.
- (−) Mức hỗ trợ `iptables`, `nftables` của gVisor chưa đầy đủ; thử nghiệm S4 xác định lệnh nào dùng được để ghi vào
  hướng dẫn soạn bài.
- (−) Học viên có `NET_ADMIN` đổi được địa chỉ của giao diện mạng; kiểm thử cách ly I12 kiểm chứng việc này không vượt
  qua được NetworkPolicy.
- (−) Pod ấm tách theo tập capability nên kho ấm phân mảnh thêm.
