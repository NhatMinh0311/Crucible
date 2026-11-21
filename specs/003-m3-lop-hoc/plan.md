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

Hiện thực chức năng quản lý lớp học, tài khoản hạn mức lớp và nhập danh sách sinh viên bằng CSV cho module M3 (`learning-service`). Sử dụng Spring Boot 4, đặc tả `services/learning/api/openapi.yaml` sinh interface controller `ClassesApi`, Flyway migration tạo các bảng `classes` và `class_members`, phân tích CSV streaming với Apache Commons CSV và `BOMInputStream` (hỗ trợ CSV UTF-8 có BOM từ Excel), đồng bộ tài khoản hạn mức `class:<id>` bên M1 qua `putQuotaAccount`, tự động liên kết tài khoản qua `UserAutoEnrollFilter`, và cung cấp giao diện React hoàn chỉnh thuộc SCR-21 (Quản lý lớp học, thành viên, hạn mức, nhập CSV cho GV/TG) và SCR-11 (Lớp học cho SV).

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
├── api/
│   └── openapi.yaml                                    # Đặc tả OpenAPI nội bộ của M3 (/v1/classes/**)
├── build.gradle.kts                                    # Sinh ClassesApi, commons-csv, commons-io, WireMock
└── src/
    ├── main/
    │   ├── java/dev/crucible/learning/
    │   │   ├── api/                                    # Tầng Web / REST Controller (sinh từ openapi.yaml)
    │   │   │   ├── ClassController.java                # Hiện thực ClassesApi (/v1/classes/**)
    │   │   │   ├── ProblemException.java
    │   │   │   └── GlobalExceptionHandler.java         # Xử lý lỗi RFC 9457 kèm traceId
    │   │   ├── application/                            # Tầng nghiệp vụ Use Cases
    │   │   │   ├── ClassService.java                   # Quản lý lớp học, thành viên & đồng bộ putQuotaAccount
    │   │   │   ├── CsvImportService.java               # Xử lý phân tích CSV (BOMInputStream) và validate dòng
    │   │   │   └── UserAutoEnrollService.java          # Tự động gán user_id, locale, joined_at theo email token
    │   │   ├── domain/                                 # Tầng thực thể & luật miền
    │   │   │   ├── ClassEntity.java                    # Thực thể lớp học
    │   │   │   ├── ClassMemberEntity.java              # Thực thể thành viên lớp
    │   │   │   ├── ClassRole.java                      # INSTRUCTOR, TEACHING_ASSISTANT, STUDENT
    │   │   │   └── ClassStatus.java                    # ACTIVE, ARCHIVED
    │   │   └── infrastructure/
    │   │       ├── persistence/                        # JPA Repositories
    │   │       │   ├── ClassRepository.java
    │   │       │   └── ClassMemberRepository.java
    │   │       └── security/                           # Tiện ích trích xuất UserContext & bộ lọc tự gắn lớp
    │   │           ├── SecurityConfig.java
    │   │           ├── UserAutoEnrollFilter.java
    │   │           └── UserContext.java
    │   └── resources/
    │       └── db/migration/
    │           └── V2__classes_and_members.sql         # Migration tạo bảng classes và class_members
    └── test/
        └── java/dev/crucible/learning/
            ├── application/
            │   ├── ClassServiceTest.java
            │   ├── CsvImportServiceTest.java
            │   └── UserAutoEnrollServiceTest.java
            └── api/
                └── ClassControllerTest.java

web/src/
├── locales/
│   ├── vi/learning.json                                # Khóa dịch tiếng Việt
│   └── en/learning.json                                # Khóa dịch tiếng Anh
├── shell/api/
│   └── client.ts                                       # learningApi (openapi-fetch)
└── modules/learning/
    ├── routes.tsx                                      # Cấu hình định tuyến màn hình lớp học
    ├── api/                                            # API Client gọi learning-service qua learningApi
    │   └── classApi.ts
    ├── components/                                     # Thành phần giao diện lớp học (SCR-21)
    │   └── CsvImportModal.tsx                          # Kéo thả CSV và bảng báo lỗi dòng
    └── pages/
        ├── ClassListPage.tsx                           # SCR-21 (GV) & SCR-11 (SV)
        └── ClassDetailPage.tsx                         # SCR-21 / SCR-11 (Chi tiết lớp, thành viên & hạn mức)
