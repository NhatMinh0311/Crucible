# Kế hoạch: Quản lý lớp học và nhập danh sách sinh viên bằng CSV

| Thuộc tính | Giá trị |
|---|---|
| Đặc tả | `specs/003-m3-lop-hoc/spec.md` |
| Module | M3 (Học tập và kỳ thi) |
| Mốc | MA |
| Ngày lập | 09/11/2026 |
| Issue | #40 |

---

## 1. Tóm tắt

Hiện thực chức năng quản lý lớp học và nhập danh sách sinh viên bằng CSV cho module M3 (`learning-service`). Sử dụng Spring Boot 4, Flyway migration tạo các bảng `classes` và `class_members`, phân tích CSV streaming với Apache Commons CSV hoặc FastCSV, xử lý tự động liên kết tài khoản sinh viên qua Keycloak JWT token, và cung cấp giao diện React hoàn chỉnh gồm các màn hình SCR-11 (Danh sách lớp giảng viên), SCR-12 (Chi tiết lớp - tab Thành viên), SCR-13 (Kéo thả và xem trước CSV), và SCR-21 (Danh sách lớp sinh viên).

---

## 2. Bối cảnh kỹ thuật

| Mục | Giá trị |
|---|---|
| Ngôn ngữ, phiên bản | Java 25 (hỗ trợ JDK 21 qua `-PjavaVersion=21`); Spring Boot 4.1.1; TypeScript 5.9, React 19.3 |
| Thư viện chính | Spring Data JPA, Hibernate, Flyway, Apache Commons CSV `1.12.0` (hoặc Jackson CSV), `crucible-platform` 0.1.0 (xác thực Keycloak JWT) |
| Lưu trữ | PostgreSQL 18: schema `learning`, các bảng `classes`, `class_members` |
| Kiểm thử | JUnit 5, Mockito, Spring Boot Test, Testcontainers PostgreSQL 18, Vitest (Giao diện) |
| Ràng buộc | Thời gian phân tích và kiểm tra tệp CSV 200 dòng ≤ 1.5 giây; mã lỗi chuẩn RFC 9457; văn bản giao diện đi qua khóa i18n (`src/locales/vi/learning.json` và `en/learning.json`) |

---

## 3. Kiểm tra theo hiến pháp

*Điều kiện phải đạt trước khi thiết kế và kiểm tra lại sau thiết kế (`.specify/memory/constitution.md`).*

| Nguyên tắc | Đạt | Ghi chú |
|---|---|---|
| I. Yêu cầu dẫn từ tài liệu đã duyệt; không còn điểm chưa làm rõ | Có | Bám sát `L01`, `L12`, `AC-L1.5`, `AC-R2.8`, `BR-M3.05`, `BR-M3.26` trong CR-DOC-01 và CR-MOD-M3 |
| II. Hợp đồng: không đổi, hoặc thay đổi được tách thành pull request riêng | Có | API phục vụ nội bộ giao diện M3, không thay đổi hợp đồng dùng chung `contracts/` |
| III. Mỗi `AC-*`, `BR-*` được hiện thực có kiểm thử mang mã | Có | Có kiểm thử riêng gắn mã `AC_L1_5`, `BR_M3_05`, `BR_M3_26` |
| IV. Chỉ sửa mã trong module; phụ thuộc module khác được ghi ở mục 7 | Có | Chỉ sửa mã trong `services/learning/` và `web/src/modules/learning/` |
| V. Bí mật, log, migration, metric theo quy định | Có | Sử dụng migration Flyway `V2__classes_and_members.sql`; log JSON có `trace_id` |

---

## 4. Cấu trúc mã

```text
services/learning/
├── build.gradle.kts                                    # Thêm thư viện commons-csv nếu cần
└── src/
    ├── main/
    │   ├── java/dev/crucible/learning/
    │   │   ├── api/                                    # Tầng Web / REST Controller
    │   │   │   ├── ClassController.java                # REST endpoints /api/v1/classes/**
    │   │   │   ├── dto/                                # Request / Response DTOs
    │   │   │   │   ├── CreateClassRequest.java
    │   │   │   │   ├── ClassResponse.java
    │   │   │   │   ├── ClassMemberResponse.java
    │   │   │   │   ├── CsvImportPreviewResponse.java
    │   │   │   │   └── AddInstructorRequest.java
    │   │   │   └── GlobalExceptionHandler.java         # Xử lý lỗi RFC 9457
    │   │   ├── application/                            # Tầng nghiệp vụ Use Cases
    │   │   │   ├── ClassService.java                   # Quản lý lớp học & thành viên
    │   │   │   ├── CsvImportService.java               # Xử lý phân tích CSV và validate dòng
    │   │   │   └── UserAutoEnrollService.java          # Tự động gán user_id theo email token
    │   │   ├── domain/                                 # Tầng thực thể & luật miền
    │   │   │   ├── ClassEntity.java                    # Thực thể lớp học
    │   │   │   ├── ClassMemberEntity.java              # Thực thể thành viên lớp
    │   │   │   ├── ClassRole.java                      # INSTRUCTOR, TEACHING_ASSISTANT, STUDENT
    │   │   │   ├── ClassStatus.java                    # ACTIVE, ARCHIVED
    │   │   │   └── CsvRowError.java                    # Đối tượng lỗi theo dòng CSV
    │   │   └── infrastructure/
    │   │       ├── persistence/                        # JPA Repositories
    │   │       │   ├── ClassRepository.java
    │   │       │   └── ClassMemberRepository.java
    │   │       └── security/                           # Tiện ích trích xuất UserContext từ JWT
    │   │           └── UserContext.java
    │   └── resources/
    │       └── db/migration/
    │           └── V2__classes_and_members.sql         # Migration tạo bảng classes và class_members
    └── test/
        └── java/dev/crucible/learning/
            ├── application/
            │   ├── ClassServiceTest.java
            │   └── CsvImportServiceTest.java
            └── api/
                └── ClassControllerIntegrationTest.java

web/src/
├── locales/
│   ├── vi/learning.json                                # Khóa dịch tiếng Việt
│   └── en/learning.json                                # Khóa dịch tiếng Anh
└── modules/learning/
    ├── routes.tsx                                      # Cấu hình định tuyến màn hình lớp học
    ├── api/                                            # API Client gọi learning-service
    │   └── classApi.ts
    ├── components/                                     # Thành phần giao diện lớp học
    │   ├── ClassCard.tsx
    │   ├── CreateClassDialog.tsx
    │   ├── MemberTable.tsx
    │   └── CsvImportModal.tsx                          # Kéo thả CSV và bảng báo lỗi dòng
    └── pages/
        ├── ClassListPage.tsx                           # SCR-11 (GV) & SCR-21 (SV)
        ├── ClassDetailPage.tsx                         # SCR-12 (Chi tiết lớp & tab Thành viên)
        └── CsvImportPage.tsx                           # SCR-13
```

