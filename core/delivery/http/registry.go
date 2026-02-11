// CODE_GENERATOR: registry
package http

import (
	"go-socket/config"
	"go-socket/core/application/usecase"
	"go-socket/core/delivery/http/handler"
)

func BuildRegistry(config *config.Config, usecase usecase.Usecase) map[string]routingConfig {
	return map[string]routingConfig{
		"POST:/api/v1/auth/login": {
			handler: handler.NewLoginHandler(usecase),
		},
		"POST:/api/v1/auth/register": {
			handler: handler.NewRegisterHandler(usecase),
		},
		"POST:/api/v1/auth/logout": {
			handler: handler.NewLogoutHandler(usecase),
		},
		"GET:/api/v1/auth/profile": {
			handler: handler.NewGetProfileHandler(usecase),
		},
		"POST:/api/v1/room/create": {
			handler: handler.NewCreateRoomHandler(usecase),
		},
		"GET:/api/v1/room/list": {
			handler: handler.NewListRoomsHandler(usecase),
		},
		"GET:/api/v1/room/get": {
			handler: handler.NewGetRoomHandler(usecase),
		},
		"PUT:/api/v1/room/update": {
			handler: handler.NewUpdateRoomHandler(usecase),
		},
		"DELETE:/api/v1/room/delete": {
			handler: handler.NewDeleteRoomHandler(usecase),
		},
		"POST:/api/v1/message/create": {
			handler: handler.NewCreateMessageHandler(usecase),
		},
	}
}
