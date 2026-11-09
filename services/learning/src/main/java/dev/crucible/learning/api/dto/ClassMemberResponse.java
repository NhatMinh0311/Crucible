package dev.crucible.learning.api.dto;

import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.domain.ClassRole;

import java.time.Instant;
import java.util.UUID;

public record ClassMemberResponse(
    UUID id,
    UUID classId,
    UUID userId,
    ClassRole role,
    String studentIdNumber,
    String fullName,
    String email,
    Instant createdAt
) {
    public static ClassMemberResponse fromEntity(ClassMemberEntity entity) {
        return new ClassMemberResponse(
            entity.getId(),
            entity.getClazz().getId(),
            entity.getUserId(),
            entity.getRole(),
            entity.getStudentIdNumber(),
            entity.getFullName(),
            entity.getEmail(),
            entity.getCreatedAt()
        );
    }
}
