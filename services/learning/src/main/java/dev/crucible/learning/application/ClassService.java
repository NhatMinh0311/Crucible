package dev.crucible.learning.application;

import dev.crucible.learning.api.ProblemException;
import dev.crucible.learning.api.generated.model.*;
import dev.crucible.learning.domain.ClassEntity;
import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.domain.ClassRole;
import dev.crucible.learning.domain.ClassStatus;
import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import dev.crucible.learning.infrastructure.persistence.ClassRepository;
import dev.crucible.learning.infrastructure.runtime.generated.QuotaApi;
import dev.crucible.learning.infrastructure.runtime.generated.model.PutQuotaAccountRequest;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.client.RestClientException;
import org.springframework.web.client.RestClientResponseException;

import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.List;
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@Service
public class ClassService {

    private static final int DEFAULT_WEEKLY_MINUTES = 1200;
    private static final int DEFAULT_MAX_CONCURRENT_SESSIONS = 30;
    private static final Pattern JSON_CODE_PATTERN = Pattern.compile("\"code\"\\s*:\\s*\"([^\"]+)\"");
    private static final Pattern JSON_DETAIL_PATTERN = Pattern.compile("\"detail\"\\s*:\\s*\"([^\"]+)\"");

    private final ClassRepository classRepository;
    private final ClassMemberRepository classMemberRepository;
    private final QuotaApi quotaApi;
    private final UserAutoEnrollService userAutoEnrollService;

    public ClassService(ClassRepository classRepository,
                        ClassMemberRepository classMemberRepository,
                        QuotaApi quotaApi) {
        this(classRepository, classMemberRepository, quotaApi, null);
    }

    @Autowired
    public ClassService(ClassRepository classRepository,
                        ClassMemberRepository classMemberRepository,
                        QuotaApi quotaApi,
                        UserAutoEnrollService userAutoEnrollService) {
        this.classRepository = classRepository;
        this.classMemberRepository = classMemberRepository;
        this.quotaApi = quotaApi;
        this.userAutoEnrollService = userAutoEnrollService;
    }

    @Transactional
    public ClassItem createClass(CreateClassRequest request,
                                 UUID creatorUserId,
                                 String creatorFullName,
                                 String creatorEmail,
                                 String creatorLocale) {
        if (creatorEmail == null || creatorEmail.isBlank()) {
            throw ProblemException.validationFailed("Tài khoản người tạo lớp chưa có địa chỉ email");
        }

        String code = request.getCode().trim();
        if (classRepository.existsByCode(code)) {
            throw ProblemException.conflict("CLASS_CODE_ALREADY_EXISTS", "Mã lớp " + code + " đã tồn tại");
        }

        int weeklyMinutes = request.getWeeklyMinutes() != null ? request.getWeeklyMinutes() : DEFAULT_WEEKLY_MINUTES;
        int maxConcurrentSessions = request.getMaxConcurrentSessions() != null
                ? request.getMaxConcurrentSessions()
                : DEFAULT_MAX_CONCURRENT_SESSIONS;

        ClassEntity clazz = new ClassEntity(
                code,
                request.getName().trim(),
                request.getSemester().trim(),
                creatorUserId,
                weeklyMinutes,
                maxConcurrentSessions
        );

        String displayName = (creatorFullName != null && !creatorFullName.isBlank()) ? creatorFullName.trim() : creatorEmail.trim();
        ClassMemberEntity instructor = new ClassMemberEntity(
                clazz,
                creatorUserId,
                ClassRole.INSTRUCTOR,
                null,
                displayName,
                creatorEmail,
                creatorLocale
        );
        clazz.addMember(instructor);

        ClassEntity saved = classRepository.save(clazz);
        syncClassQuotaAccount(saved);
        return toClassItem(saved);
    }

