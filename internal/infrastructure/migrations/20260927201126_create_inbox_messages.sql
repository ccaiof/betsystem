-- +goose Up
CREATE TYPE outbox_status AS ENUM ('PENDING', 'PUBLISHED', 'FAILED');

CREATE TABLE
    inbox_messages (
        consumer_name TEXT NOT NULL,
        message_id TEXT NOT NULL,
        payload_hash TEXT NOT NULL,
        received_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        completed_at TIMESTAMPTZ,
        PRIMARY KEY (consumer_name, message_id)
    );

CREATE TABLE
    outbox_events (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        aggregate_id UUID NOT NULL,
        event_type TEXT NOT NULL,
        correlation_id UUID,
        causation_id UUID,
        payload JSONB NOT NULL,
        occurred_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        attempts INT NOT NULL DEFAULT 0,
        next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        locked_by TEXT,
        locked_at TIMESTAMPTZ,
        published_at TIMESTAMPTZ,
        status outbox_status NOT NULL DEFAULT 'PENDING'
    );

CREATE INDEX idx_outbox_dispatch ON outbox_events (next_attempt_at)
WHERE
    status = 'PENDING';

CREATE INDEX idx_outbox_aggregate_id ON outbox_events (aggregate_id);

-- +goose Down
DROP TABLE IF EXISTS outbox_events;

DROP TABLE IF EXISTS inbox_messages;