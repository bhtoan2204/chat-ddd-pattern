package providers

import (
	"context"
	"errors"

	"wechat-clone/core/modules/payment/domain/entity"
)

var (
	ErrProviderNotFound        = errors.New("provider not found")
	ErrInvalidWebhookSignature = errors.New("invalid webhook signature")
	ErrWebhookEventIgnored     = errors.New("webhook event ignored")
)

type CreatePaymentRequest struct {
	TransactionID string
	Amount        int64
	Currency      string
	Metadata      map[string]string
}

type RefundPaymentRequest struct {
	TransactionID string
	ExternalRef   string
	Amount        int64
	Currency      string
	Reason        string
}

type CreatePaymentResponse struct {
	Provider      string `json:"provider"`
	TransactionID string `json:"transaction_id"`
	ExternalRef   string `json:"external_ref,omitempty"`
	Status        string `json:"status"`
	CheckoutURL   string `json:"checkout_url,omitempty"`
}

type RefundPaymentResponse struct {
	Provider      string `json:"provider"`
	TransactionID string `json:"transaction_id"`
	ExternalRef   string `json:"external_ref,omitempty"`
	Status        string `json:"status"`
	Amount        int64  `json:"amount,omitempty"`
	Currency      string `json:"currency,omitempty"`
}

type WebhookEvent struct {
	Provider   string
	EventID    string
	EventType  string
	Attributes map[string]string
}

type PaymentResult struct {
	TransactionID string
	EventID       string
	EventType     string
	Status        string
	Amount        int64
	Currency      string
	ExternalRef   string
}

//go:generate mockgen -package=providers -destination=provider_mock.go -source=provider.go
type PaymentProvider interface {
	Name() string
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error)
	CreateWithdrawal(ctx context.Context, intent *entity.PaymentIntent, metadata map[string]string) (*CreatePaymentResponse, error)
	RefundPayment(ctx context.Context, req RefundPaymentRequest) (*RefundPaymentResponse, error)
	VerifyWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error)
	ParseEvent(ctx context.Context, event *WebhookEvent) (*PaymentResult, error)
}
