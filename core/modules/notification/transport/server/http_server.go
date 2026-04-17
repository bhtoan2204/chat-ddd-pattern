// CODE_GENERATOR: registry
package server

import (
	"context"

	"wechat-clone/core/modules/notification/application/dto/in"
	"wechat-clone/core/modules/notification/application/dto/out"
	notificationhttp "wechat-clone/core/modules/notification/transport/http"
	"wechat-clone/core/shared/pkg/cqrs"
	infrahttp "wechat-clone/core/shared/transport/http"

	"github.com/gin-gonic/gin"
)

type notificationHTTPServer struct {
	savePushSubscription cqrs.Dispatcher[*in.SavePushSubscriptionRequest, *out.SavePushSubscriptionResponse]
	listNotification     cqrs.Dispatcher[*in.ListNotificationRequest, *out.ListNotificationResponse]
}

func NewHTTPServer(
	savePushSubscription cqrs.Dispatcher[*in.SavePushSubscriptionRequest, *out.SavePushSubscriptionResponse],
	listNotification cqrs.Dispatcher[*in.ListNotificationRequest, *out.ListNotificationResponse],
) (infrahttp.HTTPServer, error) {
	return &notificationHTTPServer{
		savePushSubscription: savePushSubscription,
		listNotification:     listNotification,
	}, nil
}

func (s *notificationHTTPServer) RegisterPublicRoutes(routes *gin.RouterGroup) {
	notificationhttp.RegisterPublicRoutes(routes)
}

func (s *notificationHTTPServer) RegisterPrivateRoutes(routes *gin.RouterGroup) {
	notificationhttp.RegisterPrivateRoutes(routes, s.savePushSubscription, s.listNotification)
}

func (s *notificationHTTPServer) RegisterSocketRoutes(routes *gin.RouterGroup) {
}

func (s *notificationHTTPServer) Stop(ctx context.Context) error {
	return nil
}
