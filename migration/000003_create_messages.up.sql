-- =========================
-- TABLE: messages
-- =========================
CREATE TABLE messages (
    id          VARCHAR2(1024) PRIMARY KEY,
    room_id     VARCHAR2(1024) NOT NULL,
    sender_id   VARCHAR2(1024) NOT NULL,
    message     VARCHAR2(4000) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
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
