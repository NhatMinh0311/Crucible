# Danh sách việc: Quản lý lớp học và nhập danh sách sinh viên bằng CSV

| Thuộc tính | Giá trị |
|---|---|
| Đặc tả | `specs/003-m3-lop-hoc/spec.md` |
| Kế hoạch | `specs/003-m3-lop-hoc/plan.md` |
| Issue | #40 |

<!--
  Định dạng mỗi việc (bắt buộc theo quy chuẩn Spec Kit):
    - [ ] T001 [P] [US1] Mô tả việc kèm đường dẫn tệp — Kiểm chứng: <cách kiểm chứng>
  Mỗi việc là một commit. Việc viết kiểm thử là bắt buộc và đứng trước việc hiện thực tương ứng.
-->

## Giai đoạn 1: Chuẩn bị CSDL và Thư viện

- [x] T001 Thêm phụ thuộc `commons-csv` vào `services/learning/build.gradle.kts` — Kiểm chứng: `./gradlew dependencies` biên dịch thành công
- [x] T002 Tạo migration Flyway `V2__classes_and_members.sql` tạo bảng `classes` và `class_members` trong `services/learning/src/main/resources/db/migration/` — Kiểm chứng: Flyway migration chạy thành công trên PostgreSQL

## Giai đoạn 2: Nền tảng thực thể và phân quyền

- [x] T003 [P] Tạo các thực thể miền và enum: `ClassEntity`, `ClassMemberEntity`, `ClassRole`, `ClassStatus`, `CsvRowError` trong `services/learning/src/main/java/dev/crucible/learning/domain/` — Kiểm chứng: biên dịch mã nguồn không lỗi
- [x] T004 [P] Tạo JPA repositories: `ClassRepository` và `ClassMemberRepository` trong `services/learning/src/main/java/dev/crucible/learning/infrastructure/persistence/` — Kiểm chứng: test repository đơn giản với DataJpaTest
- [x] T005 Tạo bộ tiện ích trích xuất `UserContext` từ JWT token OIDC của Keycloak trong `services/learning/src/main/java/dev/crucible/learning/infrastructure/security/` — Kiểm chứng: trích xuất đúng `userId`, `email`, danh sách vai trò từ mock JWT

## Giai đoạn 3: US1 — Giảng viên tạo và quản lý lớp học (P1)

- [x] T006 [P] [US1] Viết kiểm thử `ClassServiceTest` cho các ca tạo lớp, kiểm tra mã lớp duy nhất, phân bổ hạn mức giờ máy — Kiểm chứng: kiểm thử hoàn thành
- [x] T007 [US1] Hiện thực `ClassService.createClass` và `ClassService.getClassDetail` trong `services/learning/src/main/java/dev/crucible/learning/application/` — Kiểm chứng: kiểm thử T006 đạt
- [x] T008 [P] [US1] Viết DTOs và cấu trúc phản hồi lỗi theo chuẩn RFC 9457 — Kiểm chứng: biên dịch hoàn tất
- [x] T009 [US1] Hiện thực `ClassController` và GlobalExceptionHandler (RFC 9457) trong `services/learning/src/main/java/dev/crucible/learning/api/` — Kiểm chứng: endpoint xử lý đúng định dạng ProblemDetail

## Giai đoạn 4: US2 — Nhập danh sách sinh viên bằng CSV (P1; AC-L1.5, BR-M3.05)

- [x] T010 [P] [US2] Viết kiểm thử đơn vị `CsvImportServiceTest`: kiểm tra phát hiện lỗi theo dòng (thiếu cột, email sai cú pháp RFC 5322, trùng lặp mã SV) `AC_L1_5` — Kiểm chứng: kiểm thử hoàn thành
- [x] T011 [US2] Hiện thực `CsvImportService.parseAndValidate` và `CsvRowError` trong `services/learning/src/main/java/dev/crucible/learning/application/` — Kiểm chứng: kiểm thử T010 đạt
- [x] T012 [P] [US2] Thiết kế cấu trúc phản hồi `CsvImportPreviewResponse` (chế độ preview và commit) — Kiểm chứng: DTO biên dịch hợp lệ
- [x] T013 [US2] Hiện thực controller endpoint cho `POST /api/v1/classes/{classId}/members/import-csv` — Kiểm chứng: upload và xử lý CSV hoàn tất

## Giai đoạn 5: US3 — Phân quyền trợ giảng và giảng viên (P2; L01)

- [x] T014 [P] [US3] Viết kiểm thử cho chức năng thêm giảng viên / trợ giảng vào lớp và kiểm tra quyền quản trị — Kiểm chứng: kiểm thử hoàn thành
- [x] T015 [US3] Hiện thực `ClassService.addInstructorOrAssistant` và API `POST /api/v1/classes/{classId}/instructors` — Kiểm chứng: phân quyền thành công

## Giai đoạn 6: US4 — Tự động liên kết tài khoản sinh viên khi đăng nhập (P1; AC-L1.5)

- [x] T016 [P] [US4] Viết kiểm thử tích hợp `UserAutoEnrollServiceTest`: tự động cập nhật `user_id` từ `sub` khi email trùng khớp `AC_L1_5` — Kiểm chứng: kiểm thử hoàn thành
- [x] T017 [US4] Hiện thực `UserAutoEnrollService` tích hợp vào luồng `GET /api/v1/classes` — Kiểm chứng: tự động liên kết thành công

## Giai đoạn 7: Giao diện người dùng M3 (SCR-11, SCR-12, SCR-13, SCR-21)

- [x] T018 [P] Bổ sung khóa dịch song ngữ cho lớp học và nhập CSV trong `web/src/locales/vi/learning.json` và `en/learning.json` — Kiểm chứng: khớp 100% khóa dịch 2 ngôn ngữ
- [x] T019 [P] Xây dựng API client `web/src/modules/learning/api/classApi.ts` — Kiểm chứng: kết nối với learning-service
- [x] T020 Xây dựng màn hình danh sách lớp học: SCR-11 (Giảng viên) và SCR-21 (Sinh viên) trong `web/src/modules/learning/pages/ClassListPage.tsx` — Kiểm chứng: hiển thị danh sách lớp và modal tạo lớp
- [x] T021 Xây dựng màn hình chi tiết lớp và tab thành viên SCR-12 trong `web/src/modules/learning/pages/ClassDetailPage.tsx` — Kiểm chứng: hiển thị danh sách giảng viên, trợ giảng, sinh viên
- [x] T022 Xây dựng thành phần modal/trang kéo thả CSV và hiển thị bảng báo lỗi chi tiết theo dòng SCR-13 trong `web/src/modules/learning/components/CsvImportModal.tsx` — Kiểm chứng: xem trước và báo lỗi dòng đúng giao diện Bản vẽ kỹ thuật

## Giai đoạn cuối: Hoàn thiện

- [ ] T023 Cập nhật tài liệu thiết kế module M3 và README nếu có sai khác — Kiểm chứng: đọc lại
- [ ] T024 Chạy kiểm tra toàn bộ cục bộ (Gradle build, test, frontend lint & typecheck) — Kiểm chứng: toàn bộ test suite xanh
