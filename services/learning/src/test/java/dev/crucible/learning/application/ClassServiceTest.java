package dev.crucible.learning.application;

import dev.crucible.learning.api.ProblemException;
import dev.crucible.learning.api.dto.AddInstructorRequest;
import dev.crucible.learning.api.dto.ClassMemberResponse;
import dev.crucible.learning.api.dto.ClassResponse;
import dev.crucible.learning.api.dto.CreateClassRequest;
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

import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class ClassServiceTest {

    @Mock
    private ClassRepository classRepository;

    @Mock
    private ClassMemberRepository classMemberRepository;

    @Mock
    private UserAutoEnrollService userAutoEnrollService;

    @InjectMocks
    private ClassService classService;

    private UUID instructorId;

    @BeforeEach
    void setUp() {
        instructorId = UUID.randomUUID();
    }

    @Test
    @DisplayName("L01: Tạo lớp học thành công và tự động gán người tạo làm INSTRUCTOR")
    void createClass_shouldSucceed_andAssignCreatorAsInstructor_L01() {
        CreateClassRequest req = new CreateClassRequest("IT4060-2026.1", "Hệ điều hành Linux", "2026.1", "quota-1", 1);
        when(classRepository.existsByCode(req.code())).thenReturn(false);
        when(classRepository.save(any(ClassEntity.class))).thenAnswer(invocation -> {
            ClassEntity entity = invocation.getArgument(0);
            entity.setId(UUID.randomUUID());
            return entity;
        });

        ClassResponse res = classService.createClass(req, instructorId, "Thầy A", "a@hust.edu.vn");

        assertThat(res.code()).isEqualTo("IT4060-2026.1");
        assertThat(res.name()).isEqualTo("Hệ điều hành Linux");
        verify(classRepository).save(any(ClassEntity.class));
    }

    @Test
    @DisplayName("L01: Tạo lớp thất bại nếu mã lớp đã tồn tại (409 CLASS_CODE_ALREADY_EXISTS)")
    void createClass_shouldThrowConflict_whenCodeExists_L01() {
        CreateClassRequest req = new CreateClassRequest("IT4060-2026.1", "Hệ điều hành Linux", "2026.1", null, 1);
        when(classRepository.existsByCode(req.code())).thenReturn(true);

        assertThatThrownBy(() -> classService.createClass(req, instructorId, "Thầy A", "a@hust.edu.vn"))
                .isInstanceOf(ProblemException.class)
                .hasMessageContaining("CLASS_CODE_ALREADY_EXISTS");
    }

    @Test
    @DisplayName("L01: Giảng viên phụ trách thêm trợ giảng thành công")
    void addAssistant_shouldSucceed_whenRequestedByInstructor_L01() {
        UUID classId = UUID.randomUUID();
        ClassEntity clazz = new ClassEntity("IT4060", "Linux", "2026.1", null, 1);
        clazz.setId(classId);

        ClassMemberEntity requester = new ClassMemberEntity(clazz, instructorId, ClassRole.INSTRUCTOR, null, "Thầy A", "a@hust.edu.vn");

        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classMemberRepository.findByClazzIdAndUserId(classId, instructorId)).thenReturn(Optional.of(requester));
        when(classMemberRepository.existsByClazzIdAndEmail(classId, "ta@sis.hust.edu.vn")).thenReturn(false);
        when(classMemberRepository.save(any(ClassMemberEntity.class))).thenAnswer(i -> {
            ClassMemberEntity m = i.getArgument(0);
            m.setId(UUID.randomUUID());
            return m;
        });

        AddInstructorRequest addReq = new AddInstructorRequest("ta@sis.hust.edu.vn", "Trợ giảng B", ClassRole.TEACHING_ASSISTANT);
        ClassMemberResponse res = classService.addInstructorOrAssistant(classId, addReq, instructorId, false);

        assertThat(res.email()).isEqualTo("ta@sis.hust.edu.vn");
        assertThat(res.role()).isEqualTo(ClassRole.TEACHING_ASSISTANT);
    }

    @Test
    @DisplayName("L01: Trợ giảng không có quyền thêm giảng viên hoặc trợ giảng khác")
    void addAssistant_shouldThrowForbidden_whenRequestedByTeachingAssistant_L01() {
        UUID classId = UUID.randomUUID();
        ClassEntity clazz = new ClassEntity("IT4060", "Linux", "2026.1", null, 1);
        UUID taId = UUID.randomUUID();
        ClassMemberEntity ta = new ClassMemberEntity(clazz, taId, ClassRole.TEACHING_ASSISTANT, null, "Trợ giảng B", "ta@hust.edu.vn");

        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classMemberRepository.findByClazzIdAndUserId(classId, taId)).thenReturn(Optional.of(ta));

        AddInstructorRequest addReq = new AddInstructorRequest("ta2@hust.edu.vn", "Trợ giảng C", ClassRole.TEACHING_ASSISTANT);

        assertThatThrownBy(() -> classService.addInstructorOrAssistant(classId, addReq, taId, false))
                .isInstanceOf(ProblemException.class)
                .hasMessageContaining("FORBIDDEN");
    }
}
