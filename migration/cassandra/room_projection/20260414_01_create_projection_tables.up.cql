CREATE TABLE IF NOT EXISTS room_projections_by_id (
	room_id text PRIMARY KEY,
	name text,
	description text,
	room_type text,
	owner_id text,
	pinned_message_id text,
	member_count int,
	last_message_id text,
	last_message_at timestamp,
	last_message_content text,
	last_message_sender_id text,
	created_at timestamp,
	updated_at timestamp
);

CREATE TABLE IF NOT EXISTS room_projections_by_account (
	account_id text,
	room_updated_at timestamp,
	room_id text,
	name text,
	description text,
	room_type text,
	owner_id text,
	pinned_message_id text,
	member_count int,
	last_message_id text,
	last_message_at timestamp,
	last_message_content text,
	last_message_sender_id text,
	created_at timestamp,
	PRIMARY KEY ((account_id), room_updated_at, room_id)
) WITH CLUSTERING ORDER BY (room_updated_at DESC, room_id DESC);

CREATE TABLE IF NOT EXISTS room_member_projections_by_room (
	room_id text,
	account_id text,
	member_id text,
	display_name text,
	username text,
	avatar_object_key text,
	role text,
	last_delivered_at timestamp,
	last_read_at timestamp,
	created_at timestamp,
	updated_at timestamp,
	PRIMARY KEY ((room_id), account_id)
);

CREATE TABLE IF NOT EXISTS room_message_timelines (
	room_id text,
	message_sent_at timestamp,
	message_id text,
	room_name text,
	room_type text,
	message_content text,
	message_type text,
	reply_to_message_id text,
	forwarded_from_message_id text,
	file_name text,
	file_size bigint,
	mime_type text,
	object_key text,
	message_sender_id text,
	message_sender_name text,
	message_sender_email text,
	mentions_json text,
	mention_all boolean,
	mentioned_account_ids list<text>,
	edited_at timestamp,
	deleted_for_everyone_at timestamp,
	PRIMARY KEY ((room_id), message_sent_at, message_id)
) WITH CLUSTERING ORDER BY (message_sent_at DESC, message_id DESC);

CREATE TABLE IF NOT EXISTS room_messages_by_id (
	message_id text PRIMARY KEY,
	room_id text,
	room_name text,
	room_type text,
	message_content text,
	message_type text,
	reply_to_message_id text,
	forwarded_from_message_id text,
	file_name text,
	file_size bigint,
	mime_type text,
	object_key text,
	message_sender_id text,
	message_sender_name text,
	message_sender_email text,
	message_sent_at timestamp,
	mentions_json text,
	mention_all boolean,
	mentioned_account_ids list<text>,
	edited_at timestamp,
	deleted_for_everyone_at timestamp
);

CREATE TABLE IF NOT EXISTS room_message_receipts_by_message (
	message_id text,
	account_id text,
	room_id text,
	status text,
	delivered_at timestamp,
	seen_at timestamp,
	created_at timestamp,
	updated_at timestamp,
	PRIMARY KEY ((message_id), account_id)
);

CREATE TABLE IF NOT EXISTS room_message_deletions_by_account_room (
	account_id text,
	room_id text,
	message_sent_at timestamp,
	message_id text,
	created_at timestamp,
	PRIMARY KEY ((account_id, room_id), message_sent_at, message_id)
) WITH CLUSTERING ORDER BY (message_sent_at DESC, message_id DESC);
