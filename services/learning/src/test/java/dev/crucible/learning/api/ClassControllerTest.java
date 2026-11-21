package dev.crucible.learning.api;

import dev.crucible.learning.application.ClassService;
import dev.crucible.learning.application.CsvImportService;
import dev.crucible.learning.domain.ClassEntity;
import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.domain.ClassRole;
import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import dev.crucible.learning.infrastructure.persistence.ClassRepository;
import dev.crucible.learning.infrastructure.runtime.generated.QuotaApi;
import dev.crucible.learning.infrastructure.security.UserContext;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.mock.web.MockMultipartFile;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;
import org.springframework.web.client.RestClientResponseException;

import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@ExtendWith(MockitoExtension.class)
class ClassControllerTest {

    private MockMvc mockMvc;

    @Mock
    private ClassRepository classRepository;

    @Mock
    private ClassMemberRepository classMemberRepository;

    @Mock
    private QuotaApi quotaApi;

    @Mock
    private UserContext userContext;

    private final UUID userId = UUID.randomUUID();
    private final UUID classId = UUID.randomUUID();
    private ClassEntity sampleClazz;

    @BeforeEach
    void setUp() {
        ClassService classService = new ClassService(classRepository, classMemberRepository, quotaApi);
        CsvImportService csvImportService = new CsvImportService(classRepository, classMemberRepository);
        ClassController controller = new ClassController(classService, csvImportService, userContext);
        mockMvc = MockMvcBuilders.standaloneSetup(controller)
                .setControllerAdvice(new GlobalExceptionHandler())
                .build();

        sampleClazz = new ClassEntity("IT4060", "Hệ điều hành", "2026.1", userId, 1200, 30);
        sampleClazz.setId(classId);
    }

