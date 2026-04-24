package command

import (
	"context"

	"wechat-clone/core/modules/workflow/application/dto/in"
	"wechat-clone/core/modules/workflow/application/dto/out"
	workflowservice "wechat-clone/core/modules/workflow/application/service"
	"wechat-clone/core/shared/pkg/cqrs"
)

type createStripeTopUpHandler struct {
	service workflowservice.StripeTopUpWorkflowService
}

func NewCreateStripeTopUp(service workflowservice.StripeTopUpWorkflowService) cqrs.Handler[*in.CreateStripeTopUpRequest, *out.StripeTopUpResponse] {
	return &createStripeTopUpHandler{
		service: service,
	}
}

func (h *createStripeTopUpHandler) Handle(ctx context.Context, req *in.CreateStripeTopUpRequest) (*out.StripeTopUpResponse, error) {
	return h.service.CreateStripeTopUp(ctx, req)
}