```

---

## 5. Lựa chọn kỹ thuật

| Vấn đề | Lựa chọn | Lý do | Phương án khác đã xét |
|---|---|---|---|
| Thư viện đọc CSV | `org.apache.commons:commons-csv:1.12.0` | Thư viện chuẩn công nghiệp, ổn định, hỗ trợ đọc theo dòng (streaming), tự xử lý escape quotes và ngắt dòng trong ô dữ liệu | Tự viết `split(",")` (dễ lỗi khi có dấu phẩy trong họ tên); OpenCSV (nhiều phụ thuộc thừa) |
| Lưu trữ thành viên chưa đăng ký | Lưu trong `class_members` với `user_id = NULL` | Đơn giản, tự nhiên, dễ dàng tìm kiếm khi sinh viên đăng nhập lần đầu bằng câu lệnh `UPDATE ... WHERE email = :email AND user_id IS NULL` | Tạo bảng đệm riêng `pending_enrollments` (tăng độ phức tạp truy vấn) |
| Định dạng mã lỗi trả về | RFC 9457 (`application/problem+json`) thông qua `ProblemDetail` của Spring 6 | Chuẩn REST của toàn dự án Crucible (CR-DOC-03 mục 2.2); tích hợp sẵn với trình xử lý lỗi ở frontend | Trả về cấu trúc JSON tự chế |
| Đồng bộ tự động gán lớp | Bộ lọc `UserAutoEnrollFilter` chạy trước mọi yêu cầu đã xác thực | Bảo đảm sinh viên mở trực tiếp đường dẫn chi tiết lớp học (`/v1/classes/{classId}`) trước khi vào danh sách lớp vẫn được gắn `user_id`, `locale`, `joined_at` ngay lập tức | Chỉ kiểm tra ở `GET /v1/classes` (gây lỗi 403 khi mở thẳng link chi tiết lớp) |

---

## 6. Mô hình dữ liệu và Migration

### Migration: `V2__classes_and_members.sql`
```sql
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL,
    name VARCHAR(255) NOT NULL,
    semester VARCHAR(32) NOT NULL,
    lead_instructor_id UUID NOT NULL,
    quota_account_ref VARCHAR(80) NOT NULL,
    weekly_minutes INT NOT NULL DEFAULT 1200,
    max_concurrent_sessions INT NOT NULL DEFAULT 30,
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_classes_code UNIQUE (code),
    CONSTRAINT uk_classes_quota_account_ref UNIQUE (quota_account_ref)
);

CREATE TABLE IF NOT EXISTS class_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    user_id UUID,
    role VARCHAR(32) NOT NULL,
    student_id_number VARCHAR(32),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    locale VARCHAR(8) NOT NULL DEFAULT 'vi',
    joined_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_class_member_student_id UNIQUE (class_id, student_id_number)
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_class_member_email_lower ON class_members (class_id, LOWER(email));
CREATE INDEX IF NOT EXISTS idx_class_members_user ON class_members (user_id);
CREATE INDEX IF NOT EXISTS idx_class_members_unlinked_email ON class_members (LOWER(email)) WHERE user_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_classes_status ON classes (status);
```

---

## 7. Hợp đồng và phụ thuộc

| Hạng mục | Nội dung |
|---|---|
| API & Sự kiện dùng tới | Token OIDC từ Keycloak (claims `sub`, `email`, `name`, `preferred_username`, `locale`, `realm_access.roles`); `PUT /v1/quota-accounts/{quotaAccountRef}` (`putQuotaAccount`) của M1 |
| Thay đổi hợp đồng | Không có (API lớp học là API nội bộ giữa `web` và `learning-service` tại `services/learning/api/openapi.yaml`) |
| Phụ thuộc module khác | Gọi `putQuotaAccount` của M1 (`runtime-api`) khi tạo lớp, đổi hạn mức hoặc đổi giảng viên phụ trách chính (`BR-M3.26`) |

---

## 8. Kiểm thử và kiểm chứng

| Mức | Nội dung kiểm thử | Mã tiêu chí |
|---|---|---|
| **Đơn vị** | Phân tích cú pháp CSV (kể cả UTF-8 BOM từ Excel): kiểm tra phát hiện đúng số dòng lỗi khi thiếu cột, email sai định dạng, trùng lặp mã SV/email | `AC-L1.5`, `BR-M3.05` |
| **Đơn vị & Controller** | Phân quyền: kiểm tra các ca từ chối 403 (`STUDENT` tạo lớp, người ngoài lớp xem chi tiết/thành viên/xem trước CSV, trợ giảng thêm/xóa thành viên) | `L01`, `AC-R2.8` |
| **Tích hợp (WireMock)** | `ClassServiceTest`: Tạo lớp học và cập nhật lớp gọi `putQuotaAccount` sang M1, chuyển tiếp `422 QUOTA_CEILING_EXCEEDED` | `L01`, `L12`, `AC-R2.6`, `AC-R2.8`, `BR-M3.26` |
| **Tích hợp** | `UserAutoEnrollServiceTest`: Mô phỏng sinh viên đăng nhập lần đầu qua `UserAutoEnrollFilter`, kiểm chứng `user_id`, `locale`, `joined_at` được cập nhật tự động | `AC-L1.5`, `BR-M3.05` |
| **Giao diện** | Vitest & TypeScript: Kiểm tra khóa dịch song ngữ `vi`/`en`, AST không có văn bản thô, kiểu dữ liệu sinh từ `openapi.yaml` | `SCR-11`, `SCR-21` |

