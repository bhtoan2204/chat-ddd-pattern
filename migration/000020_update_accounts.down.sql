DROP INDEX idx_accounts_status;
DROP INDEX uq_accounts_username;
DROP INDEX uq_accounts_email;

ALTER TABLE accounts DROP COLUMN password_changed_at;
ALTER TABLE accounts DROP COLUMN last_login_at;
ALTER TABLE accounts DROP COLUMN email_verified_at;
ALTER TABLE accounts DROP COLUMN status;
ALTER TABLE accounts DROP COLUMN avatar_object_key;
ALTER TABLE accounts DROP COLUMN username;
ALTER TABLE accounts DROP COLUMN display_name;
