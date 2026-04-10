package assembly

import (
	appCtx "go-socket/core/context"
	paymentprocessor "go-socket/core/modules/payment/application/projection"
	paymentrepo "go-socket/core/modules/payment/infra/persistent/repository"
	"go-socket/core/shared/config"
	"go-socket/core/shared/pkg/stackErr"
	modruntime "go-socket/core/shared/runtime"
)

func buildProjectionRuntime(cfg *config.Config, appCtx *appCtx.AppContext) (modruntime.Module, error) {
	processor, err := paymentprocessor.NewProcessor(cfg, paymentrepo.NewRepoImpl(appCtx))
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return processor, nil
}
