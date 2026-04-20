CREATE TABLE payment_account_projections (
    id          VARCHAR(1024) PRIMARY KEY,
    account_id  VARCHAR(1024) NOT NULL,
    email       VARCHAR(1024) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_payment_account_projections_account_id ON payment_account_projections(account_id);

CREATE TABLE payment_aggregates (
    id             VARCHAR(1024) PRIMARY KEY,
    aggregate_id   VARCHAR(1024) NOT NULL,
    aggregate_type VARCHAR(255)  NOT NULL,
    version        BIGINT     NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_payment_aggregates_aggregate_id ON payment_aggregates(aggregate_id);

CREATE INDEX idx_payment_aggregates_aggregate_type ON payment_aggregates(aggregate_type);

CREATE INDEX idx_payment_aggregates_version ON payment_aggregates(version);

CREATE TABLE payment_balances (
    id          VARCHAR(1024) PRIMARY KEY,
    account_id  VARCHAR(1024) NOT NULL,
    amount      BIGINT     NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_payment_balances_account_id ON payment_balances(account_id);

CREATE TABLE payment_balance_snapshots (
    id            VARCHAR(1024) PRIMARY KEY,
    aggregate_id  VARCHAR(1024) NOT NULL,
    version       BIGINT     NOT NULL,
    state         TEXT           NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_snap_ver ON payment_balance_snapshots(aggregate_id, version);

CREATE TABLE payment_events (
    id             VARCHAR(1024) PRIMARY KEY,
    aggregate_id   VARCHAR(1024) NOT NULL,
    aggregate_type VARCHAR(255)  NOT NULL,
    version        BIGINT     NOT NULL,
    event_name     VARCHAR(255)  NOT NULL,
    event_data     TEXT           NOT NULL,
    metadata       TEXT           NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_agg_ver ON payment_events(aggregate_id, version);

CREATE INDEX idx_payment_events_event_name ON payment_events(event_name);

CREATE TABLE payment_event_offsets (
    consumer_name VARCHAR(1024) PRIMARY KEY,
    last_event_id BIGINT     NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE payment_transactions (
    id          VARCHAR(1024) PRIMARY KEY,
    account_id  VARCHAR(1024) NOT NULL,
    event_id    VARCHAR(1024) NOT NULL,
    amount      BIGINT     NOT NULL,
    type        VARCHAR(255)  NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_payment_transactions_account_id ON payment_transactions(account_id);

CREATE INDEX idx_payment_transactions_event_id ON payment_transactions(event_id);

CREATE TABLE payment_histories (
    id            VARCHAR(36) PRIMARY KEY,
    type          VARCHAR(50)  NOT NULL,
    amount        BIGINT    NOT NULL,
    balance       BIGINT    NOT NULL,
    sender_id     VARCHAR(36),
    receiver_id   VARCHAR(36),
    sender_name   VARCHAR(255),
    receiver_name VARCHAR(255),
    properties    TEXT          NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_payment_histories_type ON payment_histories(type);

CREATE INDEX idx_payment_histories_sender_id ON payment_histories(sender_id);

CREATE INDEX idx_payment_histories_receiver_id ON payment_histories(receiver_id);

CREATE INDEX idx_payment_histories_sender_created ON payment_histories(sender_id, created_at);

CREATE INDEX idx_payment_histories_receiver_created ON payment_histories(receiver_id, created_at);
