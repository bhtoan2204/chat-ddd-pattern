// CODE_GENERATOR: registry
package server

import (
	"context"

	"wechat-clone/core/modules/foreign_exchange/application/dto/in"
	"wechat-clone/core/modules/foreign_exchange/application/dto/out"
	foreign_exchangehttp "wechat-clone/core/modules/foreign_exchange/transport/http"
	"wechat-clone/core/shared/pkg/cqrs"
	infrahttp "wechat-clone/core/shared/transport/http"

	"github.com/gin-gonic/gin"
)

type foreign_exchangeHTTPServer struct {
	createQuote cqrs.Dispatcher[*in.CreateQuoteRequest, *out.CreateQuoteResponse]
}

func NewHTTPServer(
	createQuote cqrs.Dispatcher[*in.CreateQuoteRequest, *out.CreateQuoteResponse],
) (infrahttp.HTTPServer, error) {
	return &foreign_exchangeHTTPServer{
		createQuote: createQuote,
	}, nil
}

func (s *foreign_exchangeHTTPServer) RegisterPublicRoutes(routes *gin.RouterGroup) {
	foreign_exchangehttp.RegisterPublicRoutes(routes, s.createQuote)
}

func (s *foreign_exchangeHTTPServer) RegisterPrivateRoutes(routes *gin.RouterGroup) {
	foreign_exchangehttp.RegisterPrivateRoutes(routes)
}

func (s *foreign_exchangeHTTPServer) RegisterSocketRoutes(routes *gin.RouterGroup) {
}

func (s *foreign_exchangeHTTPServer) Stop(ctx context.Context) error {
	return nil
}
