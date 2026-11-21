package dev.crucible.learning.api;

import dev.crucible.learning.api.generated.ClassesApi;
import dev.crucible.learning.api.generated.model.*;
import dev.crucible.learning.application.ClassService;
import dev.crucible.learning.application.CsvImportService;
import dev.crucible.learning.infrastructure.security.UserContext;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.List;
import java.util.UUID;

@RestController
public class ClassController implements ClassesApi {

    private final ClassService classService;
    private final CsvImportService csvImportService;
    private final UserContext userContext;

    public ClassController(ClassService classService, CsvImportService csvImportService, UserContext userContext) {
        this.classService = classService;
        this.csvImportService = csvImportService;
        this.userContext = userContext;
    }

    @Override
    public ResponseEntity<List<ClassItem>> listClasses() {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");
        return ResponseEntity.ok(classService.getMyClasses(userId, isAdmin));
    }

    @Override
    public ResponseEntity<ClassItem> createClass(CreateClassRequest createClassRequest) {
        UUID userId = userContext.requireUserId();
        if (!userContext.hasRole("INSTRUCTOR") && !userContext.hasRole("ADMIN")) {
            throw ProblemException.forbidden("Chỉ giảng viên hoặc quản trị viên mới có quyền tạo lớp học");
        }
        String email = userContext.getEmail().orElse(null);
        String fullName = userContext.getFullName().orElse(null);
        String locale = userContext.getLocale();

        ClassItem created = classService.createClass(createClassRequest, userId, fullName, email, locale);
        return ResponseEntity.status(HttpStatus.CREATED).body(created);
    }

    @Override
    public ResponseEntity<ClassItem> getClass(UUID classId) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");
        return ResponseEntity.ok(classService.getClassDetail(classId, userId, isAdmin));
    }

    @Override
    public ResponseEntity<ClassItem> updateClass(UUID classId, UpdateClassRequest updateClassRequest) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");
        return ResponseEntity.ok(classService.updateClass(classId, updateClassRequest, userId, isAdmin));
    }

    @Override
    public ResponseEntity<List<ClassMember>> listClassMembers(UUID classId, ClassRole role) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");
        dev.crucible.learning.domain.ClassRole domainRole = role != null
                ? dev.crucible.learning.domain.ClassRole.valueOf(role.getValue())
                : null;
        return ResponseEntity.ok(classService.getClassMembers(classId, userId, domainRole, isAdmin));
    }

    @Override
    public ResponseEntity<CsvImportPreviewResponse> previewCsvMembers(UUID classId, MultipartFile file) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");

        if (file == null || file.isEmpty()) {
            throw ProblemException.validationFailed("Tệp CSV tải lên không có dữ liệu");
        }
        try {
            CsvImportPreviewResponse preview = csvImportService.previewImport(classId, file.getInputStream(), userId, isAdmin);
            return ResponseEntity.ok(preview);
        } catch (IOException e) {
            throw ProblemException.validationFailed("Không thể đọc tệp CSV: " + e.getMessage());
        }
    }

    @Override
    public ResponseEntity<CsvImportCommitResponse> importCsvMembers(UUID classId, MultipartFile file) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");

        if (file == null || file.isEmpty()) {
            throw ProblemException.validationFailed("Tệp CSV tải lên không có dữ liệu");
        }
        try {
            CsvImportCommitResponse committed = csvImportService.commitImport(classId, file.getInputStream(), userId, isAdmin);
            return ResponseEntity.ok(committed);
        } catch (IOException e) {
            throw ProblemException.validationFailed("Không thể đọc tệp CSV: " + e.getMessage());
        }
    }

    @Override
    public ResponseEntity<ClassMember> addClassStaff(UUID classId, AddClassStaffRequest addClassStaffRequest) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");
        ClassMember member = classService.addClassStaff(classId, addClassStaffRequest, userId, isAdmin);
        return ResponseEntity.status(HttpStatus.CREATED).body(member);
    }

    @Override
    public ResponseEntity<Void> removeClassMember(UUID classId, UUID memberId) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");
        classService.removeClassMember(classId, memberId, userId, isAdmin);
        return ResponseEntity.noContent().build();
    }
}
