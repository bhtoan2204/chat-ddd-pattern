ALTER TABLE rooms DROP CONSTRAINT fk_rooms_owner;

DROP INDEX idx_rooms_owner_id;

ALTER TABLE rooms DROP COLUMN owner_id;

ALTER TABLE rooms DROP COLUMN owner_type;

ALTER TABLE rooms ADD COLUMN description VARCHAR(1024) DEFAULT '' NOT NULL;

ALTER TABLE rooms ADD COLUMN room_type   VARCHAR(50)   DEFAULT 'public' NOT NULL;
