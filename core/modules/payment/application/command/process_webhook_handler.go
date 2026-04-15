package command

import (
	"context"

	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/infra/lock"
	"go-socket/core/modules/payment/application/dto/in"
	"go-socket/core/modules/payment/application/dto/out"
	paymentservice "go-socket/core/modules/payment/application/service"
	"go-socket/core/shared/pkg/cqrs"
)

type processWebhookHandler struct {
	locker                lock.Lock
	paymentCommandService paymentservice.PaymentCommandService
}

func NewProcessWebhook(appContext *appCtx.AppContext, services paymentservice.Services) cqrs.Handler[*in.ProcessWebhookRequest, *out.ProcessWebhookResponse] {
	return &processWebhookHandler{
		paymentCommandService: services.PaymentCommandService(),
		locker:                appContext.Locker(),
	}
}

func (u *processWebhookHandler) Handle(ctx context.Context, req *in.ProcessWebhookRequest) (*out.ProcessWebhookResponse, error) {
	return u.paymentCommandService.ProcessWebhook(ctx, req)
}
