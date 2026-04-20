ALTER TABLE accounts ADD COLUMN banned_reason VARCHAR(1024);

ALTER TABLE accounts ADD COLUMN banned_until  TIMESTAMPTZ;
