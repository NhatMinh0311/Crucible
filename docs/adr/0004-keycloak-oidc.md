# ADR-0004: Xác thực bằng Keycloak (OIDC), giao tiếp giữa dịch vụ bằng client credentials có scope

| Thuộc tính | Giá trị |
|---|---|
| Trạng thái | Đã chấp nhận |
| Ngày | 17/10/2026 |

## Bối cảnh
Ba dịch vụ của ba module cần cùng một cơ chế xác thực người dùng và một cách để dịch vụ gọi dịch vụ với quyền tối thiểu
(ví dụ chỉ M2 được chạy lệnh trong phiên). Thời gian của nhóm nên dành cho nghiệp vụ của module thay vì tự xây chức năng
đăng ký, đăng nhập, quên mật khẩu.

## Quyết định
- Keycloak là nhà cung cấp định danh duy nhất; realm `crucible`; vai trò `STUDENT`, `INSTRUCTOR`, `ADMIN`.
- Giao diện đăng nhập bằng Authorization Code + PKCE; dịch vụ kiểm tra JWT bằng JWKS, kiểm tra `aud` và `iss`.
- Mỗi dịch vụ có một client với client credentials; quyền gọi dịch vụ khác biểu diễn bằng scope (CR-DOC-03 mục 3.2).
- Phân quyền theo đối tượng (giảng viên của lớp nào, học viên nào là chủ lượt làm bài) do module sở hữu đối tượng quyết định.
- Đơn vị có nhà cung cấp định danh (ví dụ Microsoft Entra ID cho email trường) được liên kết vào Keycloak bằng identity
  brokering OIDC; tài khoản được gắn với sinh viên trong lớp theo email. Đơn vị không có thì dùng tài khoản cục bộ của Keycloak.
- Vai trò trong lớp (giảng viên phụ trách, trợ giảng) do M3 quản lý, không lưu trong Keycloak.
- Cấu hình realm quản lý dạng mã bằng keycloak-config-cli.

## Phương án đã cân nhắc
| Phương án | Lý do không chọn |
|---|---|
| Mỗi module tự xác thực | Lặp lại công sức, không thống nhất |
| Một dịch vụ xác thực tự viết | Tốn công sức không thuộc nghiệp vụ của module nào |
| mTLS giữa dịch vụ bằng service mesh | Chi phí tài nguyên và vận hành cao so với cụm; scope theo nghiệp vụ vẫn cần token |

## Hệ quả
- (+) Đăng nhập, quên mật khẩu, khóa tài khoản có sẵn.
- (+) Quyền giữa dịch vụ khai báo tập trung, kiểm tra được.
- (−) Thêm một thành phần chiếm khoảng 0,8 GB bộ nhớ.
- (−) Thành viên phải hiểu luồng OIDC để cấu hình dịch vụ.
