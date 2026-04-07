package service

import "go-socket/core/modules/payment/providers"

type ProviderService interface {
	Get(name string) (providers.PaymentProvider, error)
}

type providerService struct {
	providerRegistry *providers.ProviderRegistry
}

func newProviderService(providerRegistry *providers.ProviderRegistry) ProviderService {
	return providerService{providerRegistry: providerRegistry}
}

func (s providerService) Get(name string) (providers.PaymentProvider, error) {
	return s.providerRegistry.Get(name)
}
