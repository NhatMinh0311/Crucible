# Image

| Thư mục | Image | Nội dung |
|---|---|---|
| `crucible-tools/` | `ghcr.io/crucible-project/crucible/tools:1.0.0` | BusyBox, curl, jq liên kết tĩnh trong `/opt/crucible/bin`; gắn chỉ đọc vào môi trường cho lệnh kiểm tra (ADR-0011) |
| `lab-ubuntu-nginx/` | `ghcr.io/crucible-project/crucible-labs/ubuntu-nginx:1.0.0` | Image nền bài lab: Ubuntu 24.04, nginx, công cụ quản trị; PID 1 là `tini` chạy `sleep infinity` |

Build trên máy phát triển (chế độ `docker` của `runtime-api` dùng đúng các tên này):

```bash
docker build -t ghcr.io/crucible-project/crucible/tools:1.0.0 images/crucible-tools
docker build -t ghcr.io/crucible-project/crucible-labs/ubuntu-nginx:1.0.0 images/lab-ubuntu-nginx
```

Phiên bản phần mềm bên ngoài và mã băm được ghim trong `Dockerfile`; nâng phiên bản bằng cách sửa đồng thời số phiên bản
và mã băm, rồi tăng thẻ của image.
