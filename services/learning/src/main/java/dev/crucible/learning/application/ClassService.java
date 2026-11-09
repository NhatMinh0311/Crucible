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
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

@Service
public class ClassService {

    private final ClassRepository classRepository;
    private final ClassMemberRepository classMemberRepository;
    private final UserAutoEnrollService userAutoEnrollService;

    public ClassService(ClassRepository classRepository,
                        ClassMemberRepository classMemberRepository,
                        UserAutoEnrollService userAutoEnrollService) {
        this.classRepository = classRepository;
        this.classMemberRepository = classMemberRepository;
        this.userAutoEnrollService = userAutoEnrollService;
    }

    @Transactional
    public ClassResponse createClass(CreateClassRequest request, UUID creatorUserId, String creatorFullName, String creatorEmail) {
        if (classRepository.existsByCode(request.code())) {
            throw ProblemException.conflict("CLASS_CODE_ALREADY_EXISTS", "Mã lớp " + request.code() + " đã tồn tại");
        }

        int sessionLimit = request.sessionLimitPerStudent() != null ? request.sessionLimitPerStudent() : 1;
        ClassEntity clazz = new ClassEntity(
                request.code(),
                request.name(),
                request.semester(),
                request.instructorQuotaAccountRef(),
                sessionLimit
        );

        ClassMemberEntity instructor = new ClassMemberEntity(
                clazz,
                creatorUserId,
                ClassRole.INSTRUCTOR,
                null,
                creatorFullName != null && !creatorFullName.isBlank() ? creatorFullName : "Giảng viên",
                creatorEmail != null && !creatorEmail.isBlank() ? creatorEmail : "instructor@domain.edu.vn"
        );
        clazz.addMember(instructor);

        ClassEntity saved = classRepository.save(clazz);
        return ClassResponse.fromEntity(saved);
    }

    @Transactional(readOnly = true)
    public ClassResponse getClassDetail(UUID classId, UUID requestingUserId, boolean isAdmin) {
        ClassEntity clazz = classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        if (!isAdmin && classMemberRepository.findByClazzIdAndUserId(classId, requestingUserId).isEmpty()) {
            throw ProblemException.forbidden("Bạn không phải thành viên của lớp học này");
        }

        return ClassResponse.fromEntity(clazz);
    }

    @Transactional
    public List<ClassResponse> getMyClasses(UUID userId, String email, boolean isAdmin) {
        if (email != null && userId != null) {
            userAutoEnrollService.autoEnrollByEmail(email, userId);
        }

        if (isAdmin) {
            return classRepository.findAll().stream().map(ClassResponse::fromEntity).toList();
        }

        return classRepository.findByMemberUserId(userId).stream()
                .map(ClassResponse::fromEntity)
                .toList();
    }

    @Transactional(readOnly = true)
    public List<ClassMemberResponse> getClassMembers(UUID classId, UUID requestingUserId, ClassRole roleFilter, boolean isAdmin) {
        if (!isAdmin && classMemberRepository.findByClazzIdAndUserId(classId, requestingUserId).isEmpty()) {
            throw ProblemException.forbidden("Bạn không có quyền xem thành viên lớp này");
        }

        List<ClassMemberEntity> members = roleFilter != null
                ? classMemberRepository.findByClazzIdAndRole(classId, roleFilter)
                : classMemberRepository.findByClazzId(classId);

        return members.stream().map(ClassMemberResponse::fromEntity).toList();
    }

    @Transactional
    public ClassMemberResponse addInstructorOrAssistant(UUID classId, AddInstructorRequest request, UUID requestingUserId, boolean isAdmin) {
        ClassEntity clazz = classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        if (!isAdmin) {
            ClassMemberEntity requester = classMemberRepository.findByClazzIdAndUserId(classId, requestingUserId)
                    .orElseThrow(() -> ProblemException.forbidden("Bạn không phải thành viên lớp"));
            if (requester.getRole() != ClassRole.INSTRUCTOR) {
                throw ProblemException.forbidden("Chỉ giảng viên phụ trách mới có quyền thêm trợ giảng");
            }
        }

        if (request.role() == ClassRole.STUDENT) {
            throw ProblemException.badRequest("INVALID_ROLE", "Chỉ được thêm vai trò INSTRUCTOR hoặc TEACHING_ASSISTANT qua API này");
        }

        String email = request.email().trim().toLowerCase();
        if (classMemberRepository.existsByClazzIdAndEmail(classId, email)) {
            throw ProblemException.conflict("MEMBER_ALREADY_EXISTS", "Thành viên với email " + email + " đã có trong lớp");
        }

        ClassMemberEntity newMember = new ClassMemberEntity(
                clazz,
                null,
                request.role(),
                null,
                request.fullName(),
                email
        );

        ClassMemberEntity saved = classMemberRepository.save(newMember);
        return ClassMemberResponse.fromEntity(saved);
    }
}
