CREATE TABLE payment_histories (
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    amount BIGINT NOT NULL,
    balance BIGINT NOT NULL,
    sender_id VARCHAR(36),
    receiver_id VARCHAR(36),
    sender_name VARCHAR(255),
    receiver_name VARCHAR(255),
    properties TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
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
