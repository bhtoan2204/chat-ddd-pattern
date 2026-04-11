package service

import (
	repos "go-socket/core/modules/payment/domain/repos"
	domainservice "go-socket/core/modules/payment/domain/service"
)

type Services interface {
	PaymentCommandService() PaymentCommandService
}

type services struct {
	paymentCommandService PaymentCommandService
}

func NewServices(baseRepo repos.Repos, providerRegistry domainservice.PaymentProviderRegistry) Services {
	paymentCommandService := NewPaymentCommandService(baseRepo, providerRegistry)
	return &services{
		paymentCommandService: paymentCommandService,
	}
}

func (s *services) PaymentCommandService() PaymentCommandService {
	return s.paymentCommandService
}
