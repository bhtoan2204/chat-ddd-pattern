package model

import "time"

type LedgerPostedTransactionModel struct {
	ID                    string    `gorm:"primaryKey"`
	AggregateID           string    `gorm:"not null;uniqueIndex:idx_ledger_posted_tx_agg_type_tx"`
	AggregateType         string    `gorm:"not null;uniqueIndex:idx_ledger_posted_tx_agg_type_tx"`
	TransactionID         string    `gorm:"not null;uniqueIndex:idx_ledger_posted_tx_agg_type_tx"`
	ReferenceType         string    `gorm:"not null"`
	ReferenceID           string    `gorm:"not null"`
	CounterpartyAccountID string    `gorm:"not null"`
	Currency              string    `gorm:"not null"`
	AmountDelta           int64     `gorm:"not null"`
	BookedAt              time.Time `gorm:"not null"`
	CreatedAt             time.Time `gorm:"autoCreateTime"`
}

func (LedgerPostedTransactionModel) TableName() string {
	return "ledger_posted_transactions"
}
