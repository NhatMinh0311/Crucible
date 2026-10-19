# Kế hoạch: khung giao diện chung và giao diện hai ngôn ngữ

Đặc tả: `spec.md`. Kế hoạch giữ cấu trúc hiện có của `web/` (React 19, React Router 8, Tailwind 4, `react-oidc-context`),
thêm lớp ngôn ngữ và hoàn thiện khung chung.

## 1. Cấu trúc mã

```text
web/
├── index.html                    script nhỏ đặt lớp `dark` trước khi React chạy (không nháy màu)
└── src/
    ├── locales/
    │   ├── vi/{shell,runtime,content,learning}.json
    │   ├── en/{shell,runtime,content,learning}.json
    │   └── locales.test.ts       cùng tập khóa; không có chữ viết thẳng trong JSX (W3)
    ├── shell/
    │   ├── i18n/
    │   │   ├── i18n.ts           khởi tạo i18next, tài nguyên nhúng sẵn
    │   │   ├── i18next.d.ts      khai báo kiểu: khóa sai là lỗi biên dịch
    │   │   ├── language.ts       chọn ngôn ngữ lần đầu, lưu trên máy, đồng bộ với claim `locale`, ghi Account API
    │   │   ├── format.ts         formatDateTime, formatDate, formatNumber (Asia/Ho_Chi_Minh)
    │   │   └── problem.ts        ApiProblem, thông điệp theo `code`, mã tra cứu
    │   ├── theme/theme.ts        Sáng, Tối, Theo hệ thống
    │   ├── layout/
    │   │   ├── AppShell.tsx      khung, menu, lớp phủ trên điện thoại, bỏ qua tới nội dung
    │   │   └── AccountMenu.tsx   tài khoản, ngôn ngữ, chế độ màu, đăng xuất
    │   ├── auth/RequireDesktop.tsx
    │   ├── pages/ErrorPage.tsx   errorElement của router, trang không tồn tại
    │   └── ui/dropdown-menu.tsx  thành phần shadcn/ui trên Radix
    └── modules/<module>/routes.tsx  nhãn menu chuyển sang khóa dịch của không gian module
dev/keycloak/realm-crucible.json  quyền account cho người dùng
```

## 2. Lựa chọn kỹ thuật

| Vấn đề | Lựa chọn | Lý do |
|---|---|---|
| Thư viện dịch | `i18next` 26.4.2, `react-i18next` 17.0.15 | Theo CR-DOC-11 mục 4.1 |
| Tải tệp dịch | Nhúng vào bundle | Bốn không gian nhỏ; không cần tải theo yêu cầu, không có khoảng trắng khi đổi ngôn ngữ |
| Kiểm tra khóa | Khai báo `CustomTypeOptions` của i18next theo tệp `vi`; kiểm thử Vitest so tập khóa | Khóa sai bị phát hiện khi biên dịch; thiếu khóa ở `en` bị phát hiện khi chạy kiểm thử trong CI |
| Chặn chữ viết thẳng | Kiểm thử Vitest duyệt cây cú pháp bằng TypeScript compiler API | oxlint 1.81 không có quy tắc tương đương `jsx-no-literals`; không thêm ESLint |
| Phát hiện ngôn ngữ | Hàm tự viết theo `navigator.languages` | Quy tắc của AC-P2.2 đơn giản; không cần thư viện phát hiện |
| Lưu `locale` | `GET` rồi `POST /realms/crucible/account` với `username`, `email`, `firstName`, `lastName`, `attributes.locale`; sau đó `signinSilent` để token mang claim mới | Account API của Keycloak; dịch vụ đọc claim `locale` (M3 gửi thông báo theo ngôn ngữ người nhận) |
| Trang đăng nhập | `signinRedirect({ ui_locales })` | Keycloak hiển thị theo `ui_locales` khi realm bật đa ngôn ngữ |
| Menu thả xuống | `@radix-ui/react-dropdown-menu` 2.1.24, mã thành phần theo shadcn/ui | Theo CR-DOC-11 mục 4; điều hướng bàn phím, ARIA có sẵn |
| Chế độ màu | Lớp `dark` trên `<html>`; `localStorage` khóa `crucible.theme`; `matchMedia('(prefers-color-scheme: dark)')` | Token màu hiện có dùng lớp `dark` |

Phiên bản thư viện là bản ổn định mới nhất đã phát hành ít nhất hai tuần.

## 3. Thứ tự chọn ngôn ngữ

```text
claim `locale` của token (vi | en)  →  lựa chọn lưu trên máy (`crucible.locale`)  →  navigator.languages[0]
```

Trước khi đăng nhập chưa có token, ngôn ngữ lấy từ máy hoặc trình duyệt và được truyền cho trang đăng nhập.

## 4. Không gian khóa

| Không gian | Nội dung | Chủ sở hữu |
|---|---|---|
| `shell` | Khung, menu tài khoản, trang lỗi, lỗi dùng chung (`VALIDATION_FAILED`, `UNAUTHENTICATED`, `FORBIDDEN`, `NOT_FOUND`, `METHOD_NOT_ALLOWED`, `RATE_LIMITED`, `IDEMPOTENCY_KEY_REUSED`, `SETTING_OUT_OF_RANGE`) | Trưởng nhóm |
| `runtime` | Menu, màn hình M1, lỗi riêng của M1 | Phụ trách M1 |
| `content` | Menu, màn hình M2, lỗi riêng của M2 | Phụ trách M2 |
| `learning` | Menu, màn hình M3, lỗi riêng của M3 (`CAPACITY_UNAVAILABLE` khi M3 trả) | Phụ trách M3 |

Lỗi của module tra theo khóa `errors.<CODE>` trong không gian của module.

## 5. Realm phát triển

Vai trò mặc định `default-roles-crucible` và các tài khoản thử có thêm `account: manage-account, view-profile`. Realm
chỉ được nhập khi chưa tồn tại; thành viên đã chạy Keycloak trước đó nhập lại realm theo hướng dẫn bổ sung tại
CONTRIBUTING mục 5.1.

## 6. Kiểm thử

| Mức | Nội dung |
|---|---|
| Đơn vị (Vitest) | Chọn ngôn ngữ (W1, trường hợp biên), định dạng ngày giờ với múi giờ máy khác (W5), thông điệp lỗi (W4), tập khóa và chữ viết thẳng (W3) |
| Thành phần (Testing Library) | Menu tài khoản đổi ngôn ngữ, chế độ màu; menu theo vai trò; lớp phủ và Esc; `RequireDesktop` |
| Trình duyệt thật | Puppeteer với Chrome và Keycloak phát triển: W1 (trang đăng nhập theo ngôn ngữ), W2 (đổi ngôn ngữ rồi đăng nhập trên ngữ cảnh trình duyệt mới), W7 |
