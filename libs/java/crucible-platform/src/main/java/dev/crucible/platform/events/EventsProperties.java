package dev.crucible.platform.events;

import java.time.Duration;
import org.springframework.boot.context.properties.ConfigurationProperties;

/**
 * Cấu hình sự kiện của dịch vụ.
 *
 * @param source        giá trị {@code source} của CloudEvents, ví dụ {@code /crucible/content-service}
 * @param exchange      topic exchange dùng chung (CR-DOC-02 mục 6.2)
 * @param relayInterval chu kỳ chuyển tiếp outbox
 * @param relayBatch    số sự kiện tối đa mỗi lần chuyển tiếp
 * @param retention     thời gian giữ dòng outbox đã gửi và dòng inbox
 */
@ConfigurationProperties("crucible.events")
public record EventsProperties(
        String source,
        String exchange,
        Duration relayInterval,
        int relayBatch,
        Duration retention) {

    public EventsProperties {
        if (source == null || source.isBlank()) {
            throw new IllegalArgumentException("crucible.events.source phải được cấu hình");
        }
        exchange = exchange == null ? "crucible.events" : exchange;
        relayInterval = relayInterval == null ? Duration.ofSeconds(1) : relayInterval;
        relayBatch = relayBatch <= 0 ? 100 : relayBatch;
        retention = retention == null ? Duration.ofDays(7) : retention;
    }
}
