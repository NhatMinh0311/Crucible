package dev.crucible.learning.application;

import dev.crucible.learning.infrastructure.persistence.ClassMemberRepository;
import dev.crucible.learning.infrastructure.security.UserAutoEnrollFilter;
import dev.crucible.learning.infrastructure.security.UserContext;
import jakarta.servlet.FilterChain;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.mock.web.MockHttpServletRequest;
import org.springframework.mock.web.MockHttpServletResponse;

import java.time.Instant;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UserAutoEnrollServiceTest {

    @Mock
    private ClassMemberRepository classMemberRepository;

    @InjectMocks
    private UserAutoEnrollService userAutoEnrollService;

    @Test
    @DisplayName("AC_L1_5 BR_M3_05: Tự động cập nhật user_id, locale và joined_at khi email đăng nhập trùng khớp")
    void autoEnrollByEmail_linksUserIdLocaleAndJoinedAt_AC_L1_5() {
        UUID userId = UUID.randomUUID();
        String rawEmail = "  Student01@SIS.HUST.EDU.VN ";
        when(classMemberRepository.linkUserIdByEmail(eq("student01@sis.hust.edu.vn"), eq(userId), eq("en"), any(Instant.class)))
                .thenReturn(2);

        int linked = userAutoEnrollService.autoEnrollByEmail(rawEmail, userId, "en-US");

        assertThat(linked).isEqualTo(2);
        verify(classMemberRepository).linkUserIdByEmail(
                eq("student01@sis.hust.edu.vn"),
                eq(userId),
                eq("en"),
                any(Instant.class)
        );
    }

    @Test
    @DisplayName("Bỏ qua tự động liên kết nếu email rỗng hoặc userId null")
    void autoEnrollByEmail_skipsWhenInvalid() {
        int linked = userAutoEnrollService.autoEnrollByEmail("  ", UUID.randomUUID(), "vi");
        assertThat(linked).isZero();
        verifyNoInteractions(classMemberRepository);
    }

    @Test
    @DisplayName("BR_M3_05: UserAutoEnrollFilter chạy trước mọi yêu cầu đã xác thực để sinh viên mở thẳng chi tiết lớp không bị 403")
    void userAutoEnrollFilter_invokesAutoEnrollBeforeController() throws Exception {
        UserContext userContext = mock(UserContext.class);
        UserAutoEnrollService mockEnrollService = mock(UserAutoEnrollService.class);
        FilterChain filterChain = mock(FilterChain.class);

        UUID userId = UUID.randomUUID();
        when(userContext.getUserId()).thenReturn(Optional.of(userId));
        when(userContext.getEmail()).thenReturn(Optional.of("sv@sis.hust.edu.vn"));
        when(userContext.getLocale()).thenReturn("vi");

        UserAutoEnrollFilter filter = new UserAutoEnrollFilter(userContext, mockEnrollService);
        MockHttpServletRequest request = new MockHttpServletRequest("GET", "/v1/classes/" + UUID.randomUUID());
        MockHttpServletResponse response = new MockHttpServletResponse();

        filter.doFilter(request, response, filterChain);

        verify(mockEnrollService).autoEnrollByEmail("sv@sis.hust.edu.vn", userId, "vi");
        verify(filterChain).doFilter(request, response);
    }
}
