package models

import "time"

type AccountModel struct {
	ID           string     `gorm:"primaryKey"`
	Email        string     `gorm:"not null"`
	Password     string     `gorm:"not null"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
	BannedReason string     `gorm:"column:banned_reason"`
	BannedUntil  *time.Time `gorm:"column:banned_until"`
}

func (AccountModel) TableName() string {
	return "accounts"
}
