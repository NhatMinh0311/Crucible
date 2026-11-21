package dev.crucible.learning.infrastructure.security;

import dev.crucible.learning.application.UserAutoEnrollService;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.MDC;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.Optional;
import java.util.UUID;

/**
 * Bộ lọc gắn trace_id vào MDC và tự động liên kết user_id, locale, joined_at cho thành viên lớp học
 * ngay khi người dùng đã xác thực gửi bất kỳ yêu cầu nào (BR-M3.05, AC-L1.5).
 */
@Component
public class UserAutoEnrollFilter extends OncePerRequestFilter {

    private final UserContext userContext;
    private final UserAutoEnrollService userAutoEnrollService;

    public UserAutoEnrollFilter(UserContext userContext, UserAutoEnrollService userAutoEnrollService) {
        this.userContext = userContext;
        this.userAutoEnrollService = userAutoEnrollService;
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain) throws ServletException, IOException {
        boolean putTrace = false;
        String traceId = MDC.get("traceId");
        if (traceId == null || traceId.isBlank()) {
            traceId = MDC.get("trace_id");
        }
        if (traceId == null || traceId.isBlank()) {
            String headerTrace = request.getHeader("X-Trace-Id");
            traceId = (headerTrace != null && !headerTrace.isBlank()) ? headerTrace.trim() : UUID.randomUUID().toString();
            MDC.put("traceId", traceId);
            MDC.put("trace_id", traceId);
            putTrace = true;
        }

        try {
            Optional<UUID> userIdOpt = userContext.getUserId();
            Optional<String> emailOpt = userContext.getEmail();
            if (userIdOpt.isPresent() && emailOpt.isPresent()) {
                userAutoEnrollService.autoEnrollByEmail(emailOpt.get(), userIdOpt.get(), userContext.getLocale());
            }
            filterChain.doFilter(request, response);
        } finally {
            if (putTrace) {
                MDC.remove("traceId");
                MDC.remove("trace_id");
            }
        }
    }
}
