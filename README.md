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
| Nền tảng | Đăng nhập, sự kiện, lưu trữ, quan sát, triển khai | Keycloak, RabbitMQ, PostgreSQL, SeaweedFS, Prometheus, Grafana, Argo CD, k3s |

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
| [CR-DOC-09](docs/09-giai-phap-tuong-tu.md) | Khảo sát giải pháp tương tự, khoảng trống và định vị |
| [CR-DOC-10](docs/10-phan-tich-nghiep-vu.md) | Phân tích nghiệp vụ: use case, đặc tả use case, biểu đồ hoạt động, mô hình miền |
| [CR-DOC-11](docs/11-thiet-ke-giao-dien.md) | Thiết kế giao diện: nguyên tắc trải nghiệm, hệ thống thiết kế, sơ đồ và danh mục màn hình, bố cục, trạng thái, thông báo |
| [CR-DOC-12](docs/12-ke-hoach-kiem-thu.md) | Kế hoạch kiểm thử: mức kiểm thử, bất biến hệ thống, truy vết, kịch bản đầu cuối |
| [`docs/adr/`](docs/adr/) | Hồ sơ quyết định kiến trúc ADR-0001 đến ADR-0016 |
| [`docs/modules/`](docs/modules/) | Thiết kế bên trong từng module; mẫu tại `_mau.md` |

## Cấu trúc kho mã

| Đường dẫn | Nội dung |
|---|---|
| `contracts/` | Hợp đồng giữa các module: `runtime-api.yaml`, `content-api.yaml`, `events.asyncapi.yaml`, `lab-spec/` |
| `services/runtime/`, `services/content/`, `services/learning/` | Mã nguồn dịch vụ của M1, M2, M3 |
| `web/` | Giao diện React: `src/shell/` (khung chung), `src/modules/<module>/` |
| `DESIGN.md` | Ngôn ngữ thiết kế giao diện (token, phông, thành phần, quy tắc) cho người và trợ lý AI |
| `labs/` | Bài lab mẫu và image bài lab |
| `infra/` | Terraform và cloud-init dựng máy ảo, k3s, gVisor |
| `deploy/` | Helm chart, cấu hình cụm chuẩn, ứng dụng Argo CD |
| `tests/` | Kiểm thử cách ly, tải, chaos, đầu cuối |
| `libs/java/crucible-platform` | Thư viện dùng chung cho dịch vụ Spring Boot: outbox, inbox, CloudEvents, xác thực Keycloak |
| `dev/` | Dịch vụ phụ thuộc khi phát triển (Docker Compose), realm Keycloak, script kiểm tra |
| `.devcontainer/` | Môi trường phát triển chuẩn hóa |
| `scripts/` | Script tiện ích (tạo backlog trên GitHub) |

Các thư mục mã nguồn được bổ sung theo lộ trình tại CR-DOC-06.

## Kiểm tra hợp đồng

Yêu cầu: Node.js 24, Python 3.12.

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

## Hướng phát triển

Các hạng mục sau nằm ngoài MVP. Thiết kế hiện tại đã chừa chỗ để bổ sung mà không thay đổi hợp đồng theo cách phá vỡ
tương thích.

- [ ] **Môi trường nhiều máy (R14, ADR-0015).** Một phiên gồm 2–3 máy, mỗi máy là một pod gVisor, các máy cùng phiên
  chỉ liên lạc được với nhau. Mở ra dạng bài mạng giữa nhiều máy: DNS server cho máy client, ứng dụng gọi cơ sở dữ liệu ở
  máy khác, reverse proxy, tường lửa giữa hai máy, SSH bằng khóa.
  - Hợp đồng đã có trường `machine` (đặc tả bài lab, `ExecRequest`, vé terminal, vé xem), hiện chỉ nhận `main`.
  - Cần bổ sung: khai báo `environment.machines`; operator tạo các pod của một phiên theo nguyên tắc đủ cả hoặc không
    tạo, chạy setup theo thứ tự máy; NetworkPolicy và tên DNS theo phiên; sức chứa, hạn mức và đặt trước tính theo tổng
    chỗ của các máy; kho phiên dựng sẵn theo cấu hình máy; terminal nhiều tab.
  - Không bao gồm bài tầng liên kết dữ liệu và router là máy riêng (VLAN, DHCP, định tuyến động).
- [ ] **Định tuyến nhiều mạng con trong một môi trường (thử nghiệm S7).** `setup` dựng topo mô phỏng như
  `client – r1 – r2 – server` bằng namespace mạng, `veth` và `bridge` bên trong một pod; học viên sửa tuyến tĩnh,
  gateway, bật chuyển tiếp gói, lọc gói ở router.
  - Giữ nguyên một pod cho mỗi phiên: sức chứa, kho phiên dựng sẵn, terminal và hợp đồng không đổi.
  - Topo nằm trong user namespace con nên container không cần `SYS_ADMIN` và lệnh kiểm tra vẫn được bảo vệ (ADR-0011).
  - Theo mã nguồn gVisor: tạo được `veth`, `bridge`, tuyến tĩnh, `ip_forward`; không có `ip rule` và nhiều bảng định
    tuyến, không có `nftables`; bảng `nat` của `iptables` chưa chắc chắn.
  - Cần bổ sung: image nền mạng (`iproute2`, `iptables`, `traceroute`, `tcpdump`), thư viện script dựng topo, bài mẫu.
- [ ] **Bài lab Docker (R15).** Chạy Docker Engine trong môi trường học viên, phụ thuộc kết quả thử nghiệm S6.
- [ ] **Bài lab có bước phụ thuộc (C12).** Bước sau chỉ mở khi bước trước đạt.
- [ ] **Dạng bài viết script.** Học viên viết script, hệ thống chạy script để chấm trong điều kiện kiểm soát.
- [ ] **Kỳ thi nhiều bài lab.**
- [ ] **Bản dịch vụ cho nhiều đơn vị đào tạo** trên một bản cài đặt, dựa trên mã tham chiếu mờ và tài khoản hạn mức
  (ADR-0013). Phân hai cấp quản trị: quản trị nền tảng (sức chứa cụm, đơn vị, gói dịch vụ, tham số áp dụng cho mọi đơn
  vị) và quản trị của từng đơn vị (tài khoản, hạn mức, tham số hệ thống trong giới hạn của nền tảng). Bản tự cài hiện
  tại có một vai trò `ADMIN`.

## Đóng góp

Quy trình tại CR-DOC-07; lệnh và thao tác chi tiết tại [`CONTRIBUTING.md`](CONTRIBUTING.md).

- Mọi thay đổi vào `main` đi qua pull request, được review theo CODEOWNERS và qua CI.
- Thay đổi hợp đồng được thực hiện trong pull request riêng, merge trước pull request hiện thực.
- Thay đổi lược đồ CSDL bằng migration mới; migration đã merge không được sửa.

## Giấy phép

[GNU Affero General Public License v3.0](LICENSE).
