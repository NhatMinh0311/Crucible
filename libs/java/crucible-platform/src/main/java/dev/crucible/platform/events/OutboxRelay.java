package dev.crucible.platform.events;

import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.UUID;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.core.MessageBuilder;
import org.springframework.amqp.core.MessageDeliveryMode;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.jdbc.core.simple.JdbcClient;
import org.springframework.transaction.support.TransactionTemplate;

/**
 * Chuyển sự kiện từ outbox lên RabbitMQ. Nhiều bản sao dịch vụ chạy đồng thời được nhờ {@code FOR UPDATE SKIP LOCKED};
 * một sự kiện chỉ được đánh dấu đã gửi sau khi broker xác nhận (publisher confirm).
 */
public class OutboxRelay {

    private static final Logger log = LoggerFactory.getLogger(OutboxRelay.class);

    private final JdbcClient jdbc;
    private final RabbitTemplate rabbit;
    private final TransactionTemplate tx;
    private final EventsProperties properties;

    public OutboxRelay(JdbcClient jdbc, RabbitTemplate rabbit, TransactionTemplate tx, EventsProperties properties) {
        this.jdbc = jdbc;
        this.rabbit = rabbit;
        this.tx = tx;
        this.properties = properties;
    }

    private record Pending(UUID id, String type, String payload) {
    }

    /** Gửi một lô sự kiện; trả về số sự kiện đã gửi. */
    public int relayOnce() {
        Integer sent = tx.execute(status -> {
            List<Pending> batch = jdbc.sql("""
                            SELECT id, type, payload::text AS payload FROM crucible_outbox
                            WHERE sent_at IS NULL ORDER BY created_at LIMIT :limit
                            FOR UPDATE SKIP LOCKED
                            """)
                    .param("limit", properties.relayBatch())
                    .query(Pending.class)
                    .list();
            if (batch.isEmpty()) {
                return 0;
            }
            rabbit.invoke(ops -> {
                for (Pending event : batch) {
                    ops.send(properties.exchange(), event.type(), toMessage(event));
                }
                ops.waitForConfirmsOrDie(5_000);
                return null;
            });
            jdbc.sql("UPDATE crucible_outbox SET sent_at = now() WHERE id IN (:ids)")
                    .param("ids", batch.stream().map(Pending::id).toList())
                    .update();
            return batch.size();
        });
        return sent == null ? 0 : sent;
    }

    /** Xóa dòng outbox đã gửi và dòng inbox quá thời gian lưu giữ. */
    public void cleanup() {
        long seconds = properties.retention().toSeconds();
        int outbox = jdbc.sql("DELETE FROM crucible_outbox WHERE sent_at < now() - make_interval(secs => :s)")
                .param("s", seconds).update();
        int inbox = jdbc.sql("DELETE FROM crucible_inbox WHERE processed_at < now() - make_interval(secs => :s)")
                .param("s", seconds).update();
        log.debug("Đã xóa {} dòng outbox, {} dòng inbox quá hạn", outbox, inbox);
    }

    private static Message toMessage(Pending event) {
        return MessageBuilder.withBody(event.payload().getBytes(StandardCharsets.UTF_8))
                .setContentType("application/cloudevents+json")
                .setContentEncoding("utf-8")
                .setMessageId(event.id().toString())
                .setType(event.type())
                .setDeliveryMode(MessageDeliveryMode.PERSISTENT)
                .build();
    }
}
