package server

import (
	"context"
	"go-socket/core/modules/payment/application/command"
	"go-socket/core/modules/payment/application/query"
	paymenthttp "go-socket/core/modules/payment/transport/http"
	infrahttp "go-socket/core/shared/transport/http"

	"github.com/gin-gonic/gin"
)

type paymentHTTPServer struct {
	commandBus command.Bus
	queryBus   query.Bus
}

func NewHTTPServer(commandBus command.Bus, queryBus query.Bus) (infrahttp.HTTPServer, error) {
	return &paymentHTTPServer{commandBus: commandBus, queryBus: queryBus}, nil
}

func (s *paymentHTTPServer) RegisterPublicRoutes(_ *gin.RouterGroup) {}

func (s *paymentHTTPServer) RegisterPrivateRoutes(routes *gin.RouterGroup) {
	paymenthttp.RegisterPrivateRoutes(routes, s.commandBus, s.queryBus)
}

func (s *paymentHTTPServer) Stop(_ context.Context) error {
	return nil
}
