---
version: alpha
name: Crucible-Blueprint
description: >-
  Giao diện "Bản vẽ kỹ thuật" của Crucible, nền tảng thực hành Linux: giấy kẻ lưới, mực xanh đậm, xanh bản vẽ làm màu
  chính, cam tín hiệu cho thời gian và cảnh báo. Góc vuông, đường kẻ 1 px thay cho bóng đổ, các mục đánh ký hiệu
  (A —, 01), nhãn chữ mono in hoa. Có chế độ sáng và tối; terminal luôn nền tối. Mọi văn bản có tiếng Việt và tiếng Anh.
colors:
  background: "#f3f3ee"
  background-dark: "#0a1522"
  foreground: "#0b1b2b"
  foreground-dark: "#e9f0f7"
  surface: "#fafaf6"
  surface-dark: "#0d1b2b"
  muted: "#e9ebe4"
  muted-dark: "#13243a"
  muted-foreground: "#52647a"
  muted-foreground-dark: "#8aa0b5"
  border: "#c4cfda"
  border-dark: "#1f3349"
  rule: "#0b1b2b"
  rule-dark: "#3a5470"
  primary: "#1d4ed8"
  primary-dark: "#7cb4ff"
  primary-foreground: "#f6f7f2"
  primary-foreground-dark: "#0a1522"
  accent: "#e1e8fa"
  accent-dark: "#122843"
  ring: "#1d4ed8"
  ring-dark: "#7cb4ff"
  signal: "#b93d0b"
  signal-dark: "#ff9f43"
  signal-soft: "#fbe8dc"
  signal-soft-dark: "#2a1c0d"
  signal-foreground: "#fff6ef"
  signal-foreground-dark: "#0a1522"
  success: "#11704f"
  success-dark: "#5fd3a6"
  success-soft: "#dbede4"
  success-soft-dark: "#0f2a24"
  danger: "#b42318"
  danger-dark: "#ff8a7a"
  danger-soft: "#f8e1de"
  danger-soft-dark: "#2e1618"
  info: "#1d4ed8"
  info-dark: "#7cb4ff"
  info-soft: "#e1e8fa"
  info-soft-dark: "#122843"
  terminal: "#06101a"
  terminal-foreground: "#d7e4f1"
typography:
  display-lg:
    fontFamily: "Space Grotesk, IBM Plex Sans, sans-serif"
    fontSize: 30px
    fontWeight: 700
    lineHeight: 1.12
    letterSpacing: -0.02em
  display-md:
    fontFamily: "Space Grotesk, IBM Plex Sans, sans-serif"
    fontSize: 24px
    fontWeight: 700
    lineHeight: 1.2
    letterSpacing: -0.01em
  heading:
    fontFamily: "Space Grotesk, IBM Plex Sans, sans-serif"
    fontSize: 18px
    fontWeight: 600
    lineHeight: 1.3
    letterSpacing: -0.01em
  body-md:
    fontFamily: "IBM Plex Sans, sans-serif"
    fontSize: 14px
    fontWeight: 400
    lineHeight: 1.5
  body-strong:
    fontFamily: "IBM Plex Sans, sans-serif"
    fontSize: 14px
    fontWeight: 500
    lineHeight: 1.45
  button:
    fontFamily: "Space Grotesk, IBM Plex Sans, sans-serif"
    fontSize: 14px
    fontWeight: 700
    lineHeight: 1
  label:
    fontFamily: "IBM Plex Mono, monospace"
    fontSize: 11px
    fontWeight: 500
    lineHeight: 1.4
    letterSpacing: 0.08em
    textTransform: uppercase
  code:
    fontFamily: "IBM Plex Mono, monospace"
    fontSize: 13px
    fontWeight: 400
    lineHeight: 1.75
rounded:
  none: 0px
