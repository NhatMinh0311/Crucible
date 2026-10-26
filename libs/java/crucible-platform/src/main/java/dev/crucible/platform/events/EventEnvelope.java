package dev.crucible.platform.events;

import java.time.OffsetDateTime;
import java.util.UUID;

/**
 * Sự kiện theo CloudEvents 1.0 dạng JSON có cấu trúc, kèm hai thuộc tính mở rộng {@code sequence} và
 * {@code traceparent} (CR-DOC-03 mục 5.1).
 *
 * @param <T> kiểu của {@code data}
 */
public record EventEnvelope<T>(
        String specversion,
        UUID id,
        String source,
        String type,
        OffsetDateTime time,
        String subject,
        long sequence,
        String traceparent,
        String datacontenttype,
        T data) {
}
