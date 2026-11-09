package dev.crucible.learning.application;

import dev.crucible.learning.api.ProblemException;
import dev.crucible.learning.api.dto.CsvImportPreviewResponse;
import dev.crucible.learning.domain.ClassEntity;
import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.domain.ClassRole;
import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import dev.crucible.learning.infrastructure.persistence.ClassRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class CsvImportServiceTest {

    @Mock
    private ClassRepository classRepository;

    @Mock
    private ClassMemberRepository classMemberRepository;

    @InjectMocks
    private CsvImportService csvImportService;

    private UUID classId;
    private ClassEntity mockClass;

    @BeforeEach
    void setUp() {
        classId = UUID.randomUUID();
        mockClass = new ClassEntity("IT4060", "Linux OS", "2026.1", null, 1);
        mockClass.setId(classId);
    }

    @Test
    @DisplayName("AC-L1.5, BR-M3.05: Kiểm tra CSV hợp lệ trả về toàn bộ danh sách hợp lệ và 0 lỗi")
    void parseCsv_shouldSucceed_whenCsvIsValid_AC_L1_5() {
        when(classRepository.findById(classId)).thenReturn(Optional.of(mockClass));
        when(classMemberRepository.findByClazzId(classId)).thenReturn(List.of());

        String csv = """
                student_id,full_name,email
                20210001,Nguyen Van A,a.nv210001@sis.hust.edu.vn
                20210002,Tran Thi B,b.tt210002@sis.hust.edu.vn
                """;
        InputStream is = new ByteArrayInputStream(csv.getBytes(StandardCharsets.UTF_8));

        CsvImportPreviewResponse result = csvImportService.parseAndValidate(classId, is);

        assertThat(result.totalRows()).isEqualTo(2);
        assertThat(result.validRows()).isEqualTo(2);
        assertThat(result.errorRows()).isEqualTo(0);
        assertThat(result.errors()).isEmpty();
        assertThat(result.validStudents()).hasSize(2);
    }

    @Test
    @DisplayName("AC-L1.5, BR-M3.05: Phát hiện đúng lỗi theo từng dòng: thiếu cột, sai định dạng email, trùng mã SV")
    void parseCsv_shouldReportLineErrors_whenMissingFields_orInvalidEmail_orDuplicateStudentId_AC_L1_5_BR_M3_05() {
        when(classRepository.findById(classId)).thenReturn(Optional.of(mockClass));
        when(classMemberRepository.findByClazzId(classId)).thenReturn(List.of(
                new ClassMemberEntity(mockClass, null, ClassRole.STUDENT, "20219999", "Existed Student", "existed@sis.hust.edu.vn")
        ));

        // Dòng 2: Hợp lệ
        // Dòng 3: Thiếu họ tên
        // Dòng 4: Email sai định dạng RFC 5322
        // Dòng 5: Trùng mã SV với dòng 2 trong tệp
        // Dòng 6: Trùng mã SV đã có trong lớp (20219999)
        String csv = """
                student_id,full_name,email
                20210001,Nguyen Van A,a.nv210001@sis.hust.edu.vn
                20210002,,b.tt210002@sis.hust.edu.vn
                20210003,Le Van C,invalid-email-format
                20210001,Pham Van D,d.pv210001@sis.hust.edu.vn
                20219999,Hoang Van E,e.hv219999@sis.hust.edu.vn
                """;
        InputStream is = new ByteArrayInputStream(csv.getBytes(StandardCharsets.UTF_8));

        CsvImportPreviewResponse result = csvImportService.parseAndValidate(classId, is);

        assertThat(result.validRows()).isEqualTo(1);
        assertThat(result.errorRows()).isEqualTo(4);
        assertThat(result.errors()).hasSize(4);

        // Kiểm tra dòng 3
        assertThat(result.errors()).anyMatch(e -> e.line() == 3 && e.field().equals("full_name"));
        // Kiểm tra dòng 4
        assertThat(result.errors()).anyMatch(e -> e.line() == 4 && e.field().equals("email"));
        // Kiểm tra dòng 5
        assertThat(result.errors()).anyMatch(e -> e.line() == 5 && e.field().equals("student_id") && e.message().contains("Trùng mã"));
        // Kiểm tra dòng 6
        assertThat(result.errors()).anyMatch(e -> e.line() == 6 && e.field().equals("student_id") && e.message().contains("đã có trong lớp"));
    }

    @Test
    @DisplayName("AC-L1.5: Báo lỗi khi tiêu đề CSV thiếu cột bắt buộc")
    void parseCsv_shouldThrowBadRequest_whenHeaderInvalid_AC_L1_5() {
        when(classRepository.findById(classId)).thenReturn(Optional.of(mockClass));

        String csv = """
                mssv,ten
                20210001,Nguyen Van A
                """;
        InputStream is = new ByteArrayInputStream(csv.getBytes(StandardCharsets.UTF_8));

        assertThatThrownBy(() -> csvImportService.parseAndValidate(classId, is))
                .isInstanceOf(ProblemException.class)
                .hasMessageContaining("CSV_HEADER_INVALID");
    }
}
