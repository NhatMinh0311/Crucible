# Giao diện Crucible

| Thư mục | Chủ sở hữu | Nội dung |
|---|---|---|
| `src/shell/` | Trưởng nhóm | Đăng nhập, điều hướng, theme, thành phần dùng chung (`ui/`), client API sinh từ `contracts/` |
| `src/modules/runtime/` | M1 | Terminal, quản trị cụm |
| `src/modules/content/` | M2 | Thư viện, soạn bài lab, trang làm bài |
| `src/modules/learning/` | M3 | Lớp học, lộ trình, bài tập, kỳ thi, giám sát, bảng điểm |

Mỗi module khai báo menu (`<module>Nav`) và route (`<module>Routes`) trong `routes.tsx` của mình; khung chung tự gộp.

```bash
npm install
npm run dev        # http://localhost:5173, cần dịch vụ phụ thuộc của dev/compose.yaml (Keycloak)
npm run typecheck
npm run lint
npm test
npm run build
```

Đăng nhập bằng tài khoản thử tại CONTRIBUTING mục 5.1. Yêu cầu tới `/api/<module>` được chuyển tới dịch vụ chạy trên máy
(8081, 8082, 8083). Hướng dẫn thiết kế giao diện tại CR-DOC-11; dùng skill `frontend-design` khi tạo màn hình mới.
