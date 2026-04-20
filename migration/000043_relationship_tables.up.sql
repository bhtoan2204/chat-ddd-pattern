CREATE TABLE relationship_friend_requests (
    id                  VARCHAR(36) PRIMARY KEY,
    requester_id        VARCHAR(36) NOT NULL,
    addressee_id        VARCHAR(36) NOT NULL,
    status              VARCHAR(20) NOT NULL,
    message             VARCHAR(500),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    responded_at        TIMESTAMPTZ,
    expired_at          TIMESTAMPTZ,
    cancelled_at        TIMESTAMPTZ,
    rejected_reason     VARCHAR(255),

    CONSTRAINT chk_friend_requests_status
        CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED', 'CANCELLED', 'EXPIRED')),

    CONSTRAINT chk_friend_requests_not_self
        CHECK (requester_id <> addressee_id)
);

CREATE UNIQUE INDEX uq_friend_requests_active_pair
ON relationship_friend_requests (
    LEAST(requester_id, addressee_id),
    GREATEST(requester_id, addressee_id)
)
WHERE status = 'PENDING';

CREATE INDEX idx_friend_requests_requester_status_created
ON relationship_friend_requests (requester_id, status, created_at DESC);

CREATE INDEX idx_friend_requests_addressee_status_created
ON relationship_friend_requests (addressee_id, status, created_at DESC);

CREATE TABLE relationship_friendships (
    id                      VARCHAR(36) PRIMARY KEY,
    user_low_id             VARCHAR(36) NOT NULL,
    user_high_id            VARCHAR(36) NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_from_request_id VARCHAR(36),

    CONSTRAINT chk_friendships_order
        CHECK (user_low_id < user_high_id)
);

CREATE UNIQUE INDEX uq_friendships_pair
ON relationship_friendships (user_low_id, user_high_id);

CREATE INDEX idx_friendships_user_low_created
ON relationship_friendships (user_low_id, created_at DESC);

CREATE INDEX idx_friendships_user_high_created
ON relationship_friendships (user_high_id, created_at DESC);

CREATE TABLE relationship_follows (
    id                  VARCHAR(36) PRIMARY KEY,
    follower_id         VARCHAR(36) NOT NULL,
    followee_id         VARCHAR(36) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_follows_not_self
        CHECK (follower_id <> followee_id)
);

CREATE UNIQUE INDEX uq_follows_pair
ON relationship_follows (follower_id, followee_id);

CREATE INDEX idx_follows_follower_created
ON relationship_follows (follower_id, created_at DESC);

CREATE INDEX idx_follows_followee_created
ON relationship_follows (followee_id, created_at DESC);

CREATE TABLE relationship_blocks (
    id                  VARCHAR(36) PRIMARY KEY,
    blocker_id          VARCHAR(36) NOT NULL,
    blocked_id          VARCHAR(36) NOT NULL,
    reason              VARCHAR(255),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_blocks_not_self
        CHECK (blocker_id <> blocked_id)
);

CREATE UNIQUE INDEX uq_blocks_pair
ON relationship_blocks (blocker_id, blocked_id);

CREATE INDEX idx_blocks_blocker_created
ON relationship_blocks (blocker_id, created_at DESC);

CREATE INDEX idx_blocks_blocked_created
ON relationship_blocks (blocked_id, created_at DESC);

CREATE TABLE relationship_user_relationship_counters (
    user_id              VARCHAR(36) PRIMARY KEY,
    friends_count        BIGINT NOT NULL DEFAULT 0,
    followers_count      BIGINT NOT NULL DEFAULT 0,
    following_count      BIGINT NOT NULL DEFAULT 0,
    blocked_count        BIGINT NOT NULL DEFAULT 0,
    pending_in_count     BIGINT NOT NULL DEFAULT 0,
    pending_out_count    BIGINT NOT NULL DEFAULT 0,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE relationship_outbox_events (
    id                VARCHAR(36) PRIMARY KEY,
    aggregate_type    VARCHAR(50) NOT NULL,
    aggregate_id      VARCHAR(36) NOT NULL,
    event_type        VARCHAR(100) NOT NULL,
    payload           JSONB NOT NULL,
    occurred_at       TIMESTAMPTZ NOT NULL,
    published_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_outbox_unpublished
ON relationship_outbox_events (created_at)
WHERE published_at IS NULL;