// CODE_GENERATOR - do not edit: handler
package handler

import (
	"io"
	"net/http"

	"wechat-clone/core/modules/workflow/application/dto/in"
	"wechat-clone/core/modules/workflow/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type processStripeWebhookHandler struct {
	processStripeWebhook cqrs.Dispatcher[*in.ProcessStripeWebhookRequest, *out.StripeWebhookResponse]
}

func NewProcessStripeWebhookHandler(
	processStripeWebhook cqrs.Dispatcher[*in.ProcessStripeWebhookRequest, *out.StripeWebhookResponse],
) *processStripeWebhookHandler {
	return &processStripeWebhookHandler{
		processStripeWebhook: processStripeWebhook,
	}
}

func (h *processStripeWebhookHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.ProcessStripeWebhookRequest
	request.Signature = c.GetHeader("Stripe-Signature")

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorw("Read request body failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
		return nil, nil
	}
	request.Payload = string(payload)

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.processStripeWebhook.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("ProcessStripeWebhook failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
