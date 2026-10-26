# Danh sách việc: [TÊN TÍNH NĂNG]

| Thuộc tính | Giá trị |
|---|---|
| Đặc tả | `specs/[###-ten-ngan]/spec.md` |
| Kế hoạch | `specs/[###-ten-ngan]/plan.md` |

<!--
  Định dạng mỗi việc (bắt buộc, để lệnh implement, converge đánh dấu):
    - [ ] T001 [P] [US1] Mô tả việc kèm đường dẫn tệp — Kiểm chứng: <cách kiểm chứng>
  [P]: làm song song được (khác tệp, không phụ thuộc việc chưa xong). [US1]: câu chuyện người dùng trong spec.md,
  chỉ ghi ở các giai đoạn theo câu chuyện. Mỗi việc là một commit. Việc viết kiểm thử là bắt buộc (hiến pháp mục III)
  và đứng trước việc hiện thực tương ứng.
-->

## Giai đoạn 1: Chuẩn bị

- [ ] T001 [Việc chuẩn bị chung, ví dụ migration, cấu hình] trong [đường dẫn] — Kiểm chứng: [cách kiểm chứng]

## Giai đoạn 2: Nền tảng dùng chung cho các câu chuyện

- [ ] T002 [Việc mà mọi câu chuyện cần] trong [đường dẫn] — Kiểm chứng: [cách kiểm chứng]

## Giai đoạn 3: US1 — [Tên câu chuyện] (P1)

- [ ] T003 [P] [US1] Kiểm thử [AC-*] trong [đường dẫn tệp kiểm thử] — Kiểm chứng: kiểm thử thất bại trước khi hiện thực
- [ ] T004 [US1] Hiện thực [nội dung] trong [đường dẫn] — Kiểm chứng: kiểm thử T003 đạt

## Giai đoạn 4: US2 — [Tên câu chuyện] (P2)

- [ ] T005 [P] [US2] Kiểm thử [AC-*] trong [đường dẫn] — Kiểm chứng: [cách kiểm chứng]
- [ ] T006 [US2] Hiện thực [nội dung] trong [đường dẫn] — Kiểm chứng: [cách kiểm chứng]

## Giai đoạn cuối: Hoàn thiện

- [ ] T007 Cập nhật tài liệu module và `README` của dịch vụ nếu thiết kế thay đổi — Kiểm chứng: đọc lại
- [ ] T008 Chạy kiểm tra cục bộ theo CONTRIBUTING mục 2.4 — Kiểm chứng: lint, kiểm thử, build đạt

## Phụ thuộc và thứ tự

- Giai đoạn 1, 2 xong trước các giai đoạn theo câu chuyện.
- Các câu chuyện làm theo thứ tự ưu tiên; câu chuyện sau không phá kiểm thử của câu chuyện trước.

## Kết quả

[Ghi sau khi hoàn thành: số kiểm thử đạt, kết quả kiểm chứng thủ công, sai khác so với kế hoạch]
