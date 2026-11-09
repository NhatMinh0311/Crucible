package dev.crucible.learning.api;

import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.util.LinkedHashMap;
import java.util.Map;

@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(ProblemException.class)
    public ResponseEntity<Map<String, Object>> handleProblem(ProblemException ex) {
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("type", "about:blank");
        body.put("title", ex.getStatus().getReasonPhrase());
        body.put("status", ex.getStatus().value());
        body.put("code", ex.getCode());
        body.put("detail", ex.getDetail());
        return ResponseEntity.status(ex.getStatus())
                .header(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_PROBLEM_JSON_VALUE)
                .body(body);
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Map<String, Object>> handleValidation(MethodArgumentNotValidException ex) {
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("type", "about:blank");
        body.put("title", "Bad Request");
        body.put("status", HttpStatus.BAD_REQUEST.value());
        body.put("code", "VALIDATION_FAILED");
        String message = ex.getBindingResult().getAllErrors().isEmpty() ? "Lỗi dữ liệu đầu vào" :
                ex.getBindingResult().getAllErrors().getFirst().getDefaultMessage();
        body.put("detail", message);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .header(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_PROBLEM_JSON_VALUE)
                .body(body);
    }
}
