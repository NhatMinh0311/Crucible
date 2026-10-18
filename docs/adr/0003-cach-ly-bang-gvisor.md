# ADR-0003: Cách ly môi trường học viên bằng gVisor, NetworkPolicy chặn mặc định và chính sách khi tạo pod

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận; thử nghiệm S1 đạt ngày 18/10/2026 (`docs/spikes/S1.md`) |
| Ngày | 17/10/2026 |

## Bối cảnh
Học viên có quyền root trong môi trường và chạy lệnh tùy ý. Container dùng chung nhân hệ điều hành với node; một lỗ
hổng nhân đủ để thoát ra node và chiếm cụm. Cụm chạy trên máy ảo trong một máy tính cá nhân, không có ảo hóa lồng nhau
ổn định để chạy máy ảo cho mỗi phiên.

## Quyết định
- Pod phiên chạy với `RuntimeClass gvisor` (runsc, nền tảng systrap), chỉ trên node phiên được gắn nhãn.
- NetworkPolicy mặc định chặn mọi lưu lượng vào và ra trong `crucible-sessions`; chỉ mở đường ra máy chủ gói nội bộ cho
  bài lab `network: mirror`.
- Không token tài khoản dịch vụ, không capability bổ sung (ngoại lệ `NET_ADMIN` theo ADR-0016), không `hostPath`, giới hạn bộ nhớ, CPU, số tiến trình,
  dung lượng tạm.
- Các ràng buộc trên được thực thi bằng ValidatingAdmissionPolicy, độc lập với mã của operator (CR-DOC-05 mục 5).

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| runc mặc định | Chung nhân với node; không đủ cho môi trường cấp root cho người dùng |
| Kata Containers hoặc Firecracker | Cần ảo hóa lồng nhau trong máy ảo KVM; tốn bộ nhớ hơn cho mỗi phiên |
| Máy ảo cho mỗi phiên (KubeVirt) | Thời gian khởi động và bộ nhớ vượt khả năng của cụm |

## Hệ quả
- (+) Nhân riêng trong không gian người dùng cho mỗi phiên; chi phí bộ nhớ thấp (khoảng 64 MiB mỗi phiên).
- (+) Chính sách thực thi ở API server nên lỗi của operator không làm mất cách ly.
- (−) Một số lời gọi hệ thống và công cụ nhân không có trong gVisor; bài lab về nhân hệ điều hành không làm được.
- (−) Hiệu năng vào ra thấp hơn runc; không ảnh hưởng bài thực hành cơ bản.
- (−) Vượt giới hạn bộ nhớ hoặc số tiến trình làm dừng cả sandbox thay vì chỉ tiến trình vi phạm; phiên mất môi
  trường (S1).
