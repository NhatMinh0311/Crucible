-- Bảng outbox, inbox và số thứ tự sự kiện theo đối tượng (CR-DOC-03 mục 5.2).
-- Phiên bản 0.x để luôn chạy trước migration của dịch vụ (V1 trở đi).

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

CREATE TABLE crucible_inbox (
    event_id      uuid        PRIMARY KEY,
    processed_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE crucible_inbox_sequence (
    subject        text   PRIMARY KEY,
    last_sequence  bigint NOT NULL
);
