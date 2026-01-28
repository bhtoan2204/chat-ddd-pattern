-- =========================
-- TABLE: rooms (update for description & room_type)
-- =========================
ALTER TABLE rooms DROP CONSTRAINT fk_rooms_owner;
DROP INDEX idx_rooms_owner_id;

ALTER TABLE rooms DROP COLUMN owner_id;
ALTER TABLE rooms DROP COLUMN owner_type;

ALTER TABLE rooms ADD (
    description VARCHAR2(1024) DEFAULT '' NOT NULL,
    room_type   VARCHAR2(50)   DEFAULT 'public' NOT NULL
);
