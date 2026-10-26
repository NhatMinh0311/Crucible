# Hướng dẫn đóng góp

Tài liệu hướng dẫn các thao tác của thành viên khi đưa thay đổi vào kho mã. Quy trình, quyền sở hữu mã và cấu hình
bảo vệ nhánh được quy định tại CR-DOC-07 (`docs/07-quy-trinh-phat-trien.md`).

Các lệnh dùng shell bash. Trên Windows, chạy trong WSL2 và đặt mã nguồn trong hệ thống tệp của WSL2
(ví dụ `~/src/crucible`), không đặt trong `/mnt/c/`.

## 1. Chuẩn bị (thực hiện một lần)

### 1.1. Công cụ

| Công cụ | Phiên bản | Dùng cho | Ai cần |
|---|---|---|---|
| Git | 2.40 trở lên | Quản lý mã nguồn | Mọi thành viên |
| Docker (Engine trên Linux; Docker Desktop hoặc OrbStack trên macOS; Docker Desktop với WSL2 trên Windows) | 24 trở lên | Dịch vụ phụ thuộc khi phát triển, Testcontainers | Mọi thành viên |
| Node.js | 24 LTS | Kiểm tra hợp đồng, mock API, giao diện | Mọi thành viên |
| Python | 3.12 | Kiểm tra đặc tả bài lab | Mọi thành viên |
| JDK | 25 | `content-service`, `learning-service` | Phụ trách M2, M3 |
| Go | 1.27 | `runtime-api` | Phụ trách M1 |
| Tailscale | Bản hiện hành | Truy cập môi trường chuẩn (CR-DOC-04 mục 4) | Mọi thành viên |

### 1.1a. Dev Container (khuyến nghị)

Dev Container cài sẵn JDK 25, Node.js 24, Go 1.27, Python 3.12 đúng phiên bản và khởi động dịch vụ phụ thuộc
(mục 5.1) cùng lúc; máy chỉ cần Git, Docker và VS Code (tiện ích Dev Containers) hoặc IntelliJ IDEA.

1. Tải kho mã (mục 1.4). Trên Windows, tải vào hệ thống tệp của WSL2.
2. Mở thư mục kho mã bằng VS Code, chọn **Reopen in Container**. Lần đầu mất vài phút để tải image và công cụ.
3. Trong Dev Container, dịch vụ phụ thuộc chạy tại `localhost` như mục 5.1; VS Code tự chuyển tiếp các cổng ra máy.

Không chạy đồng thời dịch vụ phụ thuộc bằng `docker compose` trên máy và Dev Container: hai cách dùng chung tên dự án
Compose và volume dữ liệu.

### 1.2. Cấu hình Git

Email trong cấu hình Git phải trùng một email đã khai báo trong tài khoản GitHub (Settings → Emails), để commit được
ghi nhận cho đúng người.

```bash
git config --global user.name  "<Tên hiển thị>"
git config --global user.email "<email-tài-khoản-github>"
git config --global init.defaultBranch main
git config --global pull.ff only
```

### 1.3. Khóa SSH

```bash
ls ~/.ssh/id_ed25519.pub || ssh-keygen -t ed25519 -C "<email-tài-khoản-github>"
cat ~/.ssh/id_ed25519.pub        # dán vào GitHub → Settings → SSH and GPG keys → New SSH key
ssh -T git@github.com            # kết quả mong đợi: "Hi <username>! You've successfully authenticated"
```

### 1.4. Nhận quyền truy cập và tải kho mã

1. Chấp nhận lời mời collaborator của kho mã trên GitHub.
2. Chấp nhận lời mời tham gia tailnet của nhóm (Tailscale).
3. Tải kho mã và kiểm tra hợp đồng:

```bash
git clone git@github.com:NhatMinh0311/crucible.git
cd crucible
pip install --user jsonschema==4.23.0 pyyaml==6.0.2
python3 contracts/lab-spec/validate.py contracts/lab-spec/examples/*.yaml   # kết quả mong đợi: các tệp "Hợp lệ"
```

Hướng dẫn dựng môi trường phát triển của từng module được bổ sung vào mục 5 khi mã nguồn module có mặt.

