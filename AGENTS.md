# Hướng dẫn cho trợ lý AI

Tệp này quy định cách trợ lý AI (Claude Code, Copilot, Cursor, Codex và các công cụ tương tự) làm việc trong kho mã
Crucible. Mọi thành viên dùng chung tệp này; thay đổi tệp đi qua pull request như mã nguồn.

## 1. Đọc trước khi làm

| Việc | Tài liệu |
|---|---|
| Hiểu bài toán, chức năng, tiêu chí chấp nhận | `docs/01-bai-toan-va-pham-vi.md` |
| Hiểu ranh giới module, giao tiếp | `docs/02-kien-truc.md` |
| Gọi hoặc cung cấp API, sự kiện giữa module | `docs/03-hop-dong-giua-module.md`, `contracts/` |
| Thiết kế bên trong module đang làm | `docs/modules/<module>.md` |
| Quy trình nhánh, commit, pull request | `docs/07-quy-trinh-phat-trien.md`, `CONTRIBUTING.md` |
| Quyết định kiến trúc đã chốt | `docs/adr/` |

## 2. Quy tắc bắt buộc

1. **Chỉ sửa mã trong module được giao.** Cần thay đổi ở module khác hoặc ở `contracts/` thì dừng lại và báo người dùng;
   thay đổi hợp đồng là pull request riêng (CR-DOC-03 mục 2.1).
2. **Không tự quyết định nghiệp vụ.** Khi tiêu chí chấp nhận, quy tắc `BR-*` hoặc tài liệu module không trả lời được
   một trường hợp, hỏi người dùng và đề xuất phương án kèm lý do; không tự chọn rồi viết mã.
3. **Không truy cập dữ liệu của module khác** ngoài API và sự kiện trong `contracts/`.
4. **Mỗi tiêu chí chấp nhận hoặc quy tắc nghiệp vụ được hiện thực phải có kiểm thử** mang mã của nó trong tên kiểm thử
   hoặc chú thích (ví dụ `AC-C2.4`, `BR-M1.3`).
5. **Không đưa bí mật vào mã**, không sửa `.env`, không tắt kiểm tra bảo mật hoặc kiểm thử đang thất bại để CI qua.
6. **Không sửa migration đã merge**; thay đổi lược đồ bằng migration mới theo nguyên tắc mở rộng rồi thu hẹp.
7. **Chạy kiểm tra cục bộ** (CONTRIBUTING mục 2.4 và lệnh của module) trước khi báo hoàn thành; báo đúng kết quả,
   kể cả khi thất bại.

## 3. Quy trình theo đặc tả (Spec Kit)

Mỗi tính năng mới đi qua ba tệp trong `specs/<số>-<tên-ngắn>/`, được duyệt trước khi viết mã (CR-DOC-07 mục 7):

| Bước | Tệp | Nội dung | Người duyệt |
|---|---|---|---|
| Đặc tả | `spec.md` | Mục tiêu, tiêu chí chấp nhận (dẫn mã `AC-*`), trường hợp biên, điểm cần làm rõ | Chủ module và trưởng nhóm |
| Kế hoạch | `plan.md` | Thiết kế kỹ thuật trong phạm vi module, ảnh hưởng tới hợp đồng, dữ liệu, kiểm thử | Chủ module |
| Danh sách việc | `tasks.md` | Các việc nhỏ có thứ tự, mỗi việc kiểm chứng được | Chủ module |

Mục "điểm cần làm rõ" trong `spec.md` phải được người dùng trả lời trước khi lập kế hoạch.

## 4. Quy ước mã nguồn

| Phạm vi | Quy ước |
|---|---|
| Go (M1) | `gofmt`, `go vet`, `staticcheck`; lỗi được bọc bằng `fmt.Errorf("...: %w", err)`; không `panic` trong luồng xử lý yêu cầu |
| Java (M2, M3) | Spring Boot 4, Java 25; tầng `api` (sinh từ OpenAPI), `application`, `domain`, `infrastructure`; không gọi repository từ controller |
| Giao diện | React, TypeScript strict, shadcn/ui, Tailwind CSS; kiểu dữ liệu API sinh từ `contracts/`; văn bản hiển thị qua khóa dịch, đủ tiếng Việt và tiếng Anh (CR-DOC-11 mục 4.1) |
| Thiết kế giao diện | Dùng skill `frontend-design` khi tạo màn hình mới; theme chung tại `web/src/shell/theme`; phông Be Vietnam Pro; biểu tượng Lucide |
| Log | JSON, có `trace_id`; không ghi token, vé terminal, nội dung lời giải |
| Thời gian | Lưu UTC, hiển thị `Asia/Ho_Chi_Minh` |

## 5. Văn bản

- Tài liệu trong `docs/` viết bằng tiếng Việt, văn phong chính thức, không chứa câu hỏi hay ghi chú trao đổi.
- Commit: dòng đầu tối đa 72 ký tự, mô tả thay đổi; thân ghi lý do và mã yêu cầu liên quan.
- Thành viên ghi nhật ký sử dụng AI tại `docs/08-nhat-ky-su-dung-ai.md`.
