CREATE TABLE devices (
    id               VARCHAR(36)                        NOT NULL,
    account_id       VARCHAR(1024)                      NOT NULL,
    device_uid       VARCHAR(128)                       NOT NULL,
    device_name      VARCHAR(200),
    device_type      VARCHAR(30)    DEFAULT 'web'       NOT NULL,
    os_name          VARCHAR(50),
    os_version       VARCHAR(50),
    app_version      VARCHAR(50),
    user_agent       VARCHAR(1000),
    last_ip_address  VARCHAR(45),
    last_seen_at     TIMESTAMPTZ,
    is_trusted       NUMERIC(1,0)          DEFAULT 0 NOT NULL,
    created_at       TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at       TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_devices      PRIMARY KEY (id),
    CONSTRAINT fk_dev_acc      FOREIGN KEY (account_id)
                               REFERENCES accounts(id)
                               ON DELETE CASCADE,
    CONSTRAINT uk_dev_acc_uid  UNIQUE (account_id, device_uid),
    CONSTRAINT uk_dev_acc_id   UNIQUE (account_id, id),
    CONSTRAINT ck_dev_type     CHECK (device_type IN ('web', 'ios', 'android', 'desktop', 'other')),
    CONSTRAINT ck_dev_trusted  CHECK (is_trusted IN (0, 1))
);

CREATE INDEX ix_dev_seen ON devices (last_seen_at);

CREATE OR REPLACE FUNCTION trg_devices_bu_fn() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_devices_bu BEFORE UPDATE ON devices FOR EACH ROW EXECUTE FUNCTION trg_devices_bu_fn();
CREATE TABLE sessions (
    id                 VARCHAR(36)                        NOT NULL,
    account_id         VARCHAR(1024)                      NOT NULL,
    device_id          VARCHAR(36)                        NOT NULL,
    refresh_token_hash VARCHAR(255)                       NOT NULL,
    status             VARCHAR(20)    DEFAULT 'active'    NOT NULL,
    ip_address         VARCHAR(45),
    user_agent         VARCHAR(1000),
    last_activity_at   TIMESTAMPTZ,
    expires_at         TIMESTAMPTZ              NOT NULL,
    revoked_at         TIMESTAMPTZ,
    revoked_reason     VARCHAR(255),
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_sessions   PRIMARY KEY (id),
    CONSTRAINT fk_ses_dev    FOREIGN KEY (account_id, device_id)
                             REFERENCES devices(account_id, id)
                             ON DELETE CASCADE,
    CONSTRAINT uk_ses_rth    UNIQUE (refresh_token_hash),
    CONSTRAINT ck_ses_st     CHECK (status IN ('active', 'revoked', 'expired'))
);
CREATE INDEX ix_ses_acc_st   ON sessions (account_id, status);
CREATE INDEX ix_ses_acc_dev  ON sessions (account_id, device_id);
CREATE INDEX ix_ses_exp      ON sessions (expires_at);
CREATE INDEX ix_ses_last_act ON sessions (last_activity_at);
CREATE OR REPLACE FUNCTION trg_sessions_bu_fn() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_sessions_bu BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE FUNCTION trg_sessions_bu_fn();