### 1.5. Máy Windows

Phát triển M2, M3 và giao diện trên Windows được hỗ trợ đầy đủ khi tuân theo các điểm sau:

| Vấn đề | Cách làm |
|---|---|
| Script shell, Makefile, lệnh trong tài liệu dùng bash | Chạy trong WSL2 (Ubuntu 24.04); không dùng PowerShell hay Git Bash cho các lệnh của dự án |
| Truy cập tệp giữa Windows và Linux chậm; container đọc mã nguồn trên ổ Windows rất chậm | Đặt kho mã trong hệ thống tệp của WSL2 (`~/src/crucible`), không đặt trong `/mnt/c/` |
| IDE | VS Code với tiện ích WSL; hoặc IntelliJ IDEA mở dự án trong WSL (JDK và Gradle cài trong WSL) |
| Docker | Docker Desktop bật tích hợp WSL2, hoặc Docker Engine cài trực tiếp trong WSL2; cần bật ảo hóa trong BIOS. Testcontainers dùng được với cả hai |
| Bộ nhớ: Docker Desktop và WSL2 chạy trong máy ảo dùng chung bộ nhớ với Windows | Giới hạn bộ nhớ của WSL2 trong `%UserProfile%\.wslconfig` (ví dụ `memory=8GB` với máy 16 GB); chỉ chạy dịch vụ phụ thuộc mà module cần (mục 5) |
| Xuống dòng CRLF làm hỏng script | `.gitattributes` đã cố định LF; nếu gặp lỗi `\r`, xem mục 4 |
| Windows Defender làm chậm Gradle, `npm install` | Thêm thư mục dự án trong WSL2 và thư mục `~/.gradle` vào danh sách loại trừ |
| Đường dẫn dài trong `node_modules` | Không gặp khi kho mã nằm trong WSL2 |

Toàn bộ dịch vụ phụ thuộc (mục 5.1) dùng khoảng 0,9 GB bộ nhớ (đo ngày 17/10/2026, Keycloak chiếm khoảng 0,67 GB), nên máy
8 GB vẫn chạy được dịch vụ phụ thuộc cùng một module và IDE.

## 2. Thực hiện một thay đổi

### 2.1. Tạo nhánh từ `main` mới nhất

```bash
git switch main
git pull
git switch -c <loại>/<module>-<mô-tả-ngắn>   # ví dụ: feat/m2-variant-generator
```

Loại nhánh: `feat`, `fix`, `docs`, `test`, `chore`. Phần module: `m1`, `m2`, `m3`, `contracts`, `infra`, `docs`.
Mỗi nhánh giải quyết một việc; thay đổi lớn được chia thành nhiều pull request nhỏ.

### 2.2. Commit

```bash
git status                                 # kiểm tra tệp thay đổi; không có .env, tệp tạm, thư mục build
git add <tệp>
git commit
```

Quy ước nội dung commit:
- Dòng đầu: câu ngắn mô tả thay đổi, tối đa 72 ký tự, không có dấu chấm cuối. Ví dụ: `Sinh biến thể theo seed HMAC`.
- Để trống một dòng, sau đó giải thích lý do thay đổi và mã yêu cầu liên quan (ví dụ `US-C2`, `AC-C2.1`, `BR-M1.3`).

### 2.3. Cập nhật theo `main` trước khi push

```bash
git fetch origin
git rebase origin/main
```

Xử lý xung đột khi rebase:

```bash
git status                                 # liệt kê tệp xung đột
# sửa tệp, xóa các dấu <<<<<<< ======= >>>>>>>
git add <tệp-đã-sửa>
git rebase --continue                      # lặp lại cho tới khi hoàn tất
git rebase --abort                         # hủy rebase, quay về trạng thái trước đó
```

### 2.4. Kiểm tra cục bộ trước khi push

