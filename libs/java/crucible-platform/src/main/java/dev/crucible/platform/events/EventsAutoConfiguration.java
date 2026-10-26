package dev.crucible.platform.events;

import org.springframework.amqp.core.TopicExchange;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.jdbc.core.simple.JdbcClient;
import org.springframework.scheduling.annotation.EnableScheduling;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.transaction.support.TransactionTemplate;
import tools.jackson.databind.json.JsonMapper;

/** Đăng ký outbox, inbox, exchange và tác vụ chuyển tiếp định kỳ. */
@AutoConfiguration
@EnableScheduling
@EnableConfigurationProperties(EventsProperties.class)
public class EventsAutoConfiguration {

    @Bean
    TopicExchange crucibleEventsExchange(EventsProperties properties) {
        return new TopicExchange(properties.exchange(), true, false);
    }

    @Bean
    OutboxPublisher outboxPublisher(JdbcClient jdbc, JsonMapper json, EventsProperties properties) {
        return new OutboxPublisher(jdbc, json, properties);
    }

    @Bean
    OutboxRelay outboxRelay(JdbcClient jdbc, RabbitTemplate rabbit, TransactionTemplate tx,
                            EventsProperties properties) {
        return new OutboxRelay(jdbc, rabbit, tx, properties);
    }

    @Bean
    InboxGuard inboxGuard(JdbcClient jdbc) {
        return new InboxGuard(jdbc);
    }

    @Bean
    OutboxSchedule outboxSchedule(OutboxRelay relay) {
        return new OutboxSchedule(relay);
    }

    /** Tác vụ định kỳ; tách riêng để dịch vụ có thể tắt bằng cách ghi đè bean khi kiểm thử. */
    static class OutboxSchedule {
        private final OutboxRelay relay;

        OutboxSchedule(OutboxRelay relay) {
            this.relay = relay;
        }

        @Scheduled(fixedDelayString = "${crucible.events.relay-interval:1s}")
        void relay() {
            while (relay.relayOnce() > 0) {
                // Tiếp tục khi còn sự kiện chưa gửi.
            }
        }

        @Scheduled(fixedDelayString = "PT1H")
        void cleanup() {
            relay.cleanup();
        }
    }
}
