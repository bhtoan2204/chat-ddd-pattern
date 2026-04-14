package views

type ProjectionTableNames struct {
	RoomByID             string
	RoomByAccount        string
	RoomMemberByRoom     string
	MessageByID          string
	MessageReceipts      string
	MessageDeletions     string
	MessageTimelines     string
	GlobalRoomProjection string
	SchemaMigrations     string
}

func DefaultProjectionTableNames() ProjectionTableNames {
	return ProjectionTableNames{
		RoomByID:             "room_projections_by_id",
		RoomByAccount:        "room_projections_by_account",
		RoomMemberByRoom:     "room_member_projections_by_room",
		MessageByID:          "room_messages_by_id",
		MessageReceipts:      "room_message_receipts_by_message",
		MessageDeletions:     "room_message_deletions_by_account_room",
		MessageTimelines:     "room_message_timelines",
		GlobalRoomProjection: "room_projections_global",
		SchemaMigrations:     "room_projection_schema_migrations",
	}
}
