package model

import "time"

type LedgerEventModel struct {
	ID            string    `gorm:"primaryKey"`
	AggregateID   string    `gorm:"not null;uniqueIndex:idx_ledger_events_agg_type_ver"`
	AggregateType string    `gorm:"not null;uniqueIndex:idx_ledger_events_agg_type_ver"`
	Version       int       `gorm:"not null;uniqueIndex:idx_ledger_events_agg_type_ver"`
	EventName     string    `gorm:"not null;index"`
	EventData     string    `gorm:"type:text;not null"`
	Metadata      string    `gorm:"type:text;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (LedgerEventModel) TableName() string {
	return "ledger_events"
}
