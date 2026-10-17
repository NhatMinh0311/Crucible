package dev.crucible.platform.events;

import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.UUID;
import org.springframework.jdbc.core.simple.JdbcClient;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;
import tools.jackson.databind.json.JsonMapper;

/**
 * Ghi sự kiện vào outbox trong transaction nghiệp vụ hiện tại. Sự kiện chỉ được gửi lên broker khi transaction
 * được xác nhận (ADR-0005).
 */
public class OutboxPublisher {

    private final JdbcClient jdbc;
    private final JsonMapper json;
    private final EventsProperties properties;

    public OutboxPublisher(JdbcClient jdbc, JsonMapper json, EventsProperties properties) {
        this.jdbc = jdbc;
        this.json = json;
        this.properties = properties;
    }

    /**
     * Ghi một sự kiện.
     *
     * @param type    kiểu sự kiện, trùng routing key (ví dụ {@code content.attempt.graded})
     * @param subject mã đối tượng mà sự kiện nói về
     * @param data    nội dung {@code data}, đúng lược đồ trong {@code contracts/events.asyncapi.yaml}
     * @return mã sự kiện
     */
    @Transactional(propagation = Propagation.MANDATORY)
    public UUID publish(String type, String subject, Object data) {
        long sequence = jdbc.sql("""
                        INSERT INTO crucible_event_sequence (subject, last_sequence) VALUES (:subject, 1)
                        ON CONFLICT (subject) DO UPDATE
                            SET last_sequence = crucible_event_sequence.last_sequence + 1
                        RETURNING last_sequence
                        """)
                .param("subject", subject)
                .query(Long.class)
                .single();
        UUID id = UUID.randomUUID();
        var envelope = new EventEnvelope<>("1.0", id, properties.source(), type,
                OffsetDateTime.now(ZoneOffset.UTC), subject, sequence, null, "application/json", data);
        jdbc.sql("""
                        INSERT INTO crucible_outbox (id, type, subject, sequence, payload)
                        VALUES (:id, :type, :subject, :sequence, CAST(:payload AS jsonb))
                        """)
                .param("id", id)
                .param("type", type)
                .param("subject", subject)
                .param("sequence", sequence)
                .param("payload", json.writeValueAsString(envelope))
                .update();
        return id;
    }
}
