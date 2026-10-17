# ADR-0009: Hạ tầng bằng máy ảo KVM dựng bằng Terraform, cụm k3s một node nền tảng và ba node phiên

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Nhóm không có ngân sách cloud dài hạn. Máy của trưởng nhóm có 16 luồng CPU, 38 GB bộ nhớ, hỗ trợ KVM. Hệ thống cần cụm
nhiều node để kiểm chứng mất node, phân tách node nền tảng và node phiên, và dựng lại được từ mã.

## Quyết định
- Bốn máy ảo KVM/libvirt (CR-DOC-04 mục 3.2) dựng bằng Terraform (provider `dmacvicar/libvirt`) và cloud-init.
- k3s: một server (đồng thời chạy thành phần nền tảng), ba agent làm node phiên có taint.
- Máy của thành viên tham gia tùy chọn làm node `volatile` qua Tailscale (CR-DOC-04 mục 4).

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Cloud với tín dụng sinh viên | Tín dụng có hạn; cụm phải tắt khi không dùng; có thể chuyển sang sau vì Terraform và k3s dùng lại được |
| k3d hoặc kind trên một máy | Node là container, không giả lập được mất node, không chạy gVisor đúng như máy thật |
| Gộp ba laptop của thành viên thành cụm | Máy không bật liên tục, khác mạng, khác hệ điều hành |
| Kubernetes đầy đủ bằng kubeadm | Tốn bộ nhớ và công vận hành hơn k3s, không thêm giá trị cho bài toán |

## Hệ quả
- (+) Toàn quyền kiểm soát cụm, thí nghiệm được mất node bằng `virsh destroy`.
- (+) Chi phí bằng 0; dựng lại bằng một lệnh.
- (−) Toàn bộ hệ thống phụ thuộc một máy vật lý (CR-DOC-02 mục 11).
- (−) Máy vật lý của trưởng nhóm phải bật khi thành viên cần dùng cụm chuẩn.
