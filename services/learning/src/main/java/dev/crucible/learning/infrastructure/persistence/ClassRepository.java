package dev.crucible.learning.infrastructure.persistence;

import dev.crucible.learning.domain.ClassEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface ClassRepository extends JpaRepository<ClassEntity, UUID> {

    Optional<ClassEntity> findByCode(String code);

    boolean existsByCode(String code);

    @Query("SELECT DISTINCT c FROM ClassEntity c JOIN c.members m WHERE m.userId = :userId")
    List<ClassEntity> findByMemberUserId(@Param("userId") UUID userId);
}
