package entity

import (
	"errors"
	"strings"
	"time"
)

const (
	MessageTypeText   = "text"
	MessageTypeSystem = "system"
	MessageTypeImage  = "image"
	MessageTypeFile   = "file"
)

var (
	ErrMessageIDRequired           = errors.New("message_id is required")
	ErrMessageRoomRequired         = errors.New("room_id is required")
	ErrMessageSenderRequired       = errors.New("account_id is required")
	ErrMessageBodyRequired         = errors.New("message is required")
	ErrMessageTypeInvalid          = errors.New("message_type is invalid")
	ErrMessageObjectKeyRequired    = errors.New("object_key is required for media messages")
	ErrMessageCannotEditOther      = errors.New("cannot edit another user's message")
	ErrMessageCannotEditSystem     = errors.New("system messages cannot be edited")
	ErrMessageCannotDeleteEveryone = errors.New("cannot delete everyone for another user's message")
)

type MessageParams struct {
	Message                string
	MessageType            string
	ReplyToMessageID       string
	ForwardedFromMessageID string
	FileName               string
	FileSize               int64
	MimeType               string
	ObjectKey              string
}

func NewMessage(id, roomID, senderID string, params MessageParams, now time.Time) (*MessageEntity, error) {
	id = strings.TrimSpace(id)
	roomID = strings.TrimSpace(roomID)
	senderID = strings.TrimSpace(senderID)
	messageType := NormalizeMessageType(params.MessageType)
	content := strings.TrimSpace(params.Message)
	objectKey := strings.TrimSpace(params.ObjectKey)

	switch {
	case id == "":
		return nil, ErrMessageIDRequired
	case roomID == "":
		return nil, ErrMessageRoomRequired
	case senderID == "":
		return nil, ErrMessageSenderRequired
	case messageType == "":
		return nil, ErrMessageTypeInvalid
	case messageType == MessageTypeText && content == "":
		return nil, ErrMessageBodyRequired
	case (messageType == MessageTypeImage || messageType == MessageTypeFile) && objectKey == "":
		return nil, ErrMessageObjectKeyRequired
	}

	return &MessageEntity{
		ID:                     id,
		RoomID:                 roomID,
		SenderID:               senderID,
		Message:                content,
		MessageType:            messageType,
		ReplyToMessageID:       strings.TrimSpace(params.ReplyToMessageID),
		ForwardedFromMessageID: strings.TrimSpace(params.ForwardedFromMessageID),
		FileName:               strings.TrimSpace(params.FileName),
		FileSize:               params.FileSize,
		MimeType:               strings.TrimSpace(params.MimeType),
		ObjectKey:              objectKey,
		CreatedAt:              normalizeRoomTime(now),
	}, nil
}

func NewSystemMessage(id, roomID, senderID, body string, now time.Time) (*MessageEntity, error) {
	return NewMessage(id, roomID, senderID, MessageParams{
		Message:     body,
		MessageType: MessageTypeSystem,
	}, now)
}

func NormalizeMessageType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", MessageTypeText:
		return MessageTypeText
	case MessageTypeSystem:
		return MessageTypeSystem
	case MessageTypeImage:
		return MessageTypeImage
	case MessageTypeFile:
		return MessageTypeFile
	default:
		return ""
	}
}

func (m *MessageEntity) Edit(actorID, content string, editedAt time.Time) error {
	if strings.TrimSpace(actorID) != strings.TrimSpace(m.SenderID) {
		return ErrMessageCannotEditOther
	}
	if NormalizeMessageType(m.MessageType) == MessageTypeSystem {
		return ErrMessageCannotEditSystem
	}
	if content = strings.TrimSpace(content); content == "" {
		return ErrMessageBodyRequired
	}

	now := normalizeRoomTime(editedAt)
	m.Message = content
	m.EditedAt = &now
	return nil
}

func (m *MessageEntity) DeleteForEveryone(actorID string, deletedAt time.Time) error {
	if strings.TrimSpace(actorID) != strings.TrimSpace(m.SenderID) {
		return ErrMessageCannotDeleteEveryone
	}

	now := normalizeRoomTime(deletedAt)
	m.Message = ""
	m.DeletedForEveryoneAt = &now
	return nil
}

func (m *MessageEntity) CanBeMarkedBy(accountID string) bool {
	return strings.TrimSpace(accountID) != "" && strings.TrimSpace(accountID) != strings.TrimSpace(m.SenderID)
}
