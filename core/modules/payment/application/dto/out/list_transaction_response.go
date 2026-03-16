package out

import "time"

type ListTransactionResponse struct {
	Page    int64               `json:"page" form:"page"`
	Limit   int64               `json:"limit" form:"limit"`
	Balance int64               `json:"balance" form:"balance"`
	Record  []TransactionRecord `json:"records" form:"record"`
}

type TransactionRecord struct {
	Type       string    `json:"type" form:"type"`
	Amount     int64     `json:"amount" form:"amount"`
	Balance    int64     `json:"balance" form:"balance"`
	Date       time.Time `json:"date" form:"date"`
	Sender     string    `json:"sender" form:"sender"`
	SenderID   string    `json:"sender_id" form:"sender_id"`
	Receiver   string    `json:"receiver" form:"receiver"`
	ReceiverID string    `json:"receiver_id" form:"receiver_id"`
}
