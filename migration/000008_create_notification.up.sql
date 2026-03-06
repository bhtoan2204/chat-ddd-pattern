-- =========================
-- TABLE: notifications
-- =========================
CREATE TABLE notifications (
    id          VARCHAR2(1024) PRIMARY KEY,
    account_id  VARCHAR2(1024) NOT NULL,
    type        VARCHAR2(1024) NOT NULL,
    subject     VARCHAR2(1024) NOT NULL,
    body        VARCHAR2(1024) NOT NULL,
    is_read     BOOLEAN DEFAULT FALSE NOT NULL,
    read_at     TIMESTAMP WITH TIME ZONE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL
);

CREATE INDEX idx_notifications_account_id ON notifications(account_id);
CREATE INDEX idx_notifications_type ON notifications(type);