```bash
# Khi sửa hợp đồng
npx --yes @redocly/cli@1.34.5 lint contracts/runtime-api.yaml contracts/content-api.yaml
npx --yes @asyncapi/cli@2.16.0 validate contracts/events.asyncapi.yaml --fail-severity=error

# Khi sửa lược đồ đặc tả bài lab hoặc bài lab
python3 contracts/lab-spec/validate.py contracts/lab-spec/examples/*.yaml labs/*/lab.yaml

# Khi sửa services/runtime (cần Docker; image ở images/ phải được build trước, mục 5.2)
cd services/runtime && test -z "$(gofmt -l .)" && go vet ./... \
  && go run honnef.co/go/tools/cmd/staticcheck@v0.8.2 ./... && go test -race ./...
```

Lệnh kiểm tra của M2, M3 được bổ sung vào mục này khi mã nguồn module có mặt.

### 2.5. Push và mở pull request

```bash
git push -u origin <tên-nhánh>
```

Trên GitHub, bấm **Compare & pull request**, kiểm tra **base: `main`**, điền mô tả theo mẫu, bấm **Create pull request**.
Pull request chưa sẵn sàng để review được mở ở trạng thái **Draft**.

### 2.6. Theo dõi CI và review

- CI chạy tự động. Kiểm tra thất bại: bấm **Details**, đọc log, sửa, commit, push lại.
- GitHub tự yêu cầu review từ chủ sở hữu theo `.github/CODEOWNERS` (CR-DOC-07 mục 4).
- Mỗi góp ý được sửa bằng commit mới hoặc trả lời, sau đó bấm **Resolve conversation**.
- Commit mới sau khi đã được duyệt làm mất lượt duyệt; pull request cần được duyệt lại.
- Sau khi rebase, push lại bằng lệnh sau; lệnh từ chối ghi đè nếu trên GitHub có commit chưa tải về:

```bash
git push --force-with-lease
```

Nhánh `main` không nhận push trực tiếp và không nhận force push.

### 2.7. Merge và dọn dẹp

Người mở pull request bấm **Squash and merge** khi đủ điều kiện: lượt duyệt của chủ sở hữu, CI thành công, mọi trao đổi
đã được giải quyết. Nhánh trên GitHub được xóa khi merge.

```bash
git switch main
git pull
git fetch --prune
git branch -D <tên-nhánh>                  # squash merge tạo commit mới nên Git không nhận ra nhánh đã merge
```

## 3. Quy tắc theo loại thay đổi

### 3.1. Hợp đồng

- Thay đổi hợp đồng nằm trong pull request riêng sửa `contracts/`, merge **trước** pull request hiện thực.
- Mô tả pull request ghi rõ thay đổi tương thích hay không tương thích (CR-DOC-03 mục 2.2) và đánh dấu chủ các module
  sử dụng; pull request chỉ merge khi họ đã duyệt.

### 3.2. Migration cơ sở dữ liệu

- Mỗi module quản lý migration trong thư mục của mình (`services/<module>/.../migration/`).
- Migration đã merge không được sửa; sửa lỗi bằng migration mới.
- Migration theo nguyên tắc mở rộng rồi thu hẹp: không xóa hoặc đổi tên cột trong cùng phiên bản với thay đổi mã
  ngừng sử dụng cột đó.

### 3.3. Bài lab

- Mỗi bài lab là một thư mục `labs/<slug>/` gồm `lab.yaml` và tệp phụ trợ.
- Bài lab phải qua kiểm tra lược đồ và quy tắc ngữ nghĩa; lời giải mẫu phải đạt trên phiên chạy thử (CR-DOC-03 mục 4.5).

### 3.4. Hạ tầng và triển khai

- Thay đổi `infra/`, `deploy/` được áp dụng lên môi trường chuẩn bằng script triển khai sau khi merge (CR-DOC-04);
  không sửa trực tiếp trên máy chủ.
- Bí mật không được commit; giá trị nằm trong tệp `.env` trên máy chủ, kho mã chỉ chứa tệp mẫu.

### 3.5. Tài liệu

- Thay đổi hành vi hệ thống được cập nhật vào tài liệu tương ứng trong cùng pull request.
- Tài liệu dùng văn phong chính thức của dự án; câu hỏi và điểm chưa chắc được trao đổi trong pull request hoặc kênh
  của nhóm, không ghi vào tài liệu.
- Thành viên sử dụng công cụ AI ghi nhật ký tại CR-DOC-08 theo quy trình mục 2.1 của tài liệu đó.