    @Transactional
    public ClassItem updateClass(UUID classId,
                                 UpdateClassRequest request,
                                 UUID requestingUserId,
                                 boolean isAdmin) {
        ClassEntity clazz = classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        requireInstructorPermission(classId, requestingUserId, isAdmin);

        if (request.getName() != null && !request.getName().isBlank()) {
            clazz.setName(request.getName().trim());
        }
        if (request.getSemester() != null && !request.getSemester().isBlank()) {
            clazz.setSemester(request.getSemester().trim());
        }
        if (request.getStatus() != null) {
            clazz.setStatus(ClassStatus.valueOf(request.getStatus().getValue()));
        }

        boolean quotaNeedSync = false;

        if (request.getLeadInstructorId() != null && !request.getLeadInstructorId().equals(clazz.getLeadInstructorId())) {
            if (!isAdmin && !clazz.getLeadInstructorId().equals(requestingUserId)) {
                throw ProblemException.forbidden("Chỉ giảng viên phụ trách chính hiện tại hoặc quản trị viên mới được đổi giảng viên phụ trách chính");
            }
            ClassMemberEntity targetMember = classMemberRepository.findByClazzIdAndUserId(classId, request.getLeadInstructorId())
                    .orElseThrow(() -> ProblemException.validationFailed("Giảng viên phụ trách chính mới phải là thành viên đã liên kết trong lớp"));
            if (targetMember.getRole() != ClassRole.INSTRUCTOR) {
                throw ProblemException.validationFailed("Người được chọn làm giảng viên phụ trách chính phải có vai trò INSTRUCTOR");
            }
            clazz.setLeadInstructorId(request.getLeadInstructorId());
            quotaNeedSync = true;
        }

        if (request.getWeeklyMinutes() != null && request.getWeeklyMinutes() != clazz.getWeeklyMinutes()) {
            clazz.setWeeklyMinutes(request.getWeeklyMinutes());
            quotaNeedSync = true;
        }
        if (request.getMaxConcurrentSessions() != null && request.getMaxConcurrentSessions() != clazz.getMaxConcurrentSessions()) {
            clazz.setMaxConcurrentSessions(request.getMaxConcurrentSessions());
            quotaNeedSync = true;
        }

        clazz.setUpdatedAt(Instant.now());
        ClassEntity saved = classRepository.save(clazz);

        if (quotaNeedSync) {
            syncClassQuotaAccount(saved);
        }

        return toClassItem(saved);
    }

    @Transactional(readOnly = true)
    public ClassItem getClassDetail(UUID classId, UUID requestingUserId, boolean isAdmin) {
        ClassEntity clazz = classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        if (!isAdmin && classMemberRepository.findByClazzIdAndUserId(classId, requestingUserId).isEmpty()) {
            throw ProblemException.forbidden("Bạn không phải thành viên của lớp học này");
        }

        return toClassItem(clazz);
    }

    @Transactional(readOnly = true)
    public List<ClassItem> getMyClasses(UUID userId, boolean isAdmin) {
        if (isAdmin) {
            return classRepository.findAll().stream().map(ClassService::toClassItem).toList();
        }
        return classRepository.findByMemberUserId(userId).stream()
                .map(ClassService::toClassItem)
                .toList();
    }

    @Transactional(readOnly = true)
    public List<ClassMember> getClassMembers(UUID classId, UUID requestingUserId, ClassRole roleFilter, boolean isAdmin) {
        if (!classRepository.existsById(classId)) {
            throw ProblemException.notFound("Không tìm thấy lớp học");
        }
        if (!isAdmin && classMemberRepository.findByClazzIdAndUserId(classId, requestingUserId).isEmpty()) {
            throw ProblemException.forbidden("Bạn không có quyền xem thành viên lớp này");
        }

        List<ClassMemberEntity> members = roleFilter != null
                ? classMemberRepository.findByClazzIdAndRole(classId, roleFilter)
                : classMemberRepository.findByClazzId(classId);

        return members.stream().map(ClassService::toClassMember).toList();
    }

    @Transactional
    public ClassMember addClassStaff(UUID classId, AddClassStaffRequest request, UUID requestingUserId, boolean isAdmin) {
        ClassEntity clazz = classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        requireInstructorPermission(classId, requestingUserId, isAdmin);

        if (request.getRole() == dev.crucible.learning.api.generated.model.ClassRole.STUDENT) {
            throw ProblemException.validationFailed("Chỉ được thêm vai trò INSTRUCTOR hoặc TEACHING_ASSISTANT qua API này");
        }

        String email = request.getEmail().trim().toLowerCase();
        if (classMemberRepository.existsByClazzIdAndEmail(classId, email)) {
            throw ProblemException.conflict("MEMBER_ALREADY_EXISTS", "Thành viên với email " + email + " đã có trong lớp");
        }

        ClassMemberEntity newMember = new ClassMemberEntity(
                clazz,
                null,
                ClassRole.valueOf(request.getRole().getValue()),
                null,
                request.getFullName().trim(),
                email
        );

        ClassMemberEntity saved = classMemberRepository.save(newMember);
        if (userAutoEnrollService != null) {
            userAutoEnrollService.invalidateCache();
        }
        return toClassMember(saved);
    }

