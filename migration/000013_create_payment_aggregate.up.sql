CREATE TABLE payment_aggregates (
    id            VARCHAR(1024) PRIMARY KEY,
    aggregate_id  VARCHAR(1024) NOT NULL,
    aggregate_type VARCHAR(255)  NOT NULL,
    version       BIGINT     NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_payment_aggregates_aggregate_id ON payment_aggregates(aggregate_id);

CREATE INDEX idx_payment_aggregates_aggregate_type ON payment_aggregates(aggregate_type);

CREATE INDEX idx_payment_aggregates_version ON payment_aggregates(version);
