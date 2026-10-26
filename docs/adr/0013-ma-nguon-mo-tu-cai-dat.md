# ADR-0013: Phát hành mã nguồn mở để mỗi đơn vị tự cài đặt; bản dịch vụ nhiều tổ chức là hướng phát triển

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Đối tượng chính là học phần thực hành tại trường đại học (CR-DOC-01 mục 1.1). Môi trường thực hành tiêu tốn tài nguyên
tính toán theo số học viên; dữ liệu bài thi thường được yêu cầu lưu trong đơn vị. Nhóm cũng cân nhắc khả năng cung cấp
một bản dịch vụ thu phí trong tương lai.

## Quyết định
- Mỗi bản cài đặt phục vụ một đơn vị đào tạo, trên hạ tầng của đơn vị. Nhóm phát hành phiên bản có đánh số kèm cấu
  hình Docker Compose, mã dựng hạ tầng tham chiếu và tài liệu cài đặt, nâng cấp (CR-DOC-04 mục 6).
- Mã nguồn phát hành theo giấy phép GNU AGPL-3.0: đơn vị tự cài đặt dùng miễn phí; bên cung cấp Crucible đã sửa đổi dưới dạng
  dịch vụ mạng phải công bố mã nguồn sửa đổi, bảo vệ hướng phát triển bản dịch vụ của nhóm.
- Bản dịch vụ nhiều tổ chức, gói và thanh toán không thực hiện trong đồ án. Thiết kế giữ khả năng mở rộng: module
  tham chiếu nhau bằng mã mờ; hạn mức tài nguyên gắn với tài khoản hạn mức, có thể ánh xạ tới tổ chức sau này.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Chỉ cung cấp dịch vụ do nhóm vận hành, bán gói theo số học viên | Thêm multi-tenant cho cả ba module và module thanh toán, lệch trọng tâm khỏi nghiệp vụ thực hành và chấm điểm; chi phí tính toán do nhóm chịu |
| Giấy phép Apache-2.0 | Bên khác có thể cung cấp dịch vụ từ mã nguồn mà không chia sẻ sửa đổi |
| Thực hiện đồng thời cả hai trong đồ án | Phạm vi vượt khả năng của nhóm trong thời gian học phần |

## Hệ quả
- (+) Phát hành, cài đặt và nâng cấp cho người vận hành bên ngoài trở thành một bài toán vận hành cần kiểm chứng (Q13, Q14).
- (+) Không phát sinh multi-tenant và thanh toán trong phạm vi đồ án.
- (−) Cần duy trì tài liệu cài đặt và tương thích nâng cấp giữa các phiên bản.
- (−) Chuyển sang bản dịch vụ sau này cần thêm tầng tổ chức trong cả ba module.