    @Transactional
    public void removeClassMember(UUID classId, UUID memberId, UUID requestingUserId, boolean isAdmin) {
        ClassEntity clazz = classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        requireInstructorPermission(classId, requestingUserId, isAdmin);

        ClassMemberEntity member = classMemberRepository.findByIdAndClazzId(memberId, classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy thành viên trong lớp"));

        if (member.getUserId() != null && member.getUserId().equals(clazz.getLeadInstructorId())) {
            throw ProblemException.validationFailed("Không thể xóa giảng viên phụ trách chính của lớp");
        }

        classMemberRepository.delete(member);
    }

    private void requireInstructorPermission(UUID classId, UUID requestingUserId, boolean isAdmin) {
        if (!isAdmin) {
            ClassMemberEntity requester = classMemberRepository.findByClazzIdAndUserId(classId, requestingUserId)
                    .orElseThrow(() -> ProblemException.forbidden("Bạn không phải thành viên lớp"));
            if (requester.getRole() != ClassRole.INSTRUCTOR) {
                throw ProblemException.forbidden("Chỉ giảng viên phụ trách mới có quyền thực hiện thao tác này");
            }
        }
    }

    private void syncClassQuotaAccount(ClassEntity clazz) {
        PutQuotaAccountRequest quotaReq = new PutQuotaAccountRequest(
                clazz.getWeeklyMinutes(),
                clazz.getMaxConcurrentSessions()
        ).parentRef("instructor:" + clazz.getLeadInstructorId());

        try {
            quotaApi.putQuotaAccount(clazz.getQuotaAccountRef(), quotaReq);
        } catch (RestClientResponseException ex) {
            int statusCode = ex.getStatusCode().value();
            HttpStatus status = HttpStatus.resolve(statusCode);
            if (status == null) {
                status = HttpStatus.INTERNAL_SERVER_ERROR;
            }
            String body = ex.getResponseBodyAsString();
            String code = extractJsonField(JSON_CODE_PATTERN, body);
            String detail = extractJsonField(JSON_DETAIL_PATTERN, body);
            if (code == null || code.isBlank()) {
                code = switch (status) {
                    case UNPROCESSABLE_CONTENT -> "QUOTA_CEILING_EXCEEDED";
                    case NOT_IMPLEMENTED -> "NOT_IMPLEMENTED";
                    case FORBIDDEN -> "FORBIDDEN";
                    case UNAUTHORIZED -> "UNAUTHENTICATED";
                    case BAD_REQUEST -> "VALIDATION_FAILED";
                    default -> "INTERNAL";
                };
            }
            if (detail == null || detail.isBlank()) {
                detail = "Lỗi đồng bộ tài khoản hạn mức của lớp bên M1 (" + statusCode + ")";
            }
            throw new ProblemException(status, code, detail);
        } catch (RestClientException ex) {
            throw new ProblemException(
                    HttpStatus.INTERNAL_SERVER_ERROR,
                    "INTERNAL",
                    "Không kết nối được tới dịch vụ runtime-api: " + ex.getMessage()
            );
        }
    }

    private static String extractJsonField(Pattern pattern, String json) {
        if (json == null || json.isBlank()) {
            return null;
        }
        Matcher matcher = pattern.matcher(json);
        return matcher.find() ? matcher.group(1) : null;
    }

    public static ClassItem toClassItem(ClassEntity entity) {
        return new ClassItem(
                entity.getId(),
                entity.getCode(),
                entity.getName(),
                entity.getSemester(),
                entity.getLeadInstructorId(),
                entity.getQuotaAccountRef(),
                entity.getWeeklyMinutes(),
                entity.getMaxConcurrentSessions(),
                dev.crucible.learning.api.generated.model.ClassStatus.fromValue(entity.getStatus().name()),
                entity.getMembers() != null ? entity.getMembers().size() : 0,
                toOffsetDateTime(entity.getCreatedAt()),
                toOffsetDateTime(entity.getUpdatedAt())
        );
    }

    public static ClassMember toClassMember(ClassMemberEntity entity) {
        ClassMember member = new ClassMember(
                entity.getId(),
                entity.getClazz().getId(),
                dev.crucible.learning.api.generated.model.ClassRole.fromValue(entity.getRole().name()),
                entity.getFullName(),
                entity.getEmail(),
                ClassMember.LocaleEnum.fromValue(ClassMemberEntity.normalizeLocale(entity.getLocale())),
                toOffsetDateTime(entity.getCreatedAt())
        );
        member.setUserId(entity.getUserId());
        member.setStudentIdNumber(entity.getStudentIdNumber());
        member.setJoinedAt(toOffsetDateTime(entity.getJoinedAt()));
        return member;
    }

    private static OffsetDateTime toOffsetDateTime(Instant instant) {
        return instant != null ? instant.atOffset(ZoneOffset.UTC) : null;
    }
}
