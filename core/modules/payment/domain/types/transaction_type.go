package types

type TransactionType string

const (
	TransactionTypeDeposited   TransactionType = "deposited"
	TransactionTypeWithdrawn   TransactionType = "withdrawn"
	TransactionTypeTransferred TransactionType = "transferred"
	TransactionTypeRefunded    TransactionType = "refunded"
)

func (t TransactionType) String() string {
	return string(t)
}