spacing:
  xs: 4px
  sm: 8px
  md: 12px
  lg: 16px
  xl: 20px
  xxl: 28px
  grid: 24px
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.primary-foreground}"
    borderColor: "{colors.primary}"
    typography: "{typography.button}"
    rounded: "{rounded.none}"
    height: 40px
    padding: 0 16px
  button-outline:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.foreground}"
    borderColor: "{colors.rule}"
    typography: "{typography.body-strong}"
    rounded: "{rounded.none}"
    height: 40px
    padding: 0 16px
  card:
    backgroundColor: "{colors.surface}"
    borderColor: "{colors.rule}"
    rounded: "{rounded.none}"
    padding: 20px
  status-badge:
    typography: "{typography.label}"
    rounded: "{rounded.none}"
    height: 24px
    padding: 0 8px
  section-label:
    textColor: "{colors.primary}"
    typography: "{typography.label}"
  nav-item-active:
    backgroundColor: "{colors.accent}"
    textColor: "{colors.foreground}"
    borderLeft: 3px solid {colors.primary}
  alert-signal:
    backgroundColor: "{colors.surface}"
    borderColor: "{colors.signal}"
    iconCellColor: "{colors.signal}"
  terminal:
    backgroundColor: "{colors.terminal}"
    textColor: "{colors.terminal-foreground}"
    typography: "{typography.code}"
    borderColor: "{colors.rule}"
---

# DESIGN.md — Crucible

Tệp này mô tả ngôn ngữ thiết kế để người và trợ lý AI dựng màn hình mới nhất quán. Token ở phần đầu tệp trùng với
`web/src/shell/theme/theme.css` (kiểm thử `web/src/shell/theme/design-md.test.ts` báo lỗi khi hai nơi lệch nhau).
Quy định đầy đủ về giao diện tại `docs/11-thiet-ke-giao-dien.md` (CR-DOC-11); khi tệp này và CR-DOC-11 khác nhau,
CR-DOC-11 là căn cứ.

## Overview

Crucible là công cụ làm việc: sinh viên làm bài lab và bài thi trong nhiều giờ, giảng viên soạn bài và giám sát kỳ thi.
Giao diện lấy cảm hứng từ bản vẽ kỹ thuật: thông tin có cấu trúc, đánh ký hiệu rõ ràng, ít trang trí. Cá tính đến từ
lưới, đường kẻ, chữ mono in hoa và độ tương phản, không đến từ hiệu ứng.

- Mật độ: trung bình (bảng, danh sách bước, terminal cạnh hướng dẫn).
- Biến thiên bố cục: thấp; cấu trúc lặp lại để người dùng định hướng nhanh.
- Chuyển động: tối thiểu; không dùng chuyển động trang trí.
- Chế độ: sáng, tối, theo hệ điều hành (người dùng chọn). Terminal luôn nền tối.

## Colors

| Vai trò | Token | Sáng | Tối |
|---|---|---|---|
| Nền trang, có lưới | `background` | #f3f3ee | #0a1522 |
| Thanh bên, header, thẻ | `surface` | #fafaf6 | #0d1b2b |
| Chữ chính | `foreground` | #0b1b2b | #e9f0f7 |
| Chữ phụ, nhãn | `muted-foreground` | #52647a | #8aa0b5 |
| Khung, đường kẻ chính | `rule` | #0b1b2b | #3a5470 |
| Đường kẻ trong (giữa các hàng) | `border` | #c4cfda | #1f3349 |
| Màu chính: nút chính, mục đang chọn, nhãn mục, liên kết | `primary` | #1d4ed8 | #7cb4ff |
| Nền mục đang chọn, hover | `accent` | #e1e8fa | #122843 |
| Tín hiệu: thời gian còn lại, cảnh báo, lỗi API | `signal` | #b93d0b | #ff9f43 |
| Trạng thái Đạt / Chưa đạt / Đang kiểm tra | `success`, `danger`, `info` + `*-soft` | | |

- Chỉ một màu chính (xanh bản vẽ). Cam tín hiệu dành cho những gì cần chú ý; không dùng cam để trang trí.
- Trạng thái luôn gồm màu, biểu tượng và nhãn chữ; nhãn trạng thái là chữ màu trên nền `*-soft` cùng sắc.
- Mọi cặp chữ trên nền đạt tỉ lệ tương phản 4,5:1 ở cả hai chế độ (kiểm thử `contrast.test.ts`).

## Typography

| Vai trò | Phông | Dùng cho |
|---|---|---|
| Tiêu đề, nút chính | Space Grotesk 600–700, giãn chữ âm nhẹ | Tiêu đề trang, tiêu đề thẻ, nút chính |
| Nội dung | IBM Plex Sans 400–600 | Đoạn văn, hướng dẫn bài lab, nhãn trường |
| Nhãn, số liệu, mã | IBM Plex Mono 400–600 | Nhãn mục, tiêu đề cột, nhãn trạng thái (in hoa, giãn 0,08em), số đo, terminal |

