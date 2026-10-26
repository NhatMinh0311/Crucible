package dev.crucible.learning;

import static org.assertj.core.api.Assertions.assertThat;

import dev.crucible.platform.events.EventEnvelope;
import dev.crucible.platform.events.InboxGuard;
import dev.crucible.platform.events.OutboxPublisher;
import dev.crucible.platform.events.OutboxRelay;
import java.nio.charset.StandardCharsets;
import java.util.Map;
import org.junit.jupiter.api.Test;
import org.springframework.amqp.core.AnonymousQueue;
import org.springframework.amqp.core.BindingBuilder;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.core.TopicExchange;
import org.springframework.amqp.rabbit.core.RabbitAdmin;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
import org.springframework.transaction.support.TransactionTemplate;
import tools.jackson.databind.json.JsonMapper;

/**
 * Kiểm tra "ống nước" của dịch vụ: khởi động với CSDL và broker thật, migration chạy, sự kiện đi từ outbox
 * lên RabbitMQ đúng định dạng CloudEvents, inbox loại trùng (CR-DOC-03 mục 5.2, INV7).
 */
@SpringBootTest(properties = "crucible.events.relay-interval=1h")
@Import(TestcontainersConfiguration.class)
class EventsPlumbingTest {

    @Autowired OutboxPublisher outbox;
    @Autowired OutboxRelay relay;
    @Autowired InboxGuard inbox;
    @Autowired RabbitAdmin admin;
    @Autowired RabbitTemplate rabbit;
    @Autowired TopicExchange exchange;
    @Autowired TransactionTemplate tx;
    @Autowired JsonMapper json;

    @Test
    void outboxEventReachesBrokerAsCloudEvent() {
        var queue = new AnonymousQueue();
        admin.declareQueue(queue);
        admin.declareBinding(BindingBuilder.bind(queue).to(exchange).with("test.plumbing.checked"));

        tx.executeWithoutResult(s -> outbox.publish("test.plumbing.checked", "subject-1", Map.of("ok", true)));
        assertThat(relay.relayOnce()).isEqualTo(1);

        Message message = rabbit.receive(queue.getName(), 5_000);
        assertThat(message).isNotNull();
        var event = json.readValue(new String(message.getBody(), StandardCharsets.UTF_8), EventEnvelope.class);
        assertThat(event.type()).isEqualTo("test.plumbing.checked");
        assertThat(event.subject()).isEqualTo("subject-1");
        assertThat(event.sequence()).isEqualTo(1);
        assertThat(event.source()).isEqualTo("/crucible/learning-service");
    }

    @Test
    void inboxRejectsDuplicateAndOlderEvents() {
        var first = envelope("11111111-1111-1111-1111-111111111111", 2);
        var duplicate = envelope("11111111-1111-1111-1111-111111111111", 2);
        var older = envelope("22222222-2222-2222-2222-222222222222", 1);

        assertThat(accept(first)).isTrue();
        assertThat(accept(duplicate)).isFalse();
        assertThat(accept(older)).isFalse();
    }

    private boolean accept(EventEnvelope<?> event) {
        return Boolean.TRUE.equals(tx.execute(s -> inbox.accept(event)));
    }

    private static EventEnvelope<Object> envelope(String id, long sequence) {
        return new EventEnvelope<>("1.0", java.util.UUID.fromString(id), "/test", "test.inbox", null,
                "subject-2", sequence, null, "application/json", Map.of());
    }
}
