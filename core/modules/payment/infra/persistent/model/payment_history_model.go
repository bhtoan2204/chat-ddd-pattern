package model

import "time"

type PaymentHistoryModel struct {
	ID           string    `gorm:"primaryKey;type:uuid"`
	Type         string    `gorm:"type:varchar(50);not null;index"`
	Amount       int64     `gorm:"not null"`
	Balance      int64     `gorm:"not null"`
	SenderID     *string   `gorm:"type:uuid;index"`
	ReceiverID   *string   `gorm:"type:uuid;index"`
	SenderName   *string   `gorm:"type:varchar(255)"`
	ReceiverName *string   `gorm:"type:varchar(255)"`
	Properties   string    `gorm:"type:JSON;not null;default:'{}'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (PaymentHistoryModel) TableName() string {
	return "payment_histories"
}
