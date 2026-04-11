package command

import (
	"context"

	"go-socket/core/modules/payment/application/dto/in"
	"go-socket/core/modules/payment/application/dto/out"
	paymentservice "go-socket/core/modules/payment/application/service"
	"go-socket/core/shared/pkg/cqrs"
)

type createPaymentHandler struct {
	paymentCommandService paymentservice.PaymentCommandService
}

func NewCreatePayment(services paymentservice.Services) cqrs.Handler[*in.CreatePaymentRequest, *out.CreatePaymentResponse] {
	return &createPaymentHandler{
		paymentCommandService: services.PaymentCommandService(),
	}
}

func (u *createPaymentHandler) Handle(ctx context.Context, req *in.CreatePaymentRequest) (*out.CreatePaymentResponse, error) {
	return u.paymentCommandService.CreatePayment(ctx, req)
}
