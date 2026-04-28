package command

import (
	"context"

	"wechat-clone/core/modules/payment/application/dto/in"
	"wechat-clone/core/modules/payment/application/dto/out"
	paymentservice "wechat-clone/core/modules/payment/application/service"
	"wechat-clone/core/shared/pkg/cqrs"
)

type refundPaymentHandler struct {
	paymentCommandService paymentservice.PaymentCommandService
}

func NewRefundPayment(paymentCommandService paymentservice.PaymentCommandService) cqrs.Handler[*in.RefundPaymentRequest, *out.RefundPaymentResponse] {
	return &refundPaymentHandler{
		paymentCommandService: paymentCommandService,
	}
}

func (u *refundPaymentHandler) Handle(ctx context.Context, req *in.RefundPaymentRequest) (*out.RefundPaymentResponse, error) {
	return u.paymentCommandService.RefundPayment(ctx, req)
}
