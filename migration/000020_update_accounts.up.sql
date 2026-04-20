ALTER TABLE accounts ADD COLUMN display_name        VARCHAR(255);

ALTER TABLE accounts ADD COLUMN username            VARCHAR(100);

ALTER TABLE accounts ADD COLUMN avatar_object_key   VARCHAR(1024);

ALTER TABLE accounts ADD COLUMN status              VARCHAR(32) DEFAULT 'active' NOT NULL;

ALTER TABLE accounts ADD COLUMN email_verified_at   TIMESTAMPTZ;

ALTER TABLE accounts ADD COLUMN last_login_at       TIMESTAMPTZ;

ALTER TABLE accounts ADD COLUMN password_changed_at TIMESTAMPTZ;

UPDATE accounts
SET display_name = email
WHERE display_name IS NULL;

ALTER TABLE accounts ALTER COLUMN display_name TYPE VARCHAR(255);

ALTER TABLE accounts ALTER COLUMN display_name SET NOT NULL;

CREATE UNIQUE INDEX uq_accounts_email ON accounts(email);

CREATE UNIQUE INDEX uq_accounts_username ON accounts(username);

CREATE INDEX idx_accounts_status ON accounts(status);
