# ADR-0009: Hạ tầng gồm hai máy ảo KVM dựng bằng script và cloud-init

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 26/10/2026 |

## Bối cảnh
Nhóm không có ngân sách cloud dài hạn. Máy của trưởng nhóm có 16 luồng CPU, 38 GB bộ nhớ, hỗ trợ KVM. Môi trường phiên
phải tách khỏi nền tảng (ADR-0003) và dựng lại được từ mã (Q12).

## Quyết định
- Hai máy ảo KVM/libvirt (CR-DOC-04 mục 3.2): `cr-server` chạy nền tảng và dịch vụ, `cr-lab` chạy Docker Engine với
  gVisor cho phiên.
- Máy ảo được dựng bằng script `infra/create-vms.sh` (`virt-install`) và cloud-init (`infra/cloud-init/`): cài Docker,
  gVisor, cấu hình Docker của máy phiên, tài khoản SSH cho `runtime-api`.
- `cr-lab` chỉ nhận kết nối SSH từ `cr-server` và từ máy quản trị.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Một máy ảo cho mọi thành phần | Thoát khỏi sandbox chạm tới CSDL và Keycloak; phiên tranh bộ nhớ với nền tảng |
| Terraform với provider libvirt | Thêm công cụ phải học cho hai máy ảo; script và cloud-init đủ để dựng lại |
| Cloud với tín dụng sinh viên | Tín dụng có hạn; máy phải tắt khi không dùng; có thể chuyển sang sau vì cloud-init dùng lại được |

## Hệ quả
- (+) Chi phí bằng 0; dựng lại bằng một lệnh.
- (+) Thí nghiệm mất máy phiên bằng `virsh destroy cr-lab` hoặc dừng Docker Engine.
- (−) Toàn bộ hệ thống phụ thuộc một máy vật lý (CR-DOC-02 mục 11).
- (−) Máy vật lý của trưởng nhóm phải bật khi thành viên cần dùng môi trường chuẩn.
