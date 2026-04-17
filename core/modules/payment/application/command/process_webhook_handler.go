package command

import (
	"context"

	"wechat-clone/core/modules/payment/application/dto/in"
	"wechat-clone/core/modules/payment/application/dto/out"
	paymentservice "wechat-clone/core/modules/payment/application/service"
	"wechat-clone/core/shared/pkg/cqrs"
)

type processWebhookHandler struct {
	paymentCommandService paymentservice.PaymentCommandService
}

func NewProcessWebhook(paymentCommandService paymentservice.PaymentCommandService) cqrs.Handler[*in.ProcessWebhookRequest, *out.ProcessWebhookResponse] {
	return &processWebhookHandler{
		paymentCommandService: paymentCommandService,
	}
}

func (u *processWebhookHandler) Handle(ctx context.Context, req *in.ProcessWebhookRequest) (*out.ProcessWebhookResponse, error) {
	return u.paymentCommandService.ProcessWebhook(ctx, req)
}
