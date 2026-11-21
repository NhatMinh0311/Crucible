package dev.crucible.learning.application;

import dev.crucible.learning.domain.ClassMemberEntity;
import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.Set;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class UserAutoEnrollService {

    private static final Logger log = LoggerFactory.getLogger(UserAutoEnrollService.class);
    private final ClassMemberRepository classMemberRepository;
    private final Set<String> enrolledSessionKeys = ConcurrentHashMap.newKeySet();

    public UserAutoEnrollService(ClassMemberRepository classMemberRepository) {
        this.classMemberRepository = classMemberRepository;
    }

    /**
     * Tự động liên kết tài khoản người dùng đăng nhập vào các lớp có email trùng khớp,
     * đồng thời cập nhật locale và joined_at (AC-L1.5, BR-M3.05).
     */
    @Transactional
    public int autoEnrollByEmail(String email, UUID userId, String locale) {
        if (email == null || email.isBlank() || userId == null) {
            return 0;
        }
        String normalizedEmail = email.trim().toLowerCase();
        String normalizedLocale = ClassMemberEntity.normalizeLocale(locale);
        String cacheKey = userId + "|" + normalizedEmail + "|" + normalizedLocale;
        if (!enrolledSessionKeys.add(cacheKey)) {
            return 0;
        }
        int linkedCount = classMemberRepository.linkUserIdByEmail(normalizedEmail, userId, normalizedLocale, Instant.now());
        if (linkedCount > 0) {
            log.info("Auto-enrolled user {} to {} class(es) by email {}", userId, linkedCount, normalizedEmail);
        }
        return linkedCount;
    }

    @Transactional
    public int autoEnrollByEmail(String email, UUID userId) {
        return autoEnrollByEmail(email, userId, "vi");
    }

    /**
     * Làm mới bộ nhớ đệm khi có thành viên mới được thêm vào lớp để người dùng đang đăng nhập vẫn được gắn ngay.
     */
    public void invalidateCache() {
        enrolledSessionKeys.clear();
    }
}
