-- =========================
-- TABLE: accounts (ban fields)
-- =========================
ALTER TABLE accounts
    ADD (
        banned_reason VARCHAR2(1024),
        banned_until  TIMESTAMP WITH TIME ZONE
    );
