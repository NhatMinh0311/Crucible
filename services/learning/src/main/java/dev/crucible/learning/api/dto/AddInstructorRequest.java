package dev.crucible.learning.api.dto;

import dev.crucible.learning.domain.ClassRole;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

public record AddInstructorRequest(
    @NotBlank @Email String email,
    @NotBlank String fullName,
    @NotNull ClassRole role
) {}
