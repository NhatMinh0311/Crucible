package dev.crucible.learning.infrastructure.security;

import dev.crucible.learning.api.ProblemException;
import dev.crucible.learning.domain.ClassMemberEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.stereotype.Component;

import java.util.Optional;
import java.util.UUID;

@Component
public class UserContext {

    public Optional<Jwt> getJwt() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth != null && auth.getPrincipal() instanceof Jwt jwt) {
            return Optional.of(jwt);
        }
        return Optional.empty();
    }

    public UUID requireUserId() {
        return getUserId().orElseThrow(() ->
                ProblemException.unauthenticated("Token không hợp lệ hoặc thiếu định danh người dùng"));
    }

    public Optional<UUID> getUserId() {
        return getJwt().map(jwt -> {
            String sub = jwt.getSubject();
            if (sub == null || sub.isBlank()) {
                return null;
            }
            try {
                return UUID.fromString(sub);
            } catch (IllegalArgumentException e) {
                return null;
            }
        });
    }

    public Optional<String> getEmail() {
        return getJwt().map(jwt -> {
            String email = jwt.getClaimAsString("email");
            return (email != null && !email.isBlank()) ? email.trim().toLowerCase() : null;
        });
    }

    public Optional<String> getFullName() {
        return getJwt().map(jwt -> {
            String name = jwt.getClaimAsString("name");
            if (name != null && !name.isBlank()) {
                return name.trim();
            }
            String preferred = jwt.getClaimAsString("preferred_username");
            return (preferred != null && !preferred.isBlank()) ? preferred.trim() : null;
        });
    }

    public Optional<String> getPreferredUsername() {
        return getFullName();
    }

    public String getLocale() {
        return getJwt()
                .map(jwt -> ClassMemberEntity.normalizeLocale(jwt.getClaimAsString("locale")))
                .orElse("vi");
    }

    public boolean hasRole(String role) {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth == null) {
            return false;
        }
        String expectedRole = role.startsWith("ROLE_") ? role : "ROLE_" + role;
        return auth.getAuthorities().stream()
                .map(GrantedAuthority::getAuthority)
                .anyMatch(a -> a.equals(expectedRole) || a.equals(role));
    }
}
