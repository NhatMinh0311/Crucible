# Giao diện Crucible

| Thư mục | Chủ sở hữu | Nội dung |
|---|---|---|
| `src/shell/` | Trưởng nhóm | Đăng nhập, điều hướng, menu tài khoản, theme, ngôn ngữ, thành phần dùng chung (`ui/`), client API sinh từ `contracts/` |
| `src/locales/<ngôn ngữ>/shell.json` | Trưởng nhóm | Văn bản của khung chung, lỗi API dùng chung |
| `src/modules/runtime/`, `src/locales/*/runtime.json` | M1 | Terminal, quản trị cụm |
| `src/modules/content/`, `src/locales/*/content.json` | M2 | Thư viện, soạn bài lab, trang làm bài |
| `src/modules/learning/`, `src/locales/*/learning.json` | M3 | Lớp học, lộ trình, bài tập, kỳ thi, giám sát, bảng điểm, thông báo |

Mỗi module khai báo menu (`<module>Nav`) và route (`<module>Routes`) trong `routes.tsx` của mình; khung chung tự gộp.
Route chỉ dành cho một số vai trò khai báo `handle: { roles: ['INSTRUCTOR'] } satisfies RouteHandle`; khung chung hiển thị
trang "không có quyền" khi người dùng gõ thẳng đường dẫn. Quyền thật vẫn do dịch vụ kiểm tra.

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

## Văn bản hai ngôn ngữ

Mọi văn bản hiển thị đi qua khóa dịch (CR-DOC-11 mục 4.1, AC-P2.1). Thêm khóa vào cả `src/locales/vi/<không gian>.json`
và `src/locales/en/<không gian>.json` của module mình:

```tsx
import { useTranslation } from 'react-i18next'

export function LabLibrary() {
  const { t } = useTranslation('content')        // không gian của module
  return <h1>{t('library.title')}</h1>             // khóa sai là lỗi kiểu khi biên dịch
}
```

- Nhãn menu dùng khóa có tiền tố không gian: `labelKey: 'content:nav.labs'`.
- Văn bản trong JSX và các thuộc tính `aria-label`, `title`, `placeholder`, `alt` không được viết thẳng; `npm test` báo
  lỗi kèm tên tệp và dòng. Hai tệp `vi`, `en` phải có cùng tập khóa và cùng biến `{{...}}`.
- Nội dung bài lab hiển thị nguyên văn, không qua khóa dịch; vùng hiển thị đặt `lang` theo `metadata.language`.

## Ngày giờ, số

```tsx
import { useFormat } from '@/shell/i18n/format'

const f = useFormat()
f.dateTime(attempt.startedAt)   // "15:30 19 thg 10, 2026" | "19 Oct 2026, 15:30", luôn theo giờ Việt Nam
f.number(score)
```

## Lỗi API

Lỗi RFC 9457 hiển thị theo `code`, kèm mã tra cứu; không hiển thị `detail`.

```tsx
import { ProblemAlert } from '@/shell/ui/problem-alert'

const { error } = await contentApi.POST('/v1/attempts', { body })
if (error) return <ProblemAlert problem={error} ns="content" />   // không có phản hồi: problem={undefined}
```

Thông điệp tra khóa `errors.<CODE>` trong không gian của module trước, sau đó tới `shell` (các mã dùng chung). Mỗi
module thêm bản dịch cho mã lỗi riêng của mình tại CR-DOC-03 mục 3.4, kèm việc người dùng có thể làm tiếp (UX5).

## Thành phần của khung

| Thành phần | Dùng khi |
|---|---|
| `RequireRole` | Một phần của màn hình chỉ dành cho một số vai trò (cả route: `handle.roles`) |
| `RequireDesktop` | Trang làm bài, phòng chờ thi, soạn bài lab, chạy thử (từ 1024 px, CR-DOC-11 mục 3) |
| `ComingSoon` | Chỗ giữ cho màn hình chưa xây dựng |
| `StatusBadge` | Trạng thái kèm màu, biểu tượng, nhãn |
| `ui/dropdown-menu` | Menu thả xuống theo shadcn/ui |

Chế độ màu (Sáng, Tối, Theo hệ thống) và ngôn ngữ do người dùng chọn ở menu tài khoản.

## Hướng giao diện "Bản vẽ kỹ thuật"

Mô tả đầy đủ cho người và trợ lý AI tại `DESIGN.md` ở gốc kho mã. Đổi token thì sửa cả `theme.css` và `DESIGN.md`;
`design-md.test.ts` báo lỗi khi hai tệp lệch nhau.

| Quy tắc | Cách viết |
|---|---|
| Màu | Chỉ dùng token: `bg-surface`, `border-rule` (khung), `border-border` (kẻ trong), `text-primary`, `text-signal` (thời gian, cảnh báo), `bg-success-soft text-success`…; không đặt mã màu trong màn hình |
| Hình khối | Góc vuông, viền 1 px; không `rounded-*`, không `shadow-*` |
| Chữ | Tiêu đề `font-display font-bold`; nhãn mục, tiêu đề cột `label-mono`; số liệu, mã `font-mono` |
| Nhãn mục | `<SectionLabel index="B">{t('...')}</SectionLabel>` |
| Nền vùng nội dung | Khung chung đã đặt `bg-blueprint`; thẻ dùng `bg-surface` để nổi trên lưới |

`src/shell/theme/contrast.test.ts` kiểm tra độ tương phản của token; thêm token mới thì thêm cặp chữ trên nền vào danh sách.
