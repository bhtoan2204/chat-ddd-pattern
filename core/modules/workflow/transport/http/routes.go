// CODE_GENERATOR - do not edit: routing
package http

import (
	"wechat-clone/core/modules/workflow/application/dto/in"
	"wechat-clone/core/modules/workflow/application/dto/out"
	"wechat-clone/core/modules/workflow/transport/http/handler"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(
	routes *gin.RouterGroup,
	processStripeWebhook cqrs.Dispatcher[*in.ProcessStripeWebhookRequest, *out.StripeWebhookResponse],
) {
	routes.POST("/payment/webhooks/stripe", httpx.Wrap(handler.NewProcessStripeWebhookHandler(processStripeWebhook)))
}
func RegisterPrivateRoutes(
	routes *gin.RouterGroup,
	createStripeTopUp cqrs.Dispatcher[*in.CreateStripeTopUpRequest, *out.StripeTopUpResponse],
) {
	routes.POST("/payment/intents", httpx.Wrap(handler.NewCreateStripeTopUpHandler(createStripeTopUp)))
}
