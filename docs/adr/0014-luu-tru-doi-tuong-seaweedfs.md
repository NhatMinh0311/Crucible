# ADR-0014: Lưu trữ đối tượng bằng SeaweedFS thay cho MinIO

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận; thay thế lựa chọn MinIO trong CR-DOC-02, CR-DOC-04 phiên bản 1.0 |
| Ngày | 17/10/2026 |

## Bối cảnh
Hệ thống cần kho lưu trữ đối tượng có giao diện S3 cho bản ghi phiên, tệp đính kèm bài lab và đích sao lưu của
CloudNativePG, chạy được trong cụm và trên máy vật lý. Thiết kế ban đầu chọn MinIO. Khi khởi tạo môi trường phát triển,
nhóm xác minh rằng kho mã `minio/minio` đã được lưu trữ (archived) và image `minio/minio` không còn trên Docker Hub.

## Quyết định
- Dùng SeaweedFS (giấy phép Apache-2.0) với cổng S3, chế độ một tiến trình (`weed server -s3`) cho môi trường phát triển
  và cụm chuẩn; Helm chart chính thức của SeaweedFS khi triển khai lên cụm.
- Dịch vụ truy cập lưu trữ đối tượng qua thư viện S3 chuẩn (AWS SDK) và cấu hình điểm truy cập, không phụ thuộc
  thư viện riêng của nhà cung cấp; đơn vị tự cài có thể thay bằng bất kỳ dịch vụ tương thích S3 nào.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| MinIO | Ngừng phân phối image và lưu trữ kho mã |
| Garage | Phù hợp cụm nhỏ nhưng cần thêm bước cấu hình bố cục cụm khi khởi tạo; ít tài liệu tích hợp với CloudNativePG |
| RustFS | Mới đạt phiên bản 1.0, lịch sử vận hành ngắn |

## Hệ quả
- (+) Image được duy trì, giấy phép tương thích AGPL-3.0.
- (+) Dùng giao diện S3 chuẩn nên đổi nhà cung cấp chỉ là đổi cấu hình.
- (−) Tính năng S3 nâng cao (khóa đối tượng, vòng đời) cần được kiểm tra với SeaweedFS khi dùng tới.
