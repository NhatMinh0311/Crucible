package dev.crucible.learning.application;

import dev.crucible.learning.api.ProblemException;
import dev.crucible.learning.api.generated.model.CsvImportCommitResponse;
import dev.crucible.learning.api.generated.model.CsvImportPreviewResponse;
import dev.crucible.learning.api.generated.model.CsvRowError;
import dev.crucible.learning.api.generated.model.CsvStudentRow;
import dev.crucible.learning.domain.ClassEntity;
import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.domain.ClassRole;
import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import dev.crucible.learning.infrastructure.persistence.ClassRepository;
import org.apache.commons.csv.CSVFormat;
import org.apache.commons.csv.CSVParser;
import org.apache.commons.csv.CSVRecord;
import org.apache.commons.io.input.BOMInputStream;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.regex.Pattern;

@Service
public class CsvImportService {

    private static final Pattern EMAIL_PATTERN = Pattern.compile("^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$");

    private final ClassRepository classRepository;
    private final ClassMemberRepository classMemberRepository;
    private final UserAutoEnrollService userAutoEnrollService;

    public CsvImportService(ClassRepository classRepository, ClassMemberRepository classMemberRepository) {
        this(classRepository, classMemberRepository, null);
    }

    @Autowired
    public CsvImportService(ClassRepository classRepository,
                            ClassMemberRepository classMemberRepository,
                            UserAutoEnrollService userAutoEnrollService) {
        this.classRepository = classRepository;
        this.classMemberRepository = classMemberRepository;
        this.userAutoEnrollService = userAutoEnrollService;
    }

    @Transactional(readOnly = true)
    public CsvImportPreviewResponse previewImport(UUID classId, InputStream csvInputStream, UUID requestingUserId, boolean isAdmin) {
        ClassEntity clazz = classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        requireInstructorPermission(clazz.getId(), requestingUserId, isAdmin);
        return parseAndValidate(classId, csvInputStream);
    }

    public CsvImportPreviewResponse parseAndValidate(UUID classId, InputStream csvInputStream) {
        classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        List<CsvRowError> errors = new ArrayList<>();
        List<CsvStudentRow> validRows = new ArrayList<>();

        Set<String> seenStudentIds = new HashSet<>();
        Set<String> seenEmails = new HashSet<>();

        Set<String> existingEmails = new HashSet<>();
        Set<String> existingStudentIds = new HashSet<>();
        for (ClassMemberEntity m : classMemberRepository.findByClazzId(classId)) {
            if (m.getEmail() != null) existingEmails.add(m.getEmail().toLowerCase());
            if (m.getStudentIdNumber() != null) existingStudentIds.add(m.getStudentIdNumber().toLowerCase());
        }

        try (BOMInputStream bomStream = BOMInputStream.builder().setInputStream(csvInputStream).get();
             CSVParser parser = CSVFormat.DEFAULT.builder()
                     .setHeader()
                     .setSkipHeaderRecord(true)
                     .setIgnoreHeaderCase(true)
                     .setTrim(true)
                     .setIgnoreEmptyLines(true)
                     .build()
                     .parse(new InputStreamReader(bomStream, StandardCharsets.UTF_8))) {

            Map<String, Integer> headerMap = parser.getHeaderMap();
            if (headerMap == null || headerMap.isEmpty()) {
                throw ProblemException.validationFailed("Tệp CSV không chứa tiêu đề");
            }

            String studentIdHeader = findHeader(headerMap, "student_id", "mssv", "studentId");
            String fullNameHeader = findHeader(headerMap, "full_name", "hoten", "name", "fullName");
            String emailHeader = findHeader(headerMap, "email", "mail");

            if (studentIdHeader == null || fullNameHeader == null || emailHeader == null) {
                throw ProblemException.validationFailed(
                        "Tiêu đề CSV không đúng định dạng. Cần chứa các cột: student_id, full_name, email");
            }

            int lineNum = 1;
            for (CSVRecord record : parser) {
                lineNum++;
                String studentId = record.isSet(studentIdHeader) ? record.get(studentIdHeader).trim() : "";
                String fullName = record.isSet(fullNameHeader) ? record.get(fullNameHeader).trim() : "";
                String email = record.isSet(emailHeader) ? record.get(emailHeader).trim().toLowerCase() : "";

                boolean lineHasError = false;

                if (studentId.isEmpty()) {
                    errors.add(new CsvRowError(lineNum, "student_id", "Thiếu mã số sinh viên"));
                    lineHasError = true;
                }
                if (fullName.isEmpty()) {
                    errors.add(new CsvRowError(lineNum, "full_name", "Thiếu họ và tên"));
                    lineHasError = true;
                }
                if (email.isEmpty()) {
                    errors.add(new CsvRowError(lineNum, "email", "Thiếu địa chỉ email"));
                    lineHasError = true;
                } else if (!EMAIL_PATTERN.matcher(email).matches()) {
                    errors.add(new CsvRowError(lineNum, "email", "Email không đúng định dạng: " + email));
                    lineHasError = true;
                }

                if (!studentId.isEmpty()) {
                    if (seenStudentIds.contains(studentId.toLowerCase())) {
                        errors.add(new CsvRowError(lineNum, "student_id", "Trùng mã sinh viên trong tệp CSV: " + studentId));
                        lineHasError = true;
                    } else if (existingStudentIds.contains(studentId.toLowerCase())) {
                        errors.add(new CsvRowError(lineNum, "student_id", "Mã sinh viên đã có trong lớp: " + studentId));
                        lineHasError = true;
                    } else {
                        seenStudentIds.add(studentId.toLowerCase());
                    }
                }

                if (!email.isEmpty() && EMAIL_PATTERN.matcher(email).matches()) {
                    if (seenEmails.contains(email)) {
                        errors.add(new CsvRowError(lineNum, "email", "Trùng email trong tệp CSV: " + email));
                        lineHasError = true;
                    } else if (existingEmails.contains(email)) {
                        errors.add(new CsvRowError(lineNum, "email", "Email đã có trong lớp: " + email));
                        lineHasError = true;
                    } else {
                        seenEmails.add(email);
                    }
                }

                if (!lineHasError) {
                    validRows.add(new CsvStudentRow(lineNum, studentId, fullName, email));
                }
            }
        } catch (ProblemException pe) {
            throw pe;
        } catch (Exception e) {
            throw ProblemException.validationFailed("Lỗi đọc tệp CSV: " + e.getMessage());
        }

        int errorRowCount = (int) errors.stream().map(CsvRowError::getLine).distinct().count();
        int total = validRows.size() + errorRowCount;

        return new CsvImportPreviewResponse(total, validRows.size(), errorRowCount, validRows, errors);
    }

