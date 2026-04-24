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
	PaymentNotificationService() PaymentNotificationService
}

type services struct {
	emailVerificationService EmailVerificationService
	pushDeliveryService      PushDeliveryService
	realtimeService          RealtimeService
	paymentNotification      PaymentNotificationService
}

func NewServices(appCtx *appCtx.AppContext, repos repos.Repos) Services {
	emailVerificationService := newEmailVerificationService(appCtx.GetSMTP())
	pushDeliveryService := newPushDeliveryService(repos.PushSubscriptionRepository(), appCtx.GetWebPush())
	realtimeService := newRealtimeService(appCtx)
	paymentNotification := newPaymentNotificationService(repos, realtimeService, pushDeliveryService)
	return &services{
		emailVerificationService: emailVerificationService,
		pushDeliveryService:      pushDeliveryService,
		realtimeService:          realtimeService,
		paymentNotification:      paymentNotification,
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

func (s *services) PaymentNotificationService() PaymentNotificationService {
	return s.paymentNotification
}
