package dev.crucible.learning.api;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.UUID;

@RestControllerAdvice
public class GlobalExceptionHandler {

    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(ProblemException.class)
    public ResponseEntity<Map<String, Object>> handleProblem(ProblemException ex) {
        return buildProblemResponse(ex.getStatus(), ex.getCode(), ex.getDetail());
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Map<String, Object>> handleValidation(MethodArgumentNotValidException ex) {
        String message = ex.getBindingResult().getAllErrors().isEmpty()
                ? "Lỗi dữ liệu đầu vào"
                : ex.getBindingResult().getAllErrors().getFirst().getDefaultMessage();
        return buildProblemResponse(HttpStatus.BAD_REQUEST, "VALIDATION_FAILED", message);
    }

    @ExceptionHandler(DataIntegrityViolationException.class)
    public ResponseEntity<Map<String, Object>> handleDataIntegrityViolation(DataIntegrityViolationException ex) {
        String msg = ex.getMostSpecificCause() != null ? ex.getMostSpecificCause().getMessage() : ex.getMessage();
        String code = "MEMBER_ALREADY_EXISTS";
        String detail = "Dữ liệu bị trùng lặp do thao tác đồng thời";
        if (msg != null && msg.contains("uk_classes_code")) {
            code = "CLASS_CODE_ALREADY_EXISTS";
            detail = "Mã lớp học đã tồn tại";
        }
        return buildProblemResponse(HttpStatus.CONFLICT, code, detail);
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<Map<String, Object>> handleUnexpected(Exception ex) {
        String traceId = currentTraceId();
        log.error("Unhandled exception [traceId={}]: {}", traceId, ex.getMessage(), ex);
        return buildProblemResponse(HttpStatus.INTERNAL_SERVER_ERROR, "INTERNAL", "Lỗi nội bộ", traceId);
    }

    private ResponseEntity<Map<String, Object>> buildProblemResponse(HttpStatus status, String code, String detail) {
        return buildProblemResponse(status, code, detail, currentTraceId());
    }

    private ResponseEntity<Map<String, Object>> buildProblemResponse(HttpStatus status, String code, String detail, String traceId) {
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("type", "about:blank");
        body.put("title", status.getReasonPhrase());
        body.put("status", status.value());
        body.put("code", code);
        body.put("detail", detail);
        body.put("traceId", traceId);
        return ResponseEntity.status(status)
                .header(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_PROBLEM_JSON_VALUE)
                .body(body);
    }

    private String currentTraceId() {
        String traceId = MDC.get("traceId");
        if (traceId == null || traceId.isBlank()) {
            traceId = MDC.get("trace_id");
        }
        if (traceId == null || traceId.isBlank()) {
            traceId = UUID.randomUUID().toString();
            MDC.put("traceId", traceId);
            MDC.put("trace_id", traceId);
        }
        return traceId;
    }
}
