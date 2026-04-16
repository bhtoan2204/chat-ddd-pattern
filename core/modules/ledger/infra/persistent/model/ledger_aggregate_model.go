package model

import "time"

type LedgerAggregateModel struct {
	ID            string    `gorm:"primaryKey"`
	AggregateID   string    `gorm:"not null;uniqueIndex:idx_ledger_aggregates_agg_type_id"`
	AggregateType string    `gorm:"not null;uniqueIndex:idx_ledger_aggregates_agg_type_id"`
	Version       int       `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (LedgerAggregateModel) TableName() string {
	return "ledger_aggregates"
}
