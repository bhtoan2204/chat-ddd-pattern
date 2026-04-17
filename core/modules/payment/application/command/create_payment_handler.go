package command

import (
	"context"

	"wechat-clone/core/modules/payment/application/dto/in"
	"wechat-clone/core/modules/payment/application/dto/out"
	paymentservice "wechat-clone/core/modules/payment/application/service"
	"wechat-clone/core/shared/pkg/cqrs"
)

type createPaymentHandler struct {
	paymentCommandService paymentservice.PaymentCommandService
}

func NewCreatePayment(paymentCommandService paymentservice.PaymentCommandService) cqrs.Handler[*in.CreatePaymentRequest, *out.CreatePaymentResponse] {
	return &createPaymentHandler{
		paymentCommandService: paymentCommandService,
	}
}

func (u *createPaymentHandler) Handle(ctx context.Context, req *in.CreatePaymentRequest) (*out.CreatePaymentResponse, error) {
	return u.paymentCommandService.CreatePayment(ctx, req)
}
