// CODE_GENERATOR - do not edit: routing
package http

import (
	"wechat-clone/core/modules/foreign_exchange/application/dto/in"
	"wechat-clone/core/modules/foreign_exchange/application/dto/out"
	"wechat-clone/core/modules/foreign_exchange/transport/http/handler"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(
	routes *gin.RouterGroup,
	createQuote cqrs.Dispatcher[*in.CreateQuoteRequest, *out.CreateQuoteResponse],
) {
	routes.POST("/fx/quotes", httpx.Wrap(handler.NewCreateQuoteHandler(createQuote)))
}
func RegisterPrivateRoutes(_ *gin.RouterGroup) {}
