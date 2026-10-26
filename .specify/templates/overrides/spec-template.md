# Đặc tả: [TÊN TÍNH NĂNG]

| Thuộc tính | Giá trị |
|---|---|
| Mã | [###-ten-ngan] |
| Module | [M1, M2, M3 hoặc Nền tảng] |
| Mốc | [MA, MB, MC hoặc MD] |
| Trạng thái | Nháp |
| Ngày lập | [NGÀY] |
| Issue | [#số] |
| Tài liệu liên quan | [Mã chức năng, `US-*`, `AC-*` trong CR-DOC-01; mục của tài liệu module; `BR-*`] |
| Yêu cầu ban đầu | "$ARGUMENTS" |

<!--
  Đặc tả trả lời "làm gì, vì sao", không nói cách làm (ngôn ngữ, thư viện, cấu trúc mã thuộc plan.md).
  Dẫn mã AC-*, BR-*, US-* của tài liệu đã duyệt thay vì viết lại nội dung. Mục nào không áp dụng thì xóa hẳn.
-->

## 1. Mục tiêu

[Hai, ba câu: ai được lợi, giải quyết vấn đề gì, vì sao cần ở mốc này]

## 2. Phạm vi

| Thuộc phạm vi | Ngoài phạm vi |
|---|---|
| [Nội dung làm trong tính năng này] | [Nội dung không làm, kèm issue hoặc module phụ trách nếu có] |

## 3. Câu chuyện người dùng

<!--
  Xếp theo mức ưu tiên P1, P2, P3. Mỗi câu chuyện kiểm thử được độc lập: chỉ hiện thực một câu chuyện vẫn cho ra
  một phần dùng được. Nhãn US1, US2... dùng trong tasks.md; mã US-* trong ngoặc là câu chuyện của CR-DOC-01 nếu có.
-->

### US1 — [Tên ngắn] (ưu tiên P1; [US-* của CR-DOC-01])

[Mô tả bằng lời của người dùng: vai trò, việc cần làm, kết quả mong muốn]

**Vì sao ưu tiên này**: [Giá trị mang lại]

**Kiểm thử độc lập**: [Cách kiểm chứng riêng câu chuyện này]

**Kịch bản chấp nhận**:

1. **Cho** [trạng thái ban đầu], **khi** [hành động], **thì** [kết quả]
2. **Cho** [trạng thái ban đầu], **khi** [hành động], **thì** [kết quả]

---

### US2 — [Tên ngắn] (ưu tiên P2)

[Mô tả]

**Vì sao ưu tiên này**: [Giá trị]

**Kiểm thử độc lập**: [Cách kiểm chứng]

**Kịch bản chấp nhận**:

1. **Cho** [trạng thái ban đầu], **khi** [hành động], **thì** [kết quả]

## 4. Yêu cầu

<!--
  Mỗi yêu cầu kiểm thử được và dẫn mã tiêu chí hoặc quy tắc gốc. Chỗ tài liệu chưa trả lời được thì đánh dấu
  [NEEDS CLARIFICATION: câu hỏi], tối đa 3 chỗ.
-->

| Mã | Yêu cầu | Nguồn |
|---|---|---|
| FR-001 | Hệ thống phải [năng lực cụ thể] | [AC-*, BR-*] |
| FR-002 | Người dùng [vai trò] phải [tương tác cụ thể] | [AC-*] |
| FR-003 | Hệ thống phải [năng lực] [NEEDS CLARIFICATION: câu hỏi cụ thể] | [AC-*] |

### Thực thể chính

<!-- Chỉ khi tính năng có dữ liệu. Nêu ý nghĩa và quan hệ, không nêu kiểu cột hay bảng. -->

- **[Thực thể]**: [Ý nghĩa, thuộc tính chính, quan hệ với thực thể khác]

## 5. Tiêu chí thành công

<!-- Đo được, không phụ thuộc công nghệ. Dẫn chỉ tiêu Q* của CR-DOC-01 khi có. -->

| Mã | Tiêu chí | Cách đo |
|---|---|---|
| SC-001 | [Kết quả đo được, ví dụ: học viên vào được terminal trong vòng 5 giây kể từ khi bấm bắt đầu] | [Kiểm thử, số đo] |
| SC-002 | [Kết quả đo được] | [Cách đo] |

## 6. Trường hợp biên

| Trường hợp | Xử lý |
|---|---|
| [Điều kiện biên, lỗi, thao tác đồng thời] | [Kết quả mong đợi] |

## 7. Giả định và phụ thuộc

- [Giả định khi yêu cầu không nói rõ, kèm lý do chọn mặc định]
- [Phụ thuộc vào module khác hoặc issue khác, kèm trạng thái]

## 8. Điểm đã làm rõ

<!-- Lệnh clarify ghi vào mục này (tương ứng mục Clarifications của Spec Kit), mỗi buổi một tiểu mục. -->

### Phiên [NGÀY]

- Hỏi: [câu hỏi] → Trả lời: [quyết định của thành viên]
