# ADR-0015: Môi trường nhiều máy: chuẩn bị hợp đồng, hiện thực sau MVP

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 18/10/2026 |

## Bối cảnh
Bài quản trị mạng giữa nhiều máy (DNS server cho máy client, ứng dụng gọi cơ sở dữ liệu ở máy khác, reverse proxy,
tường lửa giữa hai máy) cần nhiều máy trong một phiên. Thiết kế hiện tại cấp một container cho một phiên (ADR-0002).
Hiện thực nhiều máy làm tăng khối lượng của M1 (nhiều container mỗi phiên, sức chứa tính theo máy, terminal và lệnh
kiểm tra chọn máy) trong khi lộ trình có 11 tuần (CR-DOC-06).

Bài ở tầng liên kết dữ liệu và định tuyến (VLAN, DHCP, định tuyến giữa các mạng con, giao thức định tuyến động) cần
mạng ảo tầng 2 riêng cho từng phiên hoặc máy ảo.

## Quyết định
- MVP chỉ có môi trường một máy. Môi trường nhiều máy (R14) hiện thực sau MVP.
- Hợp đồng có sẵn khái niệm máy để R14 không thay đổi hợp đồng theo cách phá vỡ tương thích:
  - Đặc tả bài lab: trường `machine` của bước và lỗi cài sẵn, giá trị mặc định `main`.
  - `runtime-api`: trường `machine` của `ExecRequest`, tham số `machine` khi cấp vé terminal và vé xem.
  - Hiện chỉ nhận giá trị `main`; R14 mở rộng danh sách giá trị và thêm khai báo các máy trong `environment`.
- Khi hiện thực R14: mỗi máy là một container gVisor; các container của một phiên nối vào mạng Docker `internal` riêng
  của phiên, có tên DNS theo tên máy; sức chứa của phiên bằng tổng chỗ của các máy; tối đa 3 máy.
- Bài mạng tầng liên kết dữ liệu và định tuyến nằm ngoài phạm vi (CR-DOC-01 mục 6).

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Hiện thực R14 trong MVP, tối đa 2 máy | Tăng khối lượng của M1, rủi ro trễ mốc MB, MC |
| Giữ R14 là mở rộng, không chuẩn bị hợp đồng | Khi hiện thực phải đổi hợp đồng mà M2, M3 đang dùng |
| Nhiều tiến trình trong một container mô phỏng nhiều máy | Dùng chung không gian mạng, không mô phỏng được các máy riêng |
| Mạng tầng 2 riêng mỗi phiên hoặc máy ảo | Phức tạp vận hành, vượt phạm vi và hạ tầng của đồ án |

## Hệ quả
- (+) M2, M3 viết mã theo hợp đồng có trường `machine` ngay từ đầu; R14 chỉ mở rộng giá trị cho phép.
- (+) Mạng riêng mỗi phiên đã có cho `network: mirror` (CR-MOD-M1 mục 3.2) dùng lại được cho nhiều máy.
- (−) MVP chưa có bài mạng nhiều máy; dạng bài mạng giới hạn trong một máy (CR-DOC-01 mục 5.4).
