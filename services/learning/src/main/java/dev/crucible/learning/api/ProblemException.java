package dev.crucible.learning.api;

import org.springframework.http.HttpStatus;

public class ProblemException extends RuntimeException {

    private final HttpStatus status;
    private final String code;
    private final String detail;

    public ProblemException(HttpStatus status, String code, String detail) {
        super(code + ": " + detail);
        this.status = status;
        this.code = code;
        this.detail = detail;
    }

    public HttpStatus getStatus() {
        return status;
    }

    public String getCode() {
        return code;
    }

    public String getDetail() {
        return detail;
    }

    public static ProblemException notFound(String detail) {
        return new ProblemException(HttpStatus.NOT_FOUND, "NOT_FOUND", detail);
    }

    public static ProblemException forbidden(String detail) {
        return new ProblemException(HttpStatus.FORBIDDEN, "FORBIDDEN", detail);
    }

    public static ProblemException conflict(String code, String detail) {
        return new ProblemException(HttpStatus.CONFLICT, code, detail);
    }

    public static ProblemException badRequest(String code, String detail) {
        return new ProblemException(HttpStatus.BAD_REQUEST, code, detail);
    }

    public static ProblemException unprocessable(String code, String detail) {
        return new ProblemException(HttpStatus.UNPROCESSABLE_ENTITY, code, detail);
    }
}
