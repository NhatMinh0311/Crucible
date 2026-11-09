package dev.crucible.learning.infrastructure.persistence;

import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.domain.ClassRole;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface ClassMemberRepository extends JpaRepository<ClassMemberEntity, UUID> {

    List<ClassMemberEntity> findByClazzId(UUID classId);

    List<ClassMemberEntity> findByClazzIdAndRole(UUID classId, ClassRole role);

    Optional<ClassMemberEntity> findByClazzIdAndUserId(UUID classId, UUID userId);

    Optional<ClassMemberEntity> findByClazzIdAndEmail(UUID classId, String email);

    Optional<ClassMemberEntity> findByClazzIdAndStudentIdNumber(UUID classId, String studentIdNumber);

    boolean existsByClazzIdAndEmail(UUID classId, String email);

    boolean existsByClazzIdAndStudentIdNumber(UUID classId, String studentIdNumber);

    List<ClassMemberEntity> findByEmailAndUserIdIsNull(String email);

    @Modifying
    @Query("UPDATE ClassMemberEntity m SET m.userId = :userId WHERE m.email = :email AND m.userId IS NULL")
    int linkUserIdByEmail(@Param("email") String email, @Param("userId") UUID userId);
}
