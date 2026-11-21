package dev.crucible.learning.application;

import com.sun.net.httpserver.HttpServer;
import dev.crucible.learning.api.ProblemException;
import dev.crucible.learning.api.generated.model.*;
import dev.crucible.learning.domain.ClassEntity;
import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import dev.crucible.learning.infrastructure.persistence.ClassRepository;
import dev.crucible.learning.infrastructure.runtime.generated.QuotaApi;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.http.HttpStatus;
import org.springframework.web.client.RestClient;
import org.springframework.web.client.support.RestClientAdapter;
import org.springframework.web.service.invoker.HttpServiceProxyFactory;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class ClassServiceTest {

    @Mock
    private ClassRepository classRepository;

    @Mock
    private ClassMemberRepository classMemberRepository;

    private HttpServer mockRuntimeServer;
    private final AtomicReference<String> lastQuotaPath = new AtomicReference<>();
    private final AtomicReference<String> lastQuotaBody = new AtomicReference<>();
    private final AtomicInteger mockResponseStatus = new AtomicInteger(200);
    private final AtomicReference<String> mockResponseBody = new AtomicReference<>("""
            {
              "accountRef": "class:test",
              "weeklyMinutes": 1200,
              "usedMinutesThisWeek": 0,
              "maxConcurrentSessions": 30,
              "activeSessions": 0
            }
            """);

    private ClassService classService;
    private UUID instructorId;

    @BeforeEach
    void setUp() throws IOException {
        instructorId = UUID.randomUUID();
        mockRuntimeServer = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        mockRuntimeServer.createContext("/v1/quota-accounts/", exchange -> {
            lastQuotaPath.set(exchange.getRequestURI().getPath());
            lastQuotaBody.set(new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8));
            int status = mockResponseStatus.get();
            byte[] respBytes = mockResponseBody.get().getBytes(StandardCharsets.UTF_8);
            exchange.getResponseHeaders().set("Content-Type", status >= 400 ? "application/problem+json" : "application/json");
            exchange.sendResponseHeaders(status, respBytes.length);
            try (OutputStream os = exchange.getResponseBody()) {
                os.write(respBytes);
            }
        });
        mockRuntimeServer.start();

        String baseUrl = "http://127.0.0.1:" + mockRuntimeServer.getAddress().getPort();
        RestClient restClient = RestClient.builder().baseUrl(baseUrl).build();
        QuotaApi quotaApi = HttpServiceProxyFactory.builderFor(RestClientAdapter.create(restClient))
                .build()
                .createClient(QuotaApi.class);

        classService = new ClassService(classRepository, classMemberRepository, quotaApi);
    }

    @AfterEach
    void tearDown() {
        if (mockRuntimeServer != null) {
            mockRuntimeServer.stop(0);
        }
    }

    @Test
    @DisplayName("BR_M3_26 AC_R2_6: Tạo lớp học mới gọi putQuotaAccount sang M1 với parentRef = instructor:<leadInstructorId>")
    void createClass_callsPutQuotaAccountAndSetsLeadInstructor() {
        CreateClassRequest request = new CreateClassRequest("IT4060-20261", "Hệ điều hành Linux", "2026.1")
                .weeklyMinutes(1500)
                .maxConcurrentSessions(40);
        when(classRepository.existsByCode("IT4060-20261")).thenReturn(false);
        when(classRepository.save(any(ClassEntity.class))).thenAnswer(inv -> inv.getArgument(0));

        ClassItem response = classService.createClass(request, instructorId, "Nguyễn Văn A", "a.nv@sis.hust.edu.vn", "vi");

        assertThat(response.getCode()).isEqualTo("IT4060-20261");
        assertThat(response.getLeadInstructorId()).isEqualTo(instructorId);
        assertThat(response.getQuotaAccountRef()).isEqualTo("class:" + response.getId());
        assertThat(response.getWeeklyMinutes()).isEqualTo(1500);
        assertThat(response.getMaxConcurrentSessions()).isEqualTo(40);
        assertThat(response.getMemberCount()).isEqualTo(1);

        assertThat(lastQuotaPath.get()).isEqualTo("/v1/quota-accounts/class:" + response.getId());
        assertThat(lastQuotaBody.get())
                .contains("\"parentRef\":\"instructor:" + instructorId + "\"")
                .contains("\"weeklyMinutes\":1500")
                .contains("\"maxConcurrentSessions\":40");
    }

    @Test
    @DisplayName("BR_M3_26 AC_R2_8: Khi M1 trả 422 QUOTA_CEILING_EXCEEDED, M3 chuyển tiếp nguyên mã lỗi về giao diện")
    void createClass_whenQuotaCeilingExceeded_forwards422() {
        CreateClassRequest request = new CreateClassRequest("IT4060-20261", "Hệ điều hành Linux", "2026.1")
                .weeklyMinutes(90000)
                .maxConcurrentSessions(30);
        when(classRepository.existsByCode("IT4060-20261")).thenReturn(false);
        when(classRepository.save(any(ClassEntity.class))).thenAnswer(inv -> inv.getArgument(0));

        mockResponseStatus.set(422);
        mockResponseBody.set("""
                {
                  "type": "about:blank",
                  "title": "Unprocessable Entity",
                  "status": 422,
                  "code": "QUOTA_CEILING_EXCEEDED",
                  "detail": "Weekly minutes exceed parent ceiling"
                }
                """);

        assertThatThrownBy(() -> classService.createClass(request, instructorId, "Nguyễn Văn A", "a.nv@sis.hust.edu.vn", "vi"))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> {
                    ProblemException pe = (ProblemException) ex;
                    assertThat(pe.getStatus()).isEqualTo(HttpStatus.UNPROCESSABLE_CONTENT);
                    assertThat(pe.getCode()).isEqualTo("QUOTA_CEILING_EXCEEDED");
                });
    }

    @Test
    @DisplayName("BR_M3_26: Khi M1 trả 501 NOT_IMPLEMENTED hoặc 403 FORBIDDEN, M3 chặn tạo lớp và chuyển tiếp đúng mã lỗi")
    void createClass_whenRuntimeReturns501_throwsNotImplemented() {
        CreateClassRequest request = new CreateClassRequest("IT4060-20261", "Hệ điều hành Linux", "2026.1");
        when(classRepository.existsByCode("IT4060-20261")).thenReturn(false);
        when(classRepository.save(any(ClassEntity.class))).thenAnswer(inv -> inv.getArgument(0));

        mockResponseStatus.set(501);
        mockResponseBody.set("""
                {
                  "type": "about:blank",
                  "title": "Not Implemented",
                  "status": 501,
                  "code": "NOT_IMPLEMENTED",
                  "detail": "Operation chưa được hiện thực"
                }
                """);

        assertThatThrownBy(() -> classService.createClass(request, instructorId, "Nguyễn Văn A", "a.nv@sis.hust.edu.vn", "vi"))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> {
                    ProblemException pe = (ProblemException) ex;
                    assertThat(pe.getStatus()).isEqualTo(HttpStatus.NOT_IMPLEMENTED);
                    assertThat(pe.getCode()).isEqualTo("NOT_IMPLEMENTED");
                });
    }

    @Test
    @DisplayName("BR_M3_26: Đổi giảng viên phụ trách chính cập nhật parentRef bên M1")
    void updateClass_changeLeadInstructor_updatesParentRefInM1() {
        UUID classId = UUID.randomUUID();
        UUID newLeadId = UUID.randomUUID();
        ClassEntity clazz = new ClassEntity("IT4060", "Hệ điều hành", "2026.1", instructorId, 1200, 30);
        clazz.setId(classId);

        ClassMemberEntity leadMember = new ClassMemberEntity(
                clazz, instructorId, dev.crucible.learning.domain.ClassRole.INSTRUCTOR, null, "GV 1", "gv1@hust.edu.vn");
        ClassMemberEntity secondInstructor = new ClassMemberEntity(
                clazz, newLeadId, dev.crucible.learning.domain.ClassRole.INSTRUCTOR, null, "GV 2", "gv2@hust.edu.vn");

        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classMemberRepository.findByClazzIdAndUserId(classId, instructorId)).thenReturn(Optional.of(leadMember));
        when(classMemberRepository.findByClazzIdAndUserId(classId, newLeadId)).thenReturn(Optional.of(secondInstructor));
        when(classRepository.save(any(ClassEntity.class))).thenAnswer(inv -> inv.getArgument(0));

        UpdateClassRequest updateReq = new UpdateClassRequest().leadInstructorId(newLeadId).weeklyMinutes(1800);
        ClassItem updated = classService.updateClass(classId, updateReq, instructorId, false);

        assertThat(updated.getLeadInstructorId()).isEqualTo(newLeadId);
        assertThat(updated.getWeeklyMinutes()).isEqualTo(1800);
        assertThat(lastQuotaPath.get()).isEqualTo("/v1/quota-accounts/class:" + classId);
        assertThat(lastQuotaBody.get())
                .contains("\"parentRef\":\"instructor:" + newLeadId + "\"")
                .contains("\"weeklyMinutes\":1800");
    }

    @Test
    @DisplayName("Tạo lớp học bị trùng mã lớp phải ném lỗi 409 CLASS_CODE_ALREADY_EXISTS")
    void createClass_duplicateCode_throwsConflict() {
        CreateClassRequest request = new CreateClassRequest("IT4060-20261", "Hệ điều hành Linux", "2026.1");
        when(classRepository.existsByCode("IT4060-20261")).thenReturn(true);

        assertThatThrownBy(() -> classService.createClass(request, instructorId, "GV", "gv@hust.edu.vn", "vi"))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> assertThat(((ProblemException) ex).getCode()).isEqualTo("CLASS_CODE_ALREADY_EXISTS"));

        verify(classRepository, never()).save(any());
    }

    @Test
    @DisplayName("Người ngoài lớp xem chi tiết hoặc xem thành viên bị từ chối 403 FORBIDDEN")
    void getClassDetailAndMembers_outsider_throwsForbidden() {
        UUID classId = UUID.randomUUID();
        UUID outsiderId = UUID.randomUUID();
        ClassEntity clazz = new ClassEntity("IT4060", "HĐH", "2026.1", instructorId, 1200, 30);
        clazz.setId(classId);

        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classRepository.existsById(classId)).thenReturn(true);
        when(classMemberRepository.findByClazzIdAndUserId(classId, outsiderId)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> classService.getClassDetail(classId, outsiderId, false))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> assertThat(((ProblemException) ex).getCode()).isEqualTo("FORBIDDEN"));

        assertThatThrownBy(() -> classService.getClassMembers(classId, outsiderId, null, false))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> assertThat(((ProblemException) ex).getCode()).isEqualTo("FORBIDDEN"));
    }

    @Test
    @DisplayName("Trợ giảng không có quyền thêm thành viên vào lớp")
    void addClassStaff_byAssistant_throwsForbidden() {
        UUID classId = UUID.randomUUID();
        UUID taUserId = UUID.randomUUID();
        ClassEntity clazz = new ClassEntity("IT4060", "HĐH", "2026.1", instructorId, 1200, 30);
        clazz.setId(classId);

        ClassMemberEntity taMember = new ClassMemberEntity(
                clazz, taUserId, dev.crucible.learning.domain.ClassRole.TEACHING_ASSISTANT, null, "TA", "ta@hust.edu.vn");

        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classMemberRepository.findByClazzIdAndUserId(classId, taUserId)).thenReturn(Optional.of(taMember));

        AddClassStaffRequest req = new AddClassStaffRequest("newta@hust.edu.vn", "TA 2", ClassRole.TEACHING_ASSISTANT);

        assertThatThrownBy(() -> classService.addClassStaff(classId, req, taUserId, false))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> assertThat(((ProblemException) ex).getCode()).isEqualTo("FORBIDDEN"));
    }

    @Test
    @DisplayName("Không cho phép xóa giảng viên phụ trách chính khỏi lớp")
    void removeClassMember_leadInstructor_throwsValidationFailed() {
        UUID classId = UUID.randomUUID();
        UUID memberId = UUID.randomUUID();
        ClassEntity clazz = new ClassEntity("IT4060", "HĐH", "2026.1", instructorId, 1200, 30);
        clazz.setId(classId);

        ClassMemberEntity leadMember = new ClassMemberEntity(
                clazz, instructorId, dev.crucible.learning.domain.ClassRole.INSTRUCTOR, null, "GV", "gv@hust.edu.vn");
        leadMember.setId(memberId);

        when(classRepository.findById(classId)).thenReturn(Optional.of(clazz));
        when(classMemberRepository.findByClazzIdAndUserId(classId, instructorId)).thenReturn(Optional.of(leadMember));
        when(classMemberRepository.findByIdAndClazzId(memberId, classId)).thenReturn(Optional.of(leadMember));

        assertThatThrownBy(() -> classService.removeClassMember(classId, memberId, instructorId, false))
                .isInstanceOf(ProblemException.class)
                .satisfies(ex -> assertThat(((ProblemException) ex).getCode()).isEqualTo("VALIDATION_FAILED"));
    }
}
