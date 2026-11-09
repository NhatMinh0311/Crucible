package dev.crucible.learning.application;

import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
public class UserAutoEnrollService {

    private static final Logger log = LoggerFactory.getLogger(UserAutoEnrollService.class);
    private final ClassMemberRepository classMemberRepository;

    public UserAutoEnrollService(ClassMemberRepository classMemberRepository) {
        this.classMemberRepository = classMemberRepository;
    }

    /**
     * Tự động liên kết tài khoản sinh viên đăng nhập lần đầu vào các lớp đã có email trùng khớp (AC-L1.5).
     */
    @Transactional
    public int autoEnrollByEmail(String email, UUID userId) {
        if (email == null || email.isBlank() || userId == null) {
            return 0;
        }
        int linkedCount = classMemberRepository.linkUserIdByEmail(email.trim().toLowerCase(), userId);
        if (linkedCount > 0) {
            log.info("Auto-enrolled user {} to {} class(es) by email {}", userId, linkedCount, email);
        }
        return linkedCount;
    }
}
