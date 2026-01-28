-- =========================
-- TABLE: rooms (rollback to owner-based model)
-- =========================
ALTER TABLE rooms DROP COLUMN description;
ALTER TABLE rooms DROP COLUMN room_type;

ALTER TABLE rooms ADD (
    owner_id   VARCHAR2(1024),
    owner_type VARCHAR2(1024)
);

ALTER TABLE rooms
    ADD CONSTRAINT fk_rooms_owner
    FOREIGN KEY (owner_id)
    REFERENCES accounts(id)
    ON DELETE CASCADE;

CREATE INDEX idx_rooms_owner_id ON rooms(owner_id);
