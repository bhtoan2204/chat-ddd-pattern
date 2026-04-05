package models

import "time"

type AccountModel struct {
	ID                string  `gorm:"primaryKey"`
	Email             string  `gorm:"not null;uniqueIndex"`
	Password          string  `gorm:"not null"`
	DisplayName       string  `gorm:"not null"`
	Username          *string `gorm:"uniqueIndex"`
	AvatarObjectKey   *string
	Status            string `gorm:"not null;default:active"`
	EmailVerifiedAt   *time.Time
	LastLoginAt       *time.Time
	PasswordChangedAt *time.Time
	BannedReason      string
	BannedUntil       *time.Time
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

func (AccountModel) TableName() string {
	return "accounts"
}
