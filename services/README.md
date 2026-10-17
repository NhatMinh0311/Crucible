# Dịch vụ

| Thư mục | Module | Dịch vụ | Cổng khi phát triển | Thiết kế |
|---|---|---|---|---|
| `runtime/` | M1 | `runtime-api`, `runtime-gateway`, `runtime-operator` (Go) | 8081 | `docs/modules/m1-runtime.md` |
| `content/` | M2 | `content-service` (Spring Boot) | 8082 | `docs/modules/m2-content.md` |
| `learning/` | M3 | `learning-service` (Spring Boot) | 8083 | `docs/modules/m3-learning.md` |

Khung dịch vụ Spring Boot gồm sẵn:

| Thành phần | Vị trí | Ghi chú |
|---|---|---|
| Cấu trúc tầng | `api`, `application`, `domain`, `infrastructure` | Quy ước tại `AGENTS.md` mục 4 và `package-info.java` của từng tầng |
| Interface controller sinh từ hợp đồng (chỉ M2) | `build/generated/openapi/server` | Sinh từ `contracts/content-api.yaml` khi biên dịch; controller hiện thực interface này |
| Client gọi module khác | `build/generated/openapi/<module>-client` | Spring HTTP Interface; token client credentials gắn tự động |
| Xác thực | `infrastructure/security/SecurityConfig.java` | Token Keycloak; vai trò thành `ROLE_*`, scope thành `SCOPE_*` |
| Sự kiện | Thư viện `libs/java/crucible-platform` | `OutboxPublisher.publish(...)` trong transaction; `InboxGuard.accept(...)` khi nhận |
| Migration | `src/main/resources/db/migration` | Flyway; bảng outbox, inbox do thư viện tạo |
| Kiểm thử tích hợp | `src/test/.../EventsPlumbingTest.java` | Testcontainers PostgreSQL, RabbitMQ |

Lệnh thường dùng (trong thư mục dịch vụ):

```bash
./gradlew bootRun        # chạy với dịch vụ phụ thuộc của dev/compose.yaml
./gradlew test           # kiểm thử, cần Docker cho Testcontainers
./gradlew build          # biên dịch, sinh mã từ hợp đồng, kiểm thử, đóng gói
```