    @Test
    @DisplayName("GET /v1/classes trả danh sách lớp học của người dùng")
    void getMyClasses_success() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findByMemberUserId(userId)).thenReturn(List.of(sampleClazz));

        mockMvc.perform(get("/v1/classes"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$[0].id").value(classId.toString()))
                .andExpect(jsonPath("$[0].code").value("IT4060"))
                .andExpect(jsonPath("$[0].quotaAccountRef").value("class:" + classId));
    }

    @Test
    @DisplayName("GET /v1/classes trả 401 UNAUTHENTICATED kèm traceId khi token thiếu sub hợp lệ")
    void getMyClasses_unauthenticated() throws Exception {
        when(userContext.requireUserId()).thenThrow(ProblemException.unauthenticated("Token không hợp lệ"));

        mockMvc.perform(get("/v1/classes"))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.code").value("UNAUTHENTICATED"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }

    @Test
    @DisplayName("POST /v1/classes tạo lớp học mới thành công khi có vai trò INSTRUCTOR")
    void createClass_success() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("INSTRUCTOR")).thenReturn(true);
        when(userContext.getEmail()).thenReturn(Optional.of("teacher@sis.hust.edu.vn"));
        when(userContext.getFullName()).thenReturn(Optional.of("Nguyễn Văn A"));
        when(userContext.getLocale()).thenReturn("vi");

        when(classRepository.existsByCode("IT4060")).thenReturn(false);
        when(classRepository.save(any(ClassEntity.class))).thenAnswer(inv -> inv.getArgument(0));

        String jsonPayload = """
                {
                    "code": "IT4060",
                    "name": "Hệ điều hành",
                    "semester": "2026.1",
                    "weeklyMinutes": 1200,
                    "maxConcurrentSessions": 30
                }
                """;

        mockMvc.perform(post("/v1/classes")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(jsonPayload))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.code").value("IT4060"))
                .andExpect(jsonPath("$.leadInstructorId").value(userId.toString()));
    }

    @Test
    @DisplayName("POST /v1/classes từ chối 403 FORBIDDEN khi tài khoản STUDENT gọi tạo lớp")
    void createClass_forbiddenForStudent() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("INSTRUCTOR")).thenReturn(false);
        when(userContext.hasRole("ADMIN")).thenReturn(false);

        String jsonPayload = """
                {
                    "code": "IT4060",
                    "name": "Hệ điều hành",
                    "semester": "2026.1"
                }
                """;

        mockMvc.perform(post("/v1/classes")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(jsonPayload))
                .andExpect(status().isForbidden())
                .andExpect(jsonPath("$.code").value("FORBIDDEN"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }

    @Test
    @DisplayName("POST /v1/classes chuyển tiếp 422 QUOTA_CEILING_EXCEEDED về giao diện (BR-M3.26)")
    void createClass_quotaCeilingExceeded() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("INSTRUCTOR")).thenReturn(true);
        when(userContext.getEmail()).thenReturn(Optional.of("teacher@sis.hust.edu.vn"));
        when(userContext.getFullName()).thenReturn(Optional.of("Nguyễn Văn A"));
        when(userContext.getLocale()).thenReturn("vi");

        when(classRepository.existsByCode("IT4060")).thenReturn(false);
        when(classRepository.save(any(ClassEntity.class))).thenAnswer(inv -> inv.getArgument(0));
        doThrow(new RestClientResponseException(
                "422 Unprocessable Entity", 422, "Unprocessable Entity", HttpHeaders.EMPTY,
                "{\"code\":\"QUOTA_CEILING_EXCEEDED\",\"detail\":\"Vượt trần\"}".getBytes(StandardCharsets.UTF_8),
                StandardCharsets.UTF_8
        )).when(quotaApi).putQuotaAccount(any(), any());

        String jsonPayload = """
                {
                    "code": "IT4060",
                    "name": "Hệ điều hành",
                    "semester": "2026.1",
                    "weeklyMinutes": 50000
                }
                """;

        mockMvc.perform(post("/v1/classes")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(jsonPayload))
                .andExpect(status().isUnprocessableEntity())
                .andExpect(jsonPath("$.code").value("QUOTA_CEILING_EXCEEDED"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }

    @Test
    @DisplayName("GET /v1/classes/{classId} trả chi tiết lớp học cho thành viên")
    void getClassDetail_success() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findById(classId)).thenReturn(Optional.of(sampleClazz));
        ClassMemberEntity member = new ClassMemberEntity(sampleClazz, userId, ClassRole.INSTRUCTOR, null, "GV", "gv@hust.edu.vn");
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.of(member));

        mockMvc.perform(get("/v1/classes/{classId}", classId))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(classId.toString()))
                .andExpect(jsonPath("$.code").value("IT4060"));
    }

    @Test
    @DisplayName("GET /v1/classes/{classId} trả 403 FORBIDDEN khi người ngoài lớp xem chi tiết")
    void getClassDetail_forbiddenForOutsider() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findById(classId)).thenReturn(Optional.of(sampleClazz));
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.empty());

        mockMvc.perform(get("/v1/classes/{classId}", classId))
                .andExpect(status().isForbidden())
                .andExpect(jsonPath("$.code").value("FORBIDDEN"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }

    @Test
    @DisplayName("GET /v1/classes/{classId}/members trả 403 FORBIDDEN khi người ngoài lớp xem thành viên")
    void getClassMembers_forbiddenForOutsider() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.existsById(classId)).thenReturn(true);
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.empty());

        mockMvc.perform(get("/v1/classes/{classId}/members", classId))
                .andExpect(status().isForbidden())
                .andExpect(jsonPath("$.code").value("FORBIDDEN"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }

    @Test
    @DisplayName("POST /v1/classes/{classId}/staff thêm trợ giảng thành công")
    void addClassStaff_success() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findById(classId)).thenReturn(Optional.of(sampleClazz));
        ClassMemberEntity instructor = new ClassMemberEntity(sampleClazz, userId, ClassRole.INSTRUCTOR, null, "GV", "gv@hust.edu.vn");
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.of(instructor));
        when(classMemberRepository.existsByClazzIdAndEmail(classId, "ta@sis.hust.edu.vn")).thenReturn(false);
        when(classMemberRepository.save(any(ClassMemberEntity.class))).thenAnswer(inv -> {
            ClassMemberEntity e = inv.getArgument(0);
            e.setId(UUID.randomUUID());
            return e;
        });

        String jsonPayload = """
                {
                    "email": "ta@sis.hust.edu.vn",
                    "fullName": "Trần Trợ Giảng",
                    "role": "TEACHING_ASSISTANT"
                }
                """;

        mockMvc.perform(post("/v1/classes/{classId}/staff", classId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(jsonPayload))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.role").value("TEACHING_ASSISTANT"))
                .andExpect(jsonPath("$.email").value("ta@sis.hust.edu.vn"));
    }

    @Test
    @DisplayName("POST /v1/classes/{classId}/staff trả 403 FORBIDDEN khi trợ giảng gọi thêm thành viên")
    void addClassStaff_forbiddenForTeachingAssistant() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findById(classId)).thenReturn(Optional.of(sampleClazz));
        ClassMemberEntity taMember = new ClassMemberEntity(sampleClazz, userId, ClassRole.TEACHING_ASSISTANT, null, "TA", "ta@hust.edu.vn");
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.of(taMember));

        String jsonPayload = """
                {
                    "email": "ta2@sis.hust.edu.vn",
                    "fullName": "Trợ Giảng 2",
                    "role": "TEACHING_ASSISTANT"
                }
                """;

        mockMvc.perform(post("/v1/classes/{classId}/staff", classId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(jsonPayload))
                .andExpect(status().isForbidden())
                .andExpect(jsonPath("$.code").value("FORBIDDEN"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }

    @Test
    @DisplayName("POST /v1/classes/{classId}/members/csv-preview thành công cho giảng viên của lớp")
    void previewCsvMembers_success() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findById(classId)).thenReturn(Optional.of(sampleClazz));
        ClassMemberEntity instructor = new ClassMemberEntity(sampleClazz, userId, ClassRole.INSTRUCTOR, null, "GV", "gv@hust.edu.vn");
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.of(instructor));
        when(classMemberRepository.findByClazzId(classId)).thenReturn(List.of(instructor));

        MockMultipartFile file = new MockMultipartFile(
                "file", "students.csv", "text/csv",
                "student_id,full_name,email\n20220001,Nguyen Van A,a@sis.hust.edu.vn".getBytes(StandardCharsets.UTF_8)
        );

        mockMvc.perform(multipart("/v1/classes/{classId}/members/csv-preview", classId)
                        .file(file))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.totalRows").value(1))
                .andExpect(jsonPath("$.validRows").value(1))
                .andExpect(jsonPath("$.errorRows").value(0));
    }

    @Test
    @DisplayName("POST /v1/classes/{classId}/members/csv-preview trả 403 FORBIDDEN khi người ngoài lớp xem trước CSV")
    void previewCsvMembers_forbiddenForOutsider() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findById(classId)).thenReturn(Optional.of(sampleClazz));
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.empty());

        MockMultipartFile file = new MockMultipartFile(
                "file", "students.csv", "text/csv",
                "student_id,full_name,email\n20220001,Nguyen Van A,a@sis.hust.edu.vn".getBytes(StandardCharsets.UTF_8)
        );

        mockMvc.perform(multipart("/v1/classes/{classId}/members/csv-preview", classId)
                        .file(file))
                .andExpect(status().isForbidden())
                .andExpect(jsonPath("$.code").value("FORBIDDEN"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }

    @Test
    @DisplayName("POST /v1/classes/{classId}/members/csv-import bắt DataIntegrityViolationException trả 409 CONFLICT")
    void importCsvMembers_concurrentConflictReturns409() throws Exception {
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findById(classId)).thenReturn(Optional.of(sampleClazz));
        ClassMemberEntity instructor = new ClassMemberEntity(sampleClazz, userId, ClassRole.INSTRUCTOR, null, "GV", "gv@hust.edu.vn");
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.of(instructor));
        when(classMemberRepository.findByClazzId(classId)).thenReturn(List.of(instructor));
        doThrow(new DataIntegrityViolationException("uk_class_member_email_lower"))
                .when(classMemberRepository).saveAll(any());

        MockMultipartFile file = new MockMultipartFile(
                "file", "students.csv", "text/csv",
                "student_id,full_name,email\n20220001,Nguyen Van A,a@sis.hust.edu.vn".getBytes(StandardCharsets.UTF_8)
        );

        mockMvc.perform(multipart("/v1/classes/{classId}/members/csv-import", classId)
                        .file(file))
                .andExpect(status().isConflict())
                .andExpect(jsonPath("$.code").value("MEMBER_ALREADY_EXISTS"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }

    @Test
    @DisplayName("DELETE /v1/classes/{classId}/members/{memberId} trả 403 FORBIDDEN khi không phải giảng viên")
    void removeClassMember_forbidden() throws Exception {
        UUID memberId = UUID.randomUUID();
        when(userContext.requireUserId()).thenReturn(userId);
        when(userContext.hasRole("ADMIN")).thenReturn(false);
        when(classRepository.findById(classId)).thenReturn(Optional.of(sampleClazz));
        ClassMemberEntity student = new ClassMemberEntity(sampleClazz, userId, ClassRole.STUDENT, "20220001", "SV", "sv@hust.edu.vn");
        when(classMemberRepository.findByClazzIdAndUserId(classId, userId)).thenReturn(Optional.of(student));

        mockMvc.perform(delete("/v1/classes/{classId}/members/{memberId}", classId, memberId))
                .andExpect(status().isForbidden())
                .andExpect(jsonPath("$.code").value("FORBIDDEN"))
                .andExpect(jsonPath("$.traceId").isNotEmpty());
    }
}
