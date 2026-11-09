package dev.crucible.learning.api.dto;

import dev.crucible.learning.domain.ClassEntity;
import dev.crucible.learning.domain.ClassStatus;

import java.time.Instant;
import java.util.UUID;

public record ClassResponse(
    UUID id,
    String code,
    String name,
    String semester,
    String instructorQuotaAccountRef,
    int sessionLimitPerStudent,
    ClassStatus status,
    int memberCount,
    Instant createdAt,
    Instant updatedAt
) {
    public static ClassResponse fromEntity(ClassEntity entity) {
        return new ClassResponse(
            entity.getId(),
            entity.getCode(),
            entity.getName(),
            entity.getSemester(),
            entity.getInstructorQuotaAccountRef(),
            entity.getSessionLimitPerStudent(),
            entity.getStatus(),
            entity.getMembers() != null ? entity.getMembers().size() : 0,
            entity.getCreatedAt(),
            entity.getUpdatedAt()
        );
    }
}
