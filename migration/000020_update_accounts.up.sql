ALTER TABLE accounts ADD (
    display_name        VARCHAR2(255),
    username            VARCHAR2(100),
    avatar_object_key   VARCHAR2(1024),
    status              VARCHAR2(32) DEFAULT 'active' NOT NULL,
    email_verified_at   TIMESTAMP WITH TIME ZONE,
    last_login_at       TIMESTAMP WITH TIME ZONE,
    password_changed_at TIMESTAMP WITH TIME ZONE
);

UPDATE accounts
SET display_name = email
WHERE display_name IS NULL;

ALTER TABLE accounts MODIFY (display_name VARCHAR2(255) NOT NULL);

CREATE UNIQUE INDEX uq_accounts_email ON accounts(email);
CREATE UNIQUE INDEX uq_accounts_username ON accounts(username);
CREATE INDEX idx_accounts_status ON accounts(status);
