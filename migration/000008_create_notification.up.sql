CREATE TABLE notifications (
    id          VARCHAR(1024) PRIMARY KEY,
    account_id  VARCHAR(1024) NOT NULL,
    type        VARCHAR(1024) NOT NULL,
    subject     VARCHAR(1024) NOT NULL,
    body        VARCHAR(1024) NOT NULL,
    is_read     BOOLEAN DEFAULT FALSE NOT NULL,
    read_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_notifications_account_id ON notifications(account_id);

CREATE INDEX idx_notifications_type ON notifications(type);
