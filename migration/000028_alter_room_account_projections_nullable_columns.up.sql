ALTER TABLE room_account_projections ALTER COLUMN username TYPE VARCHAR(1024);

ALTER TABLE room_account_projections ALTER COLUMN username DROP NOT NULL;

ALTER TABLE room_account_projections ALTER COLUMN avatar_object_key TYPE VARCHAR(2048);

ALTER TABLE room_account_projections ALTER COLUMN avatar_object_key DROP NOT NULL;
