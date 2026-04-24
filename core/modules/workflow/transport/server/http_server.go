// CODE_GENERATOR: registry
package server

import (
	"context"

	"wechat-clone/core/modules/workflow/application/dto/in"
	"wechat-clone/core/modules/workflow/application/dto/out"
	workflowhttp "wechat-clone/core/modules/workflow/transport/http"
	"wechat-clone/core/shared/pkg/cqrs"
	infrahttp "wechat-clone/core/shared/transport/http"

	"github.com/gin-gonic/gin"
)

type workflowHTTPServer struct {
	createStripeTopUp    cqrs.Dispatcher[*in.CreateStripeTopUpRequest, *out.StripeTopUpResponse]
	processStripeWebhook cqrs.Dispatcher[*in.ProcessStripeWebhookRequest, *out.StripeWebhookResponse]
}

func NewHTTPServer(
	createStripeTopUp cqrs.Dispatcher[*in.CreateStripeTopUpRequest, *out.StripeTopUpResponse],
	processStripeWebhook cqrs.Dispatcher[*in.ProcessStripeWebhookRequest, *out.StripeWebhookResponse],
) (infrahttp.HTTPServer, error) {
	return &workflowHTTPServer{
		createStripeTopUp:    createStripeTopUp,
		processStripeWebhook: processStripeWebhook,
	}, nil
}

func (s *workflowHTTPServer) RegisterPublicRoutes(routes *gin.RouterGroup) {
	workflowhttp.RegisterPublicRoutes(routes, s.processStripeWebhook)
}

func (s *workflowHTTPServer) RegisterPrivateRoutes(routes *gin.RouterGroup) {
	workflowhttp.RegisterPrivateRoutes(routes, s.createStripeTopUp)
}

func (s *workflowHTTPServer) RegisterSocketRoutes(routes *gin.RouterGroup) {
}

func (s *workflowHTTPServer) Stop(ctx context.Context) error {
	return nil
}
