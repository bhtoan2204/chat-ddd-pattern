-- =========================
-- TABLE: payment_aggregates
-- =========================
CREATE TABLE payment_aggregates (
    id            VARCHAR2(1024) PRIMARY KEY,
    aggregate_id  VARCHAR2(1024) NOT NULL,
    aggregate_type VARCHAR2(255)  NOT NULL,
    version       NUMBER(10)     NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_payment_aggregates_aggregate_id ON payment_aggregates(aggregate_id);
CREATE INDEX idx_payment_aggregates_aggregate_type ON payment_aggregates(aggregate_type);
CREATE INDEX idx_payment_aggregates_version ON payment_aggregates(version);