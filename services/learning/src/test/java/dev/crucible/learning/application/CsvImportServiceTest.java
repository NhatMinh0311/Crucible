package dev.crucible.learning.application;

import dev.crucible.learning.api.ProblemException;
import dev.crucible.learning.api.generated.model.CsvImportPreviewResponse;
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
    private UUID instructorId;
    private ClassEntity clazz;

    @BeforeEach
    void setUp() {
        classId = UUID.randomUUID();
        instructorId = UUID.randomUUID();
        clazz = new ClassEntity("IT4060", "Hệ điều hành", "2026.1", instructorId, 1200, 30);
        clazz.setId(classId);
    }

    @Test
    @DisplayName("AC_L1_5: Phân tích CSV phát hiện lỗi theo từng dòng (thiếu cột, sai email, trùng mã SV trong tệp và trong lớp)")
    void parseAndValidate_detectsLineByLineErrors_AC_L1_5() {
        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));

        ClassMemberEntity existingStudent = new ClassMemberEntity(
                clazz, null, ClassRole.STUDENT, "20210001", "Đã Có", "existing@sis.hust.edu.vn"
        );
        when(classMemberRepository.findByClazzId(classId)).thenReturn(List.of(existingStudent));

        String csvContent = """
                student_id,full_name,email
                20210002,Nguyễn Văn Hợp Lệ,hop.le@sis.hust.edu.vn
                ,Trần Thiếu Mã,thieu.ma@sis.hust.edu.vn
                20210003,Lê Sai Email,email-khong-hop-le
                20210002,Phạm Trùng Mã Tệp,trung.tep@sis.hust.edu.vn
                20210001,Hoàng Trùng Mã Lớp,trung.lop@sis.hust.edu.vn
                20210004,Đỗ Trùng Email Lớp,existing@sis.hust.edu.vn
                """;

        CsvImportPreviewResponse result = csvImportService.parseAndValidate(
                classId, new ByteArrayInputStream(csvContent.getBytes(StandardCharsets.UTF_8))
        );

        assertThat(result.getTotalRows()).isEqualTo(6);
        assertThat(result.getValidRows()).isEqualTo(1);
        assertThat(result.getErrorRows()).isEqualTo(5);
        assertThat(result.getValidStudents()).hasSize(1);
        assertThat(result.getValidStudents().getFirst().getStudentId()).isEqualTo("20210002");

        assertThat(result.getErrors()).anyMatch(e -> e.getLine() == 3 && e.getField().equals("student_id"));
        assertThat(result.getErrors()).anyMatch(e -> e.getLine() == 4 && e.getField().equals("email"));
        assertThat(result.getErrors()).anyMatch(e -> e.getLine() == 5 && e.getField().equals("student_id"));
        assertThat(result.getErrors()).anyMatch(e -> e.getLine() == 6 && e.getField().equals("student_id"));
        assertThat(result.getErrors()).anyMatch(e -> e.getLine() == 7 && e.getField().equals("email"));
    }

    @Test
    @DisplayName("AC_L1_5: Đọc tệp CSV UTF-8 có BOM (EF BB BF) xuất từ Excel không bị báo lỗi tiêu đề")
    void parseAndValidate_handlesUtf8BomFromExcel_AC_L1_5() {
        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classMemberRepository.findByClazzId(classId)).thenReturn(List.of());

        String csvWithBom = "\uFEFFstudent_id,full_name,email\r\n20225001,Nguyễn Văn Ánh,anh.nv225001@sis.hust.edu.vn\r\n";
        byte[] bytes = csvWithBom.getBytes(StandardCharsets.UTF_8);

        CsvImportPreviewResponse result = csvImportService.parseAndValidate(
                classId, new ByteArrayInputStream(bytes)
        );

        assertThat(result.getTotalRows()).isEqualTo(1);
        assertThat(result.getValidRows()).isEqualTo(1);
        assertThat(result.getErrorRows()).isEqualTo(0);
        assertThat(result.getValidStudents().getFirst().getStudentId()).isEqualTo("20225001");
        assertThat(result.getValidStudents().getFirst().getFullName()).isEqualTo("Nguyễn Văn Ánh");
    }

    @Test
    @DisplayName("Từ chối tệp CSV thiếu cột tiêu đề bắt buộc với mã VALIDATION_FAILED")
    void parseAndValidate_invalidHeader_throwsValidationFailed() {
        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classMemberRepository.findByClazzId(classId)).thenReturn(List.of());

        String badCsv = "id,ten\n1,Nguyen A";

        assertThatThrownBy(() -> csvImportService.parseAndValidate(
                classId, new ByteArrayInputStream(badCsv.getBytes(StandardCharsets.UTF_8))
        ))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> assertThat(((ProblemException) ex).getCode()).isEqualTo("VALIDATION_FAILED"));
    }

    @Test
    @DisplayName("Xem trước CSV kiểm tra quyền giảng viên của lớp: người ngoài lớp hoặc trợ giảng bị từ chối 403 FORBIDDEN")
    void previewImport_nonInstructor_throwsForbidden() {
        UUID outsiderId = UUID.randomUUID();
        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classMemberRepository.findByClazzIdAndUserId(classId, outsiderId)).thenReturn(Optional.empty());

        String csvContent = "student_id,full_name,email\n20210002,Nguyen A,a@sis.hust.edu.vn";

        assertThatThrownBy(() -> csvImportService.previewImport(
                classId,
                new ByteArrayInputStream(csvContent.getBytes(StandardCharsets.UTF_8)),
                outsiderId,
                false
        ))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> assertThat(((ProblemException) ex).getCode()).isEqualTo("FORBIDDEN"));
    }
}
