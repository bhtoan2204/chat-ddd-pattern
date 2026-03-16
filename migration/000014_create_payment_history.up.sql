-- =========================
-- TABLE: payment_histories
-- =========================

CREATE TABLE payment_histories (
    id VARCHAR2(36) PRIMARY KEY,
    type VARCHAR2(50) NOT NULL,
    amount NUMBER(19) NOT NULL,
    balance NUMBER(19) NOT NULL,
    sender_id VARCHAR2(36),
    receiver_id VARCHAR2(36),
    sender_name VARCHAR2(255),
    receiver_name VARCHAR2(255),
    properties CLOB NOT NULL,
    created_at TIMESTAMP DEFAULT SYSTIMESTAMP NOT NULL
);

CREATE INDEX idx_payment_histories_type
ON payment_histories(type);

CREATE INDEX idx_payment_histories_sender_id
ON payment_histories(sender_id);

CREATE INDEX idx_payment_histories_receiver_id
ON payment_histories(receiver_id);

CREATE INDEX idx_payment_histories_sender_created
ON payment_histories(sender_id, created_at);

CREATE INDEX idx_payment_histories_receiver_created
ON payment_histories(receiver_id, created_at);