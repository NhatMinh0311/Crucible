package dev.crucible.learning.api;

import dev.crucible.learning.api.dto.*;
import dev.crucible.learning.application.ClassService;
import dev.crucible.learning.application.CsvImportService;
import dev.crucible.learning.domain.ClassRole;
import dev.crucible.learning.infrastructure.security.UserContext;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/classes")
public class ClassController {

    private final ClassService classService;
    private final CsvImportService csvImportService;
    private final UserContext userContext;

    public ClassController(ClassService classService, CsvImportService csvImportService, UserContext userContext) {
        this.classService = classService;
        this.csvImportService = csvImportService;
        this.userContext = userContext;
    }

    @PostMapping
    public ResponseEntity<ClassResponse> createClass(@Valid @RequestBody CreateClassRequest request) {
        UUID userId = userContext.requireUserId();
        String email = userContext.getEmail().orElse(null);
        String name = userContext.getPreferredUsername().orElse(null);

        ClassResponse response = classService.createClass(request, userId, name, email);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @GetMapping
    public ResponseEntity<List<ClassResponse>> getMyClasses() {
        UUID userId = userContext.requireUserId();
        String email = userContext.getEmail().orElse(null);
        boolean isAdmin = userContext.hasRole("ADMIN");

        List<ClassResponse> classes = classService.getMyClasses(userId, email, isAdmin);
        return ResponseEntity.ok(classes);
    }

    @GetMapping("/{classId}")
    public ResponseEntity<ClassResponse> getClassDetail(@PathVariable UUID classId) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");

        ClassResponse detail = classService.getClassDetail(classId, userId, isAdmin);
        return ResponseEntity.ok(detail);
    }

    @GetMapping("/{classId}/members")
    public ResponseEntity<List<ClassMemberResponse>> getClassMembers(
            @PathVariable UUID classId,
            @RequestParam(required = false) ClassRole role) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");

        List<ClassMemberResponse> members = classService.getClassMembers(classId, userId, role, isAdmin);
        return ResponseEntity.ok(members);
    }

    @PostMapping("/{classId}/instructors")
    public ResponseEntity<ClassMemberResponse> addInstructorOrAssistant(
            @PathVariable UUID classId,
            @Valid @RequestBody AddInstructorRequest request) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");

        ClassMemberResponse member = classService.addInstructorOrAssistant(classId, request, userId, isAdmin);
        return ResponseEntity.status(HttpStatus.CREATED).body(member);
    }

    @PostMapping(value = "/{classId}/members/import-csv", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public ResponseEntity<?> importCsv(
            @PathVariable UUID classId,
            @RequestParam("file") MultipartFile file,
            @RequestParam(value = "commit", defaultValue = "false") boolean commit) {
        UUID userId = userContext.requireUserId();
        boolean isAdmin = userContext.hasRole("ADMIN");

        if (file.isEmpty()) {
            throw ProblemException.badRequest("CSV_EMPTY", "Tệp CSV tải lên không có dữ liệu");
        }

        try {
            if (commit) {
                int importedCount = csvImportService.commitImport(classId, file.getInputStream(), userId, isAdmin);
                return ResponseEntity.ok(Map.of(
                        "success", true,
                        "importedCount", importedCount,
                        "message", "Đã nhập thành công " + importedCount + " sinh viên vào lớp"
                ));
            } else {
                CsvImportPreviewResponse preview = csvImportService.parseAndValidate(classId, file.getInputStream());
                return ResponseEntity.ok(preview);
            }
        } catch (IOException e) {
            throw ProblemException.badRequest("FILE_READ_ERROR", "Không thể đọc tệp CSV: " + e.getMessage());
        }
    }
}