- Cả ba phông có bộ ký tự tiếng Việt. Không thêm phông khác khi chưa kiểm tra bộ ký tự tiếng Việt.
- Chữ in hoa chỉ dùng cho nhãn ngắn bằng chữ mono; không in hoa đoạn văn hay tiêu đề dài.

## Layout

- Khung chung: thanh bên 224 px (dấu hiệu, mục "A — Điều hướng", các mục đánh số 01, 02…), header cao 64 px chia thành
  ô bằng đường kẻ dọc, vùng nội dung có nền lưới 24 px.
- Vùng nội dung đệm 28 px trên máy tính, 16 px trên điện thoại. Thẻ và bảng dùng `surface` để nổi trên lưới.
- Các phần của một màn hình được đánh ký hiệu bằng nhãn mục: "B — Hướng dẫn", "C — Terminal".
- Trang làm bài, phòng chờ thi, soạn bài lab, chạy thử cần màn hình từ 1024 px (`RequireDesktop`); màn hình khác dùng
  được từ 360 px, dưới 768 px thanh bên chuyển thành menu trượt.

## Elevation & Depth

Không có bóng đổ. Phân lớp bằng đường kẻ (`rule` cho khung ngoài, `border` cho đường kẻ trong) và bằng nền (`surface`
trên `background` có lưới). Menu thả xuống và lớp phủ cũng dùng viền `rule`, không bóng.

## Shapes

Góc vuông ở mọi thành phần: nút, thẻ, ô nhập, nhãn trạng thái, ảnh đại diện. Ngoại lệ duy nhất là chấm trạng thái kết
nối. Đường kẻ 1 px; mục đang chọn trong thanh bên có vạch trái 3 px màu `primary`.

## Components

- **Nút chính**: nền `primary`, chữ Space Grotesk đậm, cao 40 px; một nút chính trên mỗi vùng thao tác.
- **Nút phụ**: viền `rule`, nền `surface`; hover nền `accent`. Nút cạnh nhau dính liền, chung đường kẻ.
- **Thẻ**: viền `rule`, nền `surface`, đệm 20 px, tiêu đề Space Grotesk 18 px.
- **Bảng, danh sách bước**: hàng tiêu đề chữ mono in hoa; cột số thứ tự chữ mono (01, 02); đường kẻ `border` giữa các hàng.
- **Nhãn trạng thái**: chữ mono in hoa 11 px, biểu tượng Lucide, nền `*-soft`.
- **Lỗi API**: ô biểu tượng nền `signal` bên trái, thông điệp theo `code`, dòng "Mã tra cứu" chữ mono; không hiển thị
  `detail`.
- **Terminal**: nền `terminal`, chữ IBM Plex Mono 13 px; thanh tiêu đề chữ mono in hoa và trạng thái kết nối.
- Thành phần có sẵn trong `web/src/shell/ui`: `Button`, `Card`, `StatusBadge`, `SectionLabel`, `ProblemAlert`,
  `dropdown-menu`. Dùng lại thay vì tự viết.

## Do's and Don'ts

**Nên**
- Dùng token qua lớp Tailwind (`bg-surface`, `border-rule`, `text-primary`, `text-signal`, `label-mono`, `font-display`).
- Đánh ký hiệu các phần của màn hình bằng `SectionLabel`.
- Đặt văn bản qua khóa dịch, đủ tiếng Việt và tiếng Anh (CR-DOC-11 mục 4.1).
- Hiển thị giờ theo `Asia/Ho_Chi_Minh` bằng `useFormat()`.

**Không nên**
- Mã màu viết thẳng trong màn hình; `rounded-*`; `shadow-*`; dải màu chuyển (gradient).
- Thẻ có vạch màu bên trái để nhấn mạnh; emoji làm biểu tượng.
- Chuyển động trang trí, hiệu ứng phóng to khi hover.
- Phông, màu nhấn thứ hai ngoài `primary` và `signal`.
- Sáo ngữ quảng cáo trong văn bản giao diện.

## Accessibility

- Tương phản 4,5:1 cho chữ; trạng thái không chỉ biểu đạt bằng màu.
- Mọi thao tác dùng được bằng bàn phím; tiêu điểm hiển thị bằng vòng 2 px màu `ring`.
- Nút chỉ có biểu tượng phải có `aria-label` qua khóa dịch.
