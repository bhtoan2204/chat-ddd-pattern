package projection

import (
	"context"
	"time"
)

type TimelineProjector interface {
	UpsertMessage(ctx context.Context, projection *TimelineMessageProjection) error
}

type MessageSearchIndexer interface {
	UpsertMessage(ctx context.Context, document *SearchMessageDocument) error
}

type ProjectionMention struct {
	AccountID   string `json:"account_id"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
}

type TimelineMessageProjection struct {
	RoomID                 string              `json:"room_id"`
	RoomName               string              `json:"room_name"`
	RoomType               string              `json:"room_type"`
	MessageID              string              `json:"message_id"`
	MessageContent         string              `json:"message_content"`
	MessageType            string              `json:"message_type"`
	ReplyToMessageID       string              `json:"reply_to_message_id"`
	ForwardedFromMessageID string              `json:"forwarded_from_message_id"`
	FileName               string              `json:"file_name"`
	FileSize               int64               `json:"file_size"`
	MimeType               string              `json:"mime_type"`
	ObjectKey              string              `json:"object_key"`
	MessageSenderID        string              `json:"message_sender_id"`
	MessageSenderName      string              `json:"message_sender_name"`
	MessageSenderEmail     string              `json:"message_sender_email"`
	MessageSentAt          time.Time           `json:"message_sent_at"`
	Mentions               []ProjectionMention `json:"mentions"`
	MentionAll             bool                `json:"mention_all"`
	MentionedAccountIDs    []string            `json:"mentioned_account_ids"`
}

type SearchMessageDocument struct {
	RoomID                 string              `json:"room_id"`
	RoomName               string              `json:"room_name"`
	RoomType               string              `json:"room_type"`
	MessageID              string              `json:"message_id"`
	MessageContent         string              `json:"message_content"`
	MessageType            string              `json:"message_type"`
	ReplyToMessageID       string              `json:"reply_to_message_id"`
	ForwardedFromMessageID string              `json:"forwarded_from_message_id"`
	FileName               string              `json:"file_name"`
	FileSize               int64               `json:"file_size"`
	MimeType               string              `json:"mime_type"`
	ObjectKey              string              `json:"object_key"`
	MessageSenderID        string              `json:"message_sender_id"`
	MessageSenderName      string              `json:"message_sender_name"`
	MessageSenderEmail     string              `json:"message_sender_email"`
	MessageSentAt          time.Time           `json:"message_sent_at"`
	Mentions               []ProjectionMention `json:"mentions"`
	MentionAll             bool                `json:"mention_all"`
	MentionedAccountIDs    []string            `json:"mentioned_account_ids"`
}
