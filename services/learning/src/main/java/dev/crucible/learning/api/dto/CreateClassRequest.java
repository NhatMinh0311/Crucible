package dev.crucible.learning.api.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Positive;
import jakarta.validation.constraints.Size;

public record CreateClassRequest(
    @NotBlank @Size(max = 32) String code,
    @NotBlank @Size(max = 255) String name,
    @NotBlank @Size(max = 32) String semester,
    String instructorQuotaAccountRef,
    @Positive Integer sessionLimitPerStudent
) {}
