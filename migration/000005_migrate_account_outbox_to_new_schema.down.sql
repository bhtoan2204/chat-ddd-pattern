-- =========================
-- ROLLBACK: account_outbox_events -> old schema
-- =========================
ALTER TABLE account_outbox_events RENAME TO account_outbox_events_new;

CREATE TABLE account_outbox_events (
    id          VARCHAR2(1024) PRIMARY KEY,
    event_name  VARCHAR2(1024) NOT NULL,
    event_data  VARCHAR2(4000) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL
);

INSERT INTO account_outbox_events (
    id,
    event_name,
    event_data,
    created_at
)
SELECT
    TO_CHAR(id) AS id,
    event_name,
    DBMS_LOB.SUBSTR(event_data, 4000, 1) AS event_data,
    created_at
FROM account_outbox_events_new;

DROP TABLE account_outbox_events_new CASCADE CONSTRAINTS;

CREATE INDEX idx_account_outbox_events_event_name ON account_outbox_events(event_name);
CREATE INDEX idx_account_outbox_events_created_at ON account_outbox_events(created_at);
