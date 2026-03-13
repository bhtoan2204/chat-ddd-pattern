package model

import "time"

type PaymentAggregateModel struct {
	ID            string    `gorm:"primaryKey"`
	AggregateID   string    `gorm:"not null;uniqueIndex:idx_agg_ver"`
	AggregateType string    `gorm:"not null"`
	Version       int       `gorm:"not null;uniqueIndex:idx_agg_ver"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (PaymentAggregateModel) TableName() string {
	return "payment_aggregates"
}