---

## 5. Lựa chọn kỹ thuật

| Vấn đề | Lựa chọn | Lý do | Phương án khác đã xét |
|---|---|---|---|
| Thư viện đọc CSV | `org.apache.commons:commons-csv:1.12.0` | Thư viện chuẩn công nghiệp, ổn định, hỗ trợ đọc theo dòng (streaming), tự xử lý escape quotes và ngắt dòng trong ô dữ liệu | Tự viết `split(",")` (dễ lỗi khi có dấu phẩy trong họ tên); OpenCSV (nhiều phụ thuộc thừa) |
| Lưu trữ thành viên chưa đăng ký | Lưu trong `class_members` với `user_id = NULL` | Đơn giản, tự nhiên, dễ dàng tìm kiếm khi sinh viên đăng nhập lần đầu bằng câu lệnh `UPDATE ... WHERE email = :email AND user_id IS NULL` | Tạo bảng đệm riêng `pending_enrollments` (tăng độ phức tạp truy vấn) |
| Định dạng mã lỗi trả về | RFC 9457 (`application/problem+json`) thông qua `ProblemDetail` của Spring 6 | Chuẩn REST của toàn dự án Crucible (CR-DOC-03 mục 2.2); tích hợp sẵn với trình xử lý lỗi ở frontend | Trả về cấu trúc JSON tự chế |
| Đồng bộ tự động gán lớp | Bộ chặn / Hook kiểm tra danh sách lớp khi gọi API `/api/v1/classes` | Sinh viên đăng nhập lần đầu chắc chắn sẽ tải danh sách lớp học của mình; thực hiện đối soát và gán `user_id` ngay tại luồng này | Dùng RabbitMQ nhận sự kiện đăng nhập từ Keycloak (phụ thuộc vào plugin Event Listener của Keycloak) |

---

## 6. Mô hình dữ liệu và Migration

### Migration: `V2__classes_and_members.sql`
```sql
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL,
    name VARCHAR(255) NOT NULL,
    semester VARCHAR(32) NOT NULL,
    instructor_quota_account_ref VARCHAR(64),
    session_limit_per_student INT NOT NULL DEFAULT 1,
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_classes_code UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS class_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    user_id UUID,
    role VARCHAR(32) NOT NULL,
    student_id_number VARCHAR(32),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_class_member_email UNIQUE (class_id, email),
    CONSTRAINT uk_class_member_student_id UNIQUE (class_id, student_id_number)
);

CREATE INDEX idx_class_members_user_id ON class_members(user_id);
CREATE INDEX idx_class_members_email ON class_members(email);
CREATE INDEX idx_classes_status ON classes(status);
```

---

## 7. Hợp đồng và phụ thuộc

| Hạng mục | Nội dung |
|---|---|
| API & Sự kiện dùng tới | Token OIDC từ Keycloak (claims `sub`, `email`, `preferred_username`, `roles`) |
| Thay đổi hợp đồng | Không có (API lớp học là API nội bộ giữa `web` và `learning-service`) |
| Phụ thuộc module khác | Không phụ thuộc M1, M2 ở tính năng này; M3 hoạt động độc lập để quản trị lớp học |

---

## 8. Kiểm thử và kiểm chứng

| Mức | Nội dung kiểm thử | Mã tiêu chí |
|---|---|---|
| **Đơn vị** | Phân tích cú pháp CSV: kiểm tra phát hiện đúng số dòng lỗi khi thiếu cột, email sai định dạng RFC 5322, trùng lặp mã SV nội bộ | `AC-L1.5`, `BR-M3.05` |
| **Đơn vị** | Phân quyền: kiểm tra trợ giảng không có quyền xóa lớp hoặc đổi hạn mức | `L01`, `AC-R2.8` |
| **Tích hợp** | `ClassServiceTest`: Tạo lớp học, đảm bảo mã lớp là duy nhất và người tạo tự động là `INSTRUCTOR` | `L01` |
| **Tích hợp** | `CsvImportServiceTest`: Nhập file CSV vào database, kiểm tra không ghi đè dữ liệu trùng, chỉ ghi nhận dòng hợp lệ | `AC-L1.5`, `BR-M3.05` |
| **Tích hợp** | `UserAutoEnrollTest`: Mô phỏng sinh viên đăng nhập lần đầu, kiểm chứng `user_id` được cập nhật tự động từ email | `AC-L1.5` |
| **Giao diện** | Vitest: Kiểm tra hiển thị màn hình lớp học tiếng Việt/tiếng Anh, modal kéo thả CSV và bảng danh sách lỗi dòng | `SCR-11`, `SCR-13` |

