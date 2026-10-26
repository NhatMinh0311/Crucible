-- Bảng phiên, vé, exec của runtime-api (CR-MOD-M1 mục 7). Bảng hạn mức, đặt trước, tham số được thêm khi hiện thực.

CREATE TABLE sessions (
    id                    uuid PRIMARY KEY,
    owner_user_id         uuid        NOT NULL,
    priority              text        NOT NULL,
    profile               text        NOT NULL,
    image                 text        NOT NULL,
    quota_account_ref     text        NOT NULL,
    reservation_ref       text,
    record_input          boolean     NOT NULL DEFAULT false,
    labels                jsonb       NOT NULL DEFAULT '{}',
    environment           jsonb       NOT NULL,
    phase                 text        NOT NULL,
    end_reason            text,
    failure_detail        text,
    created_at            timestamptz NOT NULL DEFAULT now(),
    ready_at              timestamptz,
    expires_at            timestamptz,
    ended_at              timestamptz,
    ttl_minutes           int         NOT NULL,
    idle_minutes          int         NOT NULL,
    max_extension_minutes int         NOT NULL,
    extension_minutes     int         NOT NULL DEFAULT 0,
    last_activity_at      timestamptz,
    idle_warned_at        timestamptz,
    expiring_notified     boolean     NOT NULL DEFAULT false,
    used_minutes          int,
    container_id          text,
    memory_peak_bytes     bigint,
    memory_limit_bytes    bigint      NOT NULL,
    oom_killed            boolean     NOT NULL DEFAULT false
);

-- BR-M1.1: mỗi học viên tối đa một phiên PRACTICE và một phiên EXAM chưa kết thúc.
CREATE UNIQUE INDEX sessions_one_active ON sessions (owner_user_id, priority)
    WHERE phase IN ('PENDING', 'PROVISIONING', 'READY', 'EXPIRING') AND priority <> 'LAB_CI';
CREATE INDEX sessions_active ON sessions (phase)
    WHERE phase IN ('PENDING', 'PROVISIONING', 'READY', 'EXPIRING');

CREATE TABLE terminal_tickets (
    token_hash  bytea PRIMARY KEY,
    session_id  uuid        NOT NULL REFERENCES sessions (id),
    mode        text        NOT NULL,
    user_id     uuid        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    expires_at  timestamptz NOT NULL,
    used_at     timestamptz
);

CREATE TABLE view_access_log (
    id              bigserial PRIMARY KEY,
    session_id      uuid        NOT NULL REFERENCES sessions (id),
    viewer_user_id  uuid        NOT NULL,
    issued_at       timestamptz NOT NULL DEFAULT now()
);

-- Lặp lại an toàn (CR-DOC-03 mục 3.3), khóa theo người gọi.
CREATE TABLE idempotency_keys (
    caller        text        NOT NULL,
    key           uuid        NOT NULL,
    operation     text        NOT NULL,
    request_hash  bytea       NOT NULL,
    status        int,
    response      jsonb,
    created_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (caller, key, operation)
);

-- Outbox cùng cấu trúc với thư viện dùng chung của M2, M3 (CR-DOC-03 mục 5.2).
CREATE TABLE crucible_outbox (
    id          uuid PRIMARY KEY,
    type        text        NOT NULL,
    subject     text        NOT NULL,
    sequence    bigint      NOT NULL,
    payload     jsonb       NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    sent_at     timestamptz
);
CREATE INDEX crucible_outbox_unsent ON crucible_outbox (created_at) WHERE sent_at IS NULL;

CREATE TABLE crucible_event_sequence (
    subject        text   PRIMARY KEY,
    last_sequence  bigint NOT NULL
);
