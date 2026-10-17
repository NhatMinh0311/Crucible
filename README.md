# Crucible

Crucible là nền tảng thực hành quản trị hệ thống Linux cho các học phần tại trường đại học. Mỗi học viên được cấp một
môi trường Linux riêng, cách ly, thao tác qua terminal trên trình duyệt và được chấm điểm tự động dựa trên trạng thái
của hệ thống. Giảng viên soạn bài lab có tham số để mỗi học viên nhận một biến thể riêng, tổ chức lớp, lộ trình và
kỳ thi thực hành đồng loạt.

Crucible là phần mềm mã nguồn mở theo giấy phép GNU AGPL-3.0; mỗi trường hoặc khoa tự cài đặt trên hạ tầng của mình.

## Kiến trúc tóm tắt

| Module | Trách nhiệm | Công nghệ |
|---|---|---|
| M1 `runtime` | Phiên thực hành, cấp phát tài nguyên, terminal, ghi phiên | Go, Kubernetes operator, gVisor |
| M2 `content` | Bài lab, biến thể, lượt làm bài, chấm điểm | Spring Boot |
| M3 `learning` | Lớp học, lộ trình, bài tập, kỳ thi, bảng điểm | Spring Boot |
| Nền tảng | Đăng nhập, sự kiện, lưu trữ, quan sát, triển khai | Keycloak, RabbitMQ, PostgreSQL, MinIO, Prometheus, Grafana, Argo CD, k3s |

Chi tiết tại CR-DOC-02.

## Tài liệu

| Tài liệu | Nội dung |
|---|---|
| [CR-DOC-01](docs/01-bai-toan-va-pham-vi.md) | Bài toán, đối tượng, mô hình cung cấp, chức năng, user story, tiêu chí chấp nhận, yêu cầu chất lượng |
| [CR-DOC-02](docs/02-kien-truc.md) | Kiến trúc: ngữ cảnh, container, ranh giới module, giao tiếp, vòng đời đối tượng, luồng chính |
| [CR-DOC-03](docs/03-hop-dong-giua-module.md) | Hợp đồng giữa các module: quy ước API, mã lỗi, đặc tả bài lab, sự kiện, giao thức terminal |
| [CR-DOC-04](docs/04-van-hanh.md) | Hạ tầng, triển khai, phát hành, giám sát, sao lưu, sức chứa, kiểm thử tải, chaos |
| [CR-DOC-05](docs/05-cach-ly-va-bao-mat.md) | Cách ly môi trường học viên, mối đe dọa, kiểm thử bảo mật |
| [CR-DOC-06](docs/06-rui-ro-va-lo-trinh.md) | Rủi ro, thử nghiệm sớm, lộ trình, phân công, tiêu chí hoàn thành |
| [CR-DOC-07](docs/07-quy-trinh-phat-trien.md) | Kho mã, quyền sở hữu, pull request, bảo vệ nhánh, CI |
| [CR-DOC-08](docs/08-nhat-ky-su-dung-ai.md) | Nhật ký sử dụng công cụ AI |
| [`docs/adr/`](docs/adr/) | Hồ sơ quyết định kiến trúc ADR-0001 đến ADR-0013 |
| [`docs/modules/`](docs/modules/) | Thiết kế bên trong từng module; mẫu tại `_mau.md` |

## Cấu trúc kho mã

| Đường dẫn | Nội dung |
|---|---|
| `contracts/` | Hợp đồng giữa các module: `runtime-api.yaml`, `content-api.yaml`, `events.asyncapi.yaml`, `lab-spec/` |
| `services/runtime/`, `services/content/`, `services/learning/` | Mã nguồn dịch vụ của M1, M2, M3 |
| `web/` | Giao diện React: `src/shell/` (khung chung), `src/modules/<module>/` |
| `labs/` | Bài lab mẫu và image bài lab |
| `infra/` | Terraform và cloud-init dựng máy ảo, k3s, gVisor |
| `deploy/` | Helm chart, cấu hình cụm chuẩn, ứng dụng Argo CD |
| `tests/` | Kiểm thử cách ly, tải, chaos, đầu cuối |

Các thư mục mã nguồn được bổ sung theo lộ trình tại CR-DOC-06.

## Kiểm tra hợp đồng

Yêu cầu: Node.js 22, Python 3.12.

```bash
npx --yes @redocly/cli@1.34.5 lint contracts/runtime-api.yaml contracts/content-api.yaml
npx --yes @asyncapi/cli@2.16.0 validate contracts/events.asyncapi.yaml --fail-severity=error
pip install jsonschema==4.23.0 pyyaml==6.0.2
python3 contracts/lab-spec/validate.py contracts/lab-spec/examples/*.yaml
```

Mock API của module khác khi phát triển:

```bash
npx @stoplight/prism-cli@5.14.2 mock contracts/runtime-api.yaml -p 4011
npx @stoplight/prism-cli@5.14.2 mock contracts/content-api.yaml -p 4012
```

## Đóng góp

Quy trình tại CR-DOC-07; lệnh và thao tác chi tiết tại [`CONTRIBUTING.md`](CONTRIBUTING.md).

- Mọi thay đổi vào `main` đi qua pull request, được review theo CODEOWNERS và qua CI.
- Thay đổi hợp đồng được thực hiện trong pull request riêng, merge trước pull request hiện thực.
- Thay đổi lược đồ CSDL bằng migration mới; migration đã merge không được sửa.

## Giấy phép

[GNU Affero General Public License v3.0](LICENSE).
