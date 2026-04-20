package models

import "time"

type MessageReadModel struct {
	ID                     string  `gorm:"primaryKey"`
	RoomID                 string  `gorm:"not null;index"`
	SenderID               string  `gorm:"not null;index"`
	Message                string  `gorm:"type:text;not null"`
	MessageType            string  `gorm:"type:varchar(50);default:'text';not null"`
	MentionsJSON           string  `gorm:"type:text;not null;default:'[]'"`
	MentionAll             bool    `gorm:"type:number(1);default:0;not null"`
	ReplyToMessageID       *string `gorm:"index"`
	ForwardedFromMessageID *string `gorm:"index"`
	FileName               *string `gorm:"type:varchar(1024)"`
	FileSize               *int64
	MimeType               *string `gorm:"type:varchar(255)"`
	ObjectKey              *string `gorm:"type:varchar(2048)"`
	EditedAt               *time.Time
	DeletedForEveryoneAt   *time.Time
	CreatedAt              time.Time `gorm:"autoCreateTime"`
}

func (MessageReadModel) TableName() string {
	return "message_read_models"
}
