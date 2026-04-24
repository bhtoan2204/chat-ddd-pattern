package command

import (
	"context"

	"wechat-clone/core/modules/workflow/application/dto/in"
	"wechat-clone/core/modules/workflow/application/dto/out"
	workflowservice "wechat-clone/core/modules/workflow/application/service"
	"wechat-clone/core/shared/pkg/cqrs"
)

type processStripeWebhookHandler struct {
	service workflowservice.StripeTopUpWorkflowService
}

func NewProcessStripeWebhook(service workflowservice.StripeTopUpWorkflowService) cqrs.Handler[*in.ProcessStripeWebhookRequest, *out.StripeWebhookResponse] {
	return &processStripeWebhookHandler{
		service: service,
	}
}

func (h *processStripeWebhookHandler) Handle(ctx context.Context, req *in.ProcessStripeWebhookRequest) (*out.StripeWebhookResponse, error) {
	return h.service.ProcessStripeWebhook(ctx, req)
}
