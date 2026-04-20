CREATE TABLE payment_account_projections (
    id          VARCHAR(1024) PRIMARY KEY,
    account_id  VARCHAR(1024) NOT NULL,
    email       VARCHAR(1024) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_payment_account_projections_account_id ON payment_account_projections(account_id);
