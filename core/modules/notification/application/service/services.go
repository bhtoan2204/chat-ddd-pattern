package service

import (
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/notification/domain/repos"
)

//go:generate mockgen -package=service -destination=services_mock.go -source=services.go
type Services interface {
	EmailVerificationService() EmailVerificationService
	PushDeliveryService() PushDeliveryService
	RealtimeService() RealtimeService
}

type services struct {
	emailVerificationService EmailVerificationService
	pushDeliveryService      PushDeliveryService
	realtimeService          RealtimeService
}

func NewServices(appCtx *appCtx.AppContext, repos repos.Repos) Services {
	emailVerificationService := newEmailVerificationService(appCtx.GetSMTP())
	pushDeliveryService := newPushDeliveryService(repos.PushSubscriptionRepository(), appCtx.GetWebPush())
	realtimeService := newRealtimeService(appCtx)
	return &services{
		emailVerificationService: emailVerificationService,
		pushDeliveryService:      pushDeliveryService,
		realtimeService:          realtimeService,
	}
}

func (s *services) EmailVerificationService() EmailVerificationService {
	return s.emailVerificationService
}

func (s *services) PushDeliveryService() PushDeliveryService {
	return s.pushDeliveryService
}

func (s *services) RealtimeService() RealtimeService {
	return s.realtimeService
}
