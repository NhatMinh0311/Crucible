package dev.crucible.learning.infrastructure.persistence;

import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.domain.ClassRole;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.time.Instant;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface ClassMemberRepository extends JpaRepository<ClassMemberEntity, UUID> {

    List<ClassMemberEntity> findByClazzId(UUID classId);

    List<ClassMemberEntity> findByClazzIdAndRole(UUID classId, ClassRole role);

    Optional<ClassMemberEntity> findByClazzIdAndUserId(UUID classId, UUID userId);

    Optional<ClassMemberEntity> findByIdAndClazzId(UUID id, UUID classId);

    @Query("SELECT CASE WHEN COUNT(m) > 0 THEN true ELSE false END FROM ClassMemberEntity m WHERE m.clazz.id = :classId AND LOWER(m.email) = LOWER(:email)")
    boolean existsByClazzIdAndEmail(@Param("classId") UUID classId, @Param("email") String email);

    boolean existsByClazzIdAndStudentIdNumber(UUID classId, String studentIdNumber);

    @Modifying
    @Query("UPDATE ClassMemberEntity m SET m.userId = :userId, m.locale = :locale, m.joinedAt = :joinedAt WHERE LOWER(m.email) = LOWER(:email) AND m.userId IS NULL")
    int linkUserIdByEmail(@Param("email") String email, @Param("userId") UUID userId, @Param("locale") String locale, @Param("joinedAt") Instant joinedAt);
}
