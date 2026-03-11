package server

import (
	"context"
	notificationcommand "go-socket/core/modules/notification/application/command"
	notificationquery "go-socket/core/modules/notification/application/query"
	notificationhttp "go-socket/core/modules/notification/transport/http"

	"github.com/gin-gonic/gin"
)

type HTTPServer interface {
	RegisterPublicRoutes(routes *gin.RouterGroup)
	RegisterPrivateRoutes(routes *gin.RouterGroup)
	Stop(ctx context.Context) error
}

type notificationHTTPServer struct {
	commandBus notificationcommand.Bus
	queryBus   notificationquery.Bus
}

func NewHTTPServer(commandBus notificationcommand.Bus, queryBus notificationquery.Bus) (HTTPServer, error) {
	return &notificationHTTPServer{commandBus: commandBus, queryBus: queryBus}, nil
}

func (s *notificationHTTPServer) RegisterPublicRoutes(_ *gin.RouterGroup) {}

func (s *notificationHTTPServer) RegisterPrivateRoutes(routes *gin.RouterGroup) {
	notificationhttp.RegisterPrivateRoutes(routes, s.commandBus, s.queryBus)
}

func (s *notificationHTTPServer) Stop(_ context.Context) error {
	return nil
}
