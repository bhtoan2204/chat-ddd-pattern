DROP INDEX IF EXISTS relationship_idx_outbox_unpublished;

DROP TABLE IF EXISTS relationship_outbox_events;

DROP TABLE IF EXISTS relationship_user_relationship_counters;

DROP INDEX IF EXISTS relationship_idx_blocks_blocked_created;
DROP INDEX IF EXISTS relationship_idx_blocks_blocker_created;
DROP INDEX IF EXISTS relationship_uq_blocks_pair;
DROP TABLE IF EXISTS relationship_blocks;

DROP INDEX IF EXISTS relationship_idx_follows_followee_created;
DROP INDEX IF EXISTS relationship_idx_follows_follower_created;
DROP INDEX IF EXISTS relationship_uq_follows_pair;
DROP TABLE IF EXISTS relationship_follows;

DROP INDEX IF EXISTS relationship_idx_friendships_user_high_created;
DROP INDEX IF EXISTS relationship_idx_friendships_user_low_created;
DROP INDEX IF EXISTS relationship_uq_friendships_pair;
DROP TABLE IF EXISTS relationship_friendships;

DROP INDEX IF EXISTS relationship_idx_friend_requests_addressee_status_created;
DROP INDEX IF EXISTS relationship_idx_friend_requests_requester_status_created;
DROP INDEX IF EXISTS relationship_uq_friend_requests_active_pair;
DROP TABLE IF EXISTS relationship_friend_requests;