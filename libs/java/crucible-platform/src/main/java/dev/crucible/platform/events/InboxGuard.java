package dev.crucible.platform.events;

import java.util.UUID;
import org.springframework.jdbc.core.simple.JdbcClient;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

/**
 * Bảo đảm mỗi sự kiện được xử lý đúng một lần và sự kiện cũ hơn bản đã xử lý bị bỏ qua (CR-DOC-03 mục 5.2).
 * Gọi trong cùng transaction với xử lý nghiệp vụ của sự kiện.
 */
public class InboxGuard {

    private final JdbcClient jdbc;

    public InboxGuard(JdbcClient jdbc) {
        this.jdbc = jdbc;
    }

    /**
     * Ghi nhận sự kiện và cho biết có cần xử lý không.
     *
     * @return {@code false} nếu sự kiện đã được xử lý hoặc có {@code sequence} không mới hơn bản đã xử lý
     */
    @Transactional(propagation = Propagation.MANDATORY)
    public boolean accept(EventEnvelope<?> event) {
        int inserted = jdbc.sql("INSERT INTO crucible_inbox (event_id) VALUES (:id) ON CONFLICT DO NOTHING")
                .param("id", event.id())
                .update();
        if (inserted == 0) {
            return false;
        }
        int advanced = jdbc.sql("""
                        INSERT INTO crucible_inbox_sequence (subject, last_sequence) VALUES (:subject, :seq)
                        ON CONFLICT (subject) DO UPDATE SET last_sequence = EXCLUDED.last_sequence
                            WHERE crucible_inbox_sequence.last_sequence < EXCLUDED.last_sequence
                        """)
                .param("subject", event.type() + "|" + event.subject())
                .param("seq", event.sequence())
                .update();
        return advanced == 1;
    }

    /** Dùng khi chỉ cần loại trùng, không xét thứ tự. */
    @Transactional(propagation = Propagation.MANDATORY)
    public boolean firstTime(UUID eventId) {
        return jdbc.sql("INSERT INTO crucible_inbox (event_id) VALUES (:id) ON CONFLICT DO NOTHING")
                .param("id", eventId)
                .update() == 1;
    }
}
