CREATE TABLE ledger_aggregates (
    id             VARCHAR(1024) PRIMARY KEY,
    aggregate_id   VARCHAR(1024) NOT NULL,
    aggregate_type VARCHAR(255)  NOT NULL,
    version        BIGINT     NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_ledger_aggregates_aggregate_id
    ON ledger_aggregates(aggregate_id);

CREATE INDEX idx_ledger_aggregates_aggregate_type
    ON ledger_aggregates(aggregate_type);

CREATE INDEX idx_ledger_aggregates_version
    ON ledger_aggregates(version);

CREATE TABLE ledger_events (
    id             VARCHAR(1024) PRIMARY KEY,
    aggregate_id   VARCHAR(1024) NOT NULL,
    aggregate_type VARCHAR(255)  NOT NULL,
    version        BIGINT     NOT NULL,
    event_name     VARCHAR(255)  NOT NULL,
    event_data     TEXT           NOT NULL,
    metadata       TEXT           NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_ledger_events_agg_ver
    ON ledger_events(aggregate_id, version);

CREATE INDEX idx_ledger_events_event_name
    ON ledger_events(event_name);
