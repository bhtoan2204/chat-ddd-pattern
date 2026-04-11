package service

import (
	"context"

	"go-socket/core/modules/payment/domain/entity"
)

type PaymentProvider interface {
	Name() string
	CreatePayment(ctx context.Context, intent *entity.PaymentIntent, metadata map[string]string) (*PaymentCreation, error)
	ParseWebhook(ctx context.Context, payload []byte, signature string) (*PaymentWebhook, error)
}

type PaymentProviderRegistry interface {
	Get(name string) (PaymentProvider, error)
}

type PaymentCreation struct {
	Provider    string
	Result      entity.PaymentProviderResult
	CheckoutURL string
}

type PaymentWebhook struct {
	Provider string
	Result   entity.PaymentProviderResult
}
