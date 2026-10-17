# ADR-0001: Chia hệ thống thành ba module theo nghiệp vụ, mỗi module sở hữu dữ liệu riêng

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Nhóm có ba thành viên, mỗi thành viên phụ trách trọn một module từ giao diện tới triển khai và được đánh giá theo module.
Mỗi module cần có nghiệp vụ riêng đủ độ phức tạp, phát triển và triển khai được độc lập, và ít phụ thuộc tiến độ của nhau.

## Quyết định
- Chia thành ba module theo năng lực nghiệp vụ: M1 phiên thực hành và tài nguyên, M2 bài lab và chấm điểm, M3 lớp học,
  lộ trình và thi (CR-DOC-01 mục 4).
- Mỗi module là một hoặc vài dịch vụ triển khai độc lập, một CSDL logic riêng trên cụm PostgreSQL dùng chung, một tài
  khoản CSDL chỉ có quyền trên CSDL của mình.
- Module giao tiếp qua hợp đồng công bố tại `contracts/`; phụ thuộc đồng bộ chỉ đi một chiều M3 → M2 → M1.
- Giao diện là một ứng dụng React; mỗi module sở hữu một thư mục tính năng, khung chung do trưởng nhóm sở hữu.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Chia theo tầng (backend, frontend, hạ tầng) | Trái yêu cầu chia theo module; một tính năng phụ thuộc nhiều người |
| Một ứng dụng nguyên khối chia gói theo module | Không có triển khai độc lập, không thể hiện được bài toán vận hành nhiều dịch vụ |
| Mỗi module một cụm PostgreSQL | Không đủ tài nguyên; ranh giới dữ liệu đã được bảo đảm bằng CSDL và tài khoản riêng |
| Micro-frontend cho mỗi module | Chi phí cấu hình cao so với quy mô ba module |

## Hệ quả
- (+) Mỗi thành viên có phạm vi trách nhiệm và sản phẩm rõ ràng, đánh giá được riêng.
- (+) M1 không gọi đồng bộ module nào, phát triển và kiểm thử độc lập ngay từ đầu.
- (−) Dữ liệu cần hiển thị chéo phải được sao chép qua sự kiện; nhất quán sau một khoảng trễ.
- (−) Thay đổi hợp đồng cần sự đồng ý của nhiều thành viên.
