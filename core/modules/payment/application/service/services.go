package service

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/payment/domain/repos"
	"go-socket/core/modules/payment/providers"
)

type Services interface {
	ProviderService() ProviderService
}

type services struct {
	providerService ProviderService
}

func NewServices(appCtx *appCtx.AppContext, repos repos.Repos, providerRegistry *providers.ProviderRegistry) Services {
	providerSvc := newProviderService(providerRegistry)
	_ = appCtx
	_ = repos
	return &services{
		providerService: providerSvc,
	}
}

func (s *services) ProviderService() ProviderService {
	return s.providerService
}
