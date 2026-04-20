ALTER TABLE account_outbox_events RENAME TO account_outbox_events_new;

CREATE TABLE account_outbox_events (
    id          VARCHAR(1024) PRIMARY KEY,
    event_name  VARCHAR(1024) NOT NULL,
    event_data  VARCHAR(4000) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

INSERT INTO account_outbox_events (
    id,
    event_name,
    event_data,
    created_at
)
SELECT
    id::text AS id,
    event_name,
    LEFT(event_data, 4000) AS event_data,
    created_at
FROM account_outbox_events_new;

DROP TABLE account_outbox_events_new CASCADE;

CREATE INDEX idx_account_outbox_events_event_name ON account_outbox_events(event_name);

CREATE INDEX idx_account_outbox_events_created_at ON account_outbox_events(created_at);
