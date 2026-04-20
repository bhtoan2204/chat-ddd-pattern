CREATE TABLE messages (
    id          VARCHAR(1024) PRIMARY KEY,
    room_id     VARCHAR(1024) NOT NULL,
    sender_id   VARCHAR(1024) NOT NULL,
    message     VARCHAR(4000) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_messages_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_messages_sender
        FOREIGN KEY (sender_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_messages_room_id ON messages(room_id);

CREATE INDEX idx_messages_sender_id ON messages(sender_id);

CREATE INDEX idx_messages_created_at ON messages(created_at);
