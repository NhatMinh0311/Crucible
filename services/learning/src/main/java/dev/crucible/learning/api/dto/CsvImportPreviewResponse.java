package dev.crucible.learning.api.dto;

import dev.crucible.learning.domain.CsvRowError;

import java.util.List;

public record CsvImportPreviewResponse(
    int totalRows,
    int validRows,
    int errorRows,
    List<CsvStudentRow> validStudents,
    List<CsvRowError> errors
) {
    public record CsvStudentRow(
        int line,
        String studentId,
        String fullName,
        String email
    ) {}
}
