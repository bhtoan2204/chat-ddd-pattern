package model

import "time"

type PaymentBalanceSnapshotModel struct {
	ID          string    `gorm:"primaryKey"`
	AggregateID string    `gorm:"not null;uniqueIndex:idx_snap_ver"`
	Version     int       `gorm:"not null;uniqueIndex:idx_snap_ver"`
	State       string    `gorm:"type:JSON;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (PaymentBalanceSnapshotModel) TableName() string {
	return "payment_balance_snapshots"
}
