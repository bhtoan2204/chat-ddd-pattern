-- =========================
-- ROLLBACK: rooms (remove owner_id column from room table)
-- =========================
ALTER TABLE rooms DROP COLUMN owner_id;

ALTER TABLE rooms
    DROP CONSTRAINT fk_rooms_owner;

DROP INDEX idx_rooms_owner_id;