## 4. Xử lý tình huống thường gặp

| Tình huống | Xử lý |
|---|---|
| `git push` bị từ chối do quy tắc bảo vệ `main` | Đã commit trực tiếp trên `main`. Xem dòng tiếp theo |
| Lỡ commit trên `main` cục bộ | `git switch -c <tên-nhánh>`, `git switch main`, `git reset --hard origin/main`, `git switch <tên-nhánh>`, push nhánh |
| Lỡ commit tệp `.env` hoặc mật khẩu | Chưa push: `git reset --soft HEAD~1`, `git restore --staged <tệp>`, commit lại. Đã push: báo trưởng nhóm để đổi bí mật bị lộ |
| Kiểm tra hợp đồng thất bại trên CI nhưng chạy được ở máy | Kiểm tra đúng phiên bản công cụ tại mục 2.4 |
| Script báo lỗi `\r` hoặc `bad interpreter` | Tệp bị đổi xuống dòng sang CRLF. Chạy `git add --renormalize . && git status`, commit kết quả |

## 5. Môi trường phát triển

### 5.1. Dịch vụ phụ thuộc

Dịch vụ phụ thuộc chạy bằng Docker Compose tại `dev/`; dịch vụ của module chạy từ IDE và kết nối qua `localhost`.

```bash
cp dev/.env.example dev/.env
docker compose -f dev/compose.yaml up -d --wait
./dev/check.sh                             # kết quả mong đợi: "Môi trường phát triển sẵn sàng."
docker compose -f dev/compose.yaml down    # dừng; thêm -v để xóa toàn bộ dữ liệu
```

| Dịch vụ | Địa chỉ | Thông tin đăng nhập (chỉ dùng khi phát triển) |
|---|---|---|
| PostgreSQL 18 | `localhost:5432` | CSDL và tài khoản `runtime`, `content`, `learning`; mật khẩu trong `dev/.env` |
| RabbitMQ 4.3 | AMQP `localhost:5672`, giao diện `http://localhost:15672` | `RABBITMQ_USER`, `RABBITMQ_PASSWORD` trong `dev/.env` |
| SeaweedFS (S3) | `http://localhost:8333`, bucket `recordings`, `lab-assets` | Khóa `crucible-dev` / `crucible-dev-secret` |
| Keycloak 26 | `http://localhost:8180`, realm `crucible` | Quản trị: `admin` / `KEYCLOAK_ADMIN_PASSWORD`; người dùng bên dưới |
| Mailpit | SMTP `localhost:1025`, giao diện `http://localhost:8025` | |

Tài khoản trong realm `crucible` (mật khẩu `crucible`):

| Tài khoản | Vai trò |
|---|---|
| `admin` | `ADMIN` |
| `gv1`, `gv2` | `INSTRUCTOR` |
| `tg1` | `INSTRUCTOR` (được thêm vào lớp với vai trò trợ giảng khi thử nghiệp vụ) |
| `sv1` … `sv5` | `STUDENT` |

Client trong realm: `crucible-web` (giao diện, public, PKCE), `runtime-api`, `content-service`, `learning-service`
(client credentials, secret `<tên client>-dev-secret`, scope theo CR-DOC-03 mục 3.2). Cấu hình realm tại
`dev/keycloak/realm-crucible.json`; sửa realm bằng pull request sửa tệp này. Mỗi tài khoản có vai trò `manage-account`,
`view-profile` của client `account` để giao diện lưu ngôn ngữ vào thuộc tính `locale` qua Account API. Cấu hình User
Profile chỉ cho quản trị sửa `username`, `email`, `firstName`, `lastName`; người dùng chỉ tự đổi được `locale`.

Keycloak chỉ nhập realm khi realm chưa tồn tại. Sau khi tệp realm thay đổi, nhập lại (giữ dữ liệu các CSDL khác; mã
`sub` của tài khoản thử đổi, dữ liệu phát triển tham chiếu tài khoản cũ cần tạo lại):