    @Transactional
    public CsvImportCommitResponse commitImport(UUID classId, InputStream csvInputStream, UUID requestingUserId, boolean isAdmin) {
        ClassEntity clazz = classRepository.findById(classId)
                .orElseThrow(() -> ProblemException.notFound("Không tìm thấy lớp học"));

        requireInstructorPermission(classId, requestingUserId, isAdmin);

        CsvImportPreviewResponse preview = parseAndValidate(classId, csvInputStream);
        if (!preview.getErrors().isEmpty()) {
            throw ProblemException.validationFailed(
                    "Tệp CSV có " + preview.getErrorRows() + " dòng lỗi, không thể nhập danh sách");
        }

        List<ClassMemberEntity> newMembers = preview.getValidStudents().stream()
                .map(row -> new ClassMemberEntity(
                        clazz,
                        null,
                        ClassRole.STUDENT,
                        row.getStudentId(),
                        row.getFullName(),
                        row.getEmail()
                ))
                .toList();

        classMemberRepository.saveAll(newMembers);
        if (userAutoEnrollService != null) {
            userAutoEnrollService.invalidateCache();
        }
        return new CsvImportCommitResponse(newMembers.size());
    }

    private void requireInstructorPermission(UUID classId, UUID requestingUserId, boolean isAdmin) {
        if (!isAdmin) {
            ClassMemberEntity requester = classMemberRepository.findByClazzIdAndUserId(classId, requestingUserId)
                    .orElseThrow(() -> ProblemException.forbidden("Bạn không phải thành viên lớp"));
            if (requester.getRole() != ClassRole.INSTRUCTOR) {
                throw ProblemException.forbidden("Chỉ giảng viên phụ trách mới có quyền nhập danh sách sinh viên");
            }
        }
    }

    private String findHeader(Map<String, Integer> headerMap, String... candidates) {
        for (String c : candidates) {
            for (String h : headerMap.keySet()) {
                if (h.equalsIgnoreCase(c)) {
                    return h;
                }
            }
        }
        return null;
    }
}
