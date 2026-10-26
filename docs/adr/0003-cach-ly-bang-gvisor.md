# ADR-0003: Cách ly môi trường học viên bằng gVisor, mạng riêng mỗi phiên và máy phiên tách khỏi nền tảng

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận (chờ kết quả thử nghiệm S1) |
| Ngày | 26/10/2026 |

## Bối cảnh
Học viên có quyền root trong môi trường và chạy lệnh tùy ý. Container dùng chung nhân hệ điều hành với máy chạy nó; một
lỗ hổng nhân đủ để thoát ra máy đó. Hạ tầng là máy ảo trong một máy tính cá nhân, không có ảo hóa lồng nhau ổn định để
chạy máy ảo cho mỗi phiên.

## Quyết định
- Container phiên chạy bằng runtime `runsc` của gVisor (nền tảng systrap). Docker trên máy phiên đặt `runsc` là runtime
  mặc định, nên mọi container trên máy này chạy trong gVisor kể cả khi `runtime-api` có lỗi; `runtime-api` kiểm tra
  runtime khi khởi động (BR-M1.15).
- Máy phiên `cr-lab` là máy ảo riêng, không chứa CSDL, khóa hay dịch vụ của nền tảng (ADR-0009).
- Phiên `network: none` không có mạng; phiên `network: mirror` dùng mạng Docker `internal` riêng chỉ nối với máy chủ gói
  nội bộ (CR-MOD-M1 mục 3.2).
- Container phiên: không gắn socket Docker, không gắn thư mục của máy, không capability bổ sung (ngoại lệ `NET_ADMIN`
  theo ADR-0016), `no-new-privileges`, giới hạn bộ nhớ, CPU, số tiến trình.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| runc mặc định | Chung nhân với máy phiên; không đủ cho môi trường cấp root cho người dùng |
| Kata Containers hoặc Firecracker | Cần ảo hóa lồng nhau trong máy ảo KVM; tốn bộ nhớ hơn cho mỗi phiên |
| Máy ảo cho mỗi phiên | Thời gian khởi động và bộ nhớ vượt khả năng của hạ tầng |
| Docker chế độ rootless | Hỗ trợ của gVisor với Docker rootless hạn chế; gVisor đã là ranh giới chính |

## Hệ quả
- (+) Nhân riêng trong không gian người dùng cho mỗi phiên; chi phí bộ nhớ thấp (khoảng 64 MiB mỗi phiên, đo ở S5).
- (+) Runtime mặc định của Docker trên máy phiên là lớp bảo vệ độc lập với mã của `runtime-api`.
- (−) Một số lời gọi hệ thống và công cụ nhân không có trong gVisor; bài lab về nhân hệ điều hành không làm được.
- (−) Hiệu năng vào ra thấp hơn runc; không ảnh hưởng bài thực hành cơ bản.
- (−) Vượt giới hạn bộ nhớ hoặc số tiến trình làm dừng cả sandbox thay vì chỉ tiến trình vi phạm; phiên mất môi trường.
- (−) Không có lớp chính sách ở tầng điều phối như admission của Kubernetes; quyền điều khiển Docker của `runtime-api`
  tương đương quyền root trên máy phiên (CR-DOC-05 mục 3).
