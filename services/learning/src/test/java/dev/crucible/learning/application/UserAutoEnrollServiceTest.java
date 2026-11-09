package dev.crucible.learning.application;

import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class UserAutoEnrollServiceTest {

    @Mock
    private ClassMemberRepository classMemberRepository;

    @InjectMocks
    private UserAutoEnrollService userAutoEnrollService;

    @Test
    @DisplayName("AC-L1.5: Tự động cập nhật user_id khi email người dùng trùng với danh sách lớp")
    void autoEnroll_shouldLinkUserId_whenEmailMatches_AC_L1_5() {
        UUID userId = UUID.randomUUID();
        String email = "Student1@Sis.Hust.Edu.Vn";

        when(classMemberRepository.linkUserIdByEmail("student1@sis.hust.edu.vn", userId)).thenReturn(2);

        int count = userAutoEnrollService.autoEnrollByEmail(email, userId);

        assertThat(count).isEqualTo(2);
        verify(classMemberRepository).linkUserIdByEmail("student1@sis.hust.edu.vn", userId);
    }

    @Test
    @DisplayName("AC-L1.5: Bỏ qua khi email rỗng hoặc userId null")
    void autoEnroll_shouldReturnZero_whenInputInvalid() {
        int c1 = userAutoEnrollService.autoEnrollByEmail(null, UUID.randomUUID());
        int c2 = userAutoEnrollService.autoEnrollByEmail("a@b.com", null);

        assertThat(c1).isEqualTo(0);
        assertThat(c2).isEqualTo(0);
    }
}
