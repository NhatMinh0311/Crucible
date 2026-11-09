package dev.crucible.learning.infrastructure.security;

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
        return getUserId().orElseThrow(() -> new IllegalStateException("Unauthenticated user context"));
    }

    public Optional<UUID> getUserId() {
        return getJwt().map(jwt -> {
            try {
                return UUID.fromString(jwt.getSubject());
            } catch (IllegalArgumentException e) {
                return null;
            }
        });
    }

    public Optional<String> getEmail() {
        return getJwt().map(jwt -> jwt.getClaimAsString("email"));
    }

    public Optional<String> getPreferredUsername() {
        return getJwt().map(jwt -> jwt.getClaimAsString("preferred_username"));
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