```bash
docker compose -f dev/compose.yaml stop keycloak
docker compose -f dev/compose.yaml run --rm --no-deps keycloak \
  import --file /opt/keycloak/data/import/realm-crucible.json --override true
docker compose -f dev/compose.yaml up -d --wait keycloak
```

### 5.2. Theo module

#### `runtime-api` trên máy phát triển (M1)

M2, M3 dùng một `runtime-api` thật trên máy phát triển: cùng API và sự kiện với môi trường chuẩn, mỗi phiên là một
container Docker trên máy (CR-DOC-03 mục 2.3, `specs/001-runtime-docker`). Trên Windows, chạy mọi lệnh trong WSL2.

1. Build image bộ công cụ kiểm tra và image nền bài lab (lặp lại khi `images/` thay đổi):

   ```bash
   docker build -t ghcr.io/crucible-project/crucible/tools:1.0.0 images/crucible-tools
   docker build -t ghcr.io/crucible-project/crucible-labs/ubuntu-nginx:1.0.0 images/lab-ubuntu-nginx
   ```

2. Chạy `runtime-api`, chọn một trong hai cách:

   | Cách | Lệnh | Dùng khi |
   |---|---|---|
   | Container (không cần cài Go) | `docker compose -f dev/compose.yaml --profile runtime up -d --build --wait` | Phát triển M2, M3 |
   | Từ mã nguồn | `cd services/runtime && RUNTIME_DRIVER=docker go run ./cmd/runtime-api` (dịch vụ phụ thuộc đã chạy) | Phát triển M1 |

   Dịch vụ lắng nghe `http://localhost:8081`; `content-service`, `learning-service` đã trỏ sẵn tới địa chỉ này
   (`RUNTIME_API_URL`). Terminal ở `ws://localhost:8081/v1/terminal?ticket=<vé>`.

3. Gọi thử bằng token của Keycloak phát triển:

   ```bash
   KC=http://localhost:8180/realms/crucible/protocol/openid-connect/token
   SVC=$(curl -s $KC -d grant_type=client_credentials -d client_id=content-service \
         -d client_secret=content-service-dev-secret | jq -r .access_token)          # token của M2
   STU=$(curl -s $KC -d grant_type=password -d client_id=crucible-web -d username=sv1 \
         -d password=crucible | jq -r .access_token)                                 # token của sinh viên
   OWNER=$(echo "$STU" | cut -d. -f2 | base64 -d 2>/dev/null | jq -r .sub)
   curl -s -X POST localhost:8081/v1/sessions -H "Authorization: Bearer $SVC" -H "Idempotency-Key: $(uuidgen)" \
     -H 'Content-Type: application/json' -d '{"ownerUserId":"'$OWNER'","priority":"PRACTICE",
     "quotaAccountRef":"class:1","environment":{"image":"ghcr.io/crucible-project/crucible-labs/ubuntu-nginx:1.0.0",
     "profile":"small","ttlMinutes":60,"idleMinutes":15,"maxExtensionMinutes":30,"network":"none"}}'
   ```

Khác biệt so với môi trường chuẩn:

| Mục | Máy phát triển |
|---|---|
| Cách ly | Không có gVisor; container chạy bằng runc trên máy |
| `network: mirror` | Không có máy chủ gói nội bộ; container ra Internet qua mạng `bridge` |
| Capability | Mặc định không cấp; đặt `SESSION_ALLOWED_CAPABILITIES=NET_ADMIN` trong `dev/.env` để thử bài tường lửa |
| Đặt trước, hạn mức, tham số hệ thống, bản ghi phiên | Trả 501 `NOT_IMPLEMENTED`; dùng mock Prism (CR-DOC-03 mục 2.3) |
| Hết bộ nhớ | Tiến trình vượt giới hạn bộ nhớ làm phiên kết thúc với `OUT_OF_MEMORY`, giống hành vi của gVisor trên môi trường chuẩn |

Container phiên mang nhãn `crucible.dev/managed-by=runtime-api`; `runtime-api` tự xóa container mồ côi khi khởi động.
Dọn bằng tay: `docker rm -f $(docker ps -aq --filter label=crucible.dev/managed-by=runtime-api)`.

Hướng dẫn chạy M2, M3 được bổ sung khi mã nguồn module có mặt.
