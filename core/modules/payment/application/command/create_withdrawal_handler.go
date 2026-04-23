package command

import (
	"context"

	"wechat-clone/core/modules/payment/application/dto/in"
	"wechat-clone/core/modules/payment/application/dto/out"
	paymentservice "wechat-clone/core/modules/payment/application/service"
	"wechat-clone/core/shared/pkg/cqrs"
)

type createWithdrawalHandler struct {
	paymentCommandService paymentservice.PaymentCommandService
}

func NewCreateWithdrawal(paymentCommandService paymentservice.PaymentCommandService) cqrs.Handler[*in.CreateWithdrawalRequest, *out.CreateWithdrawalResponse] {
	return &createWithdrawalHandler{
		paymentCommandService: paymentCommandService,
	}
}

func (u *createWithdrawalHandler) Handle(ctx context.Context, req *in.CreateWithdrawalRequest) (*out.CreateWithdrawalResponse, error) {
	return u.paymentCommandService.CreateWithdrawal(ctx, req)
}
