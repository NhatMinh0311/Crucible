package dev.crucible.learning.infrastructure.security;

import dev.crucible.platform.security.KeycloakJwtConverter;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.MDC;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.MediaType;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.oauth2.client.AuthorizedClientServiceOAuth2AuthorizedClientManager;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientManager;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.registration.ClientRegistrationRepository;
import org.springframework.security.oauth2.server.resource.web.authentication.BearerTokenAuthenticationFilter;
import org.springframework.security.web.SecurityFilterChain;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.UUID;

/**
 * Xác thực token Keycloak cho API của dịch vụ và quản lý token client credentials khi gọi module khác (ADR-0004).
 * Quyền theo đối tượng được kiểm tra ở tầng application, không ở đây.
 */
@Configuration
class SecurityConfig {

    @Bean
    SecurityFilterChain api(HttpSecurity http, UserAutoEnrollFilter userAutoEnrollFilter) throws Exception {
        http.csrf(csrf -> csrf.disable())
                .sessionManagement(s -> s.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/actuator/health/**", "/actuator/info").permitAll()
                        .anyRequest().authenticated())
                .oauth2ResourceServer(o -> o
                        .jwt(jwt -> jwt.jwtAuthenticationConverter(new KeycloakJwtConverter()))
                        .authenticationEntryPoint((req, res, ex) ->
                                writeProblem(res, HttpServletResponse.SC_UNAUTHORIZED, "Unauthorized", "UNAUTHENTICATED", "Thiếu hoặc sai token"))
                        .accessDeniedHandler((req, res, ex) ->
                                writeProblem(res, HttpServletResponse.SC_FORBIDDEN, "Forbidden", "FORBIDDEN", "Token thiếu vai trò hoặc scope cần thiết")))
                .exceptionHandling(e -> e
                        .authenticationEntryPoint((req, res, ex) ->
                                writeProblem(res, HttpServletResponse.SC_UNAUTHORIZED, "Unauthorized", "UNAUTHENTICATED", "Thiếu hoặc sai token"))
                        .accessDeniedHandler((req, res, ex) ->
                                writeProblem(res, HttpServletResponse.SC_FORBIDDEN, "Forbidden", "FORBIDDEN", "Token thiếu vai trò hoặc scope cần thiết")))
                .addFilterAfter(userAutoEnrollFilter, BearerTokenAuthenticationFilter.class);
        return http.build();
    }

    /** Lấy và làm mới token client credentials ngoài ngữ cảnh yêu cầu HTTP (tác vụ nền, listener sự kiện). */
    @Bean
    OAuth2AuthorizedClientManager authorizedClientManager(ClientRegistrationRepository registrations,
                                                          OAuth2AuthorizedClientService clients) {
        return new AuthorizedClientServiceOAuth2AuthorizedClientManager(registrations, clients);
    }

    private static void writeProblem(HttpServletResponse response, int status, String title, String code, String detail) throws IOException {
        String traceId = MDC.get("traceId");
        if (traceId == null || traceId.isBlank()) {
            traceId = MDC.get("trace_id");
        }
        if (traceId == null || traceId.isBlank()) {
            traceId = UUID.randomUUID().toString();
        }
        response.setStatus(status);
        response.setContentType(MediaType.APPLICATION_PROBLEM_JSON_VALUE);
        response.setCharacterEncoding(StandardCharsets.UTF_8.name());
        String json = String.format(
                "{\"type\":\"about:blank\",\"title\":\"%s\",\"status\":%d,\"code\":\"%s\",\"detail\":\"%s\",\"traceId\":\"%s\"}",
                title, status, code, detail, traceId
        );
        response.getWriter().write(json);
    }
}
