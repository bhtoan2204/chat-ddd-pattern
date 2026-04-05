package assembly

import (
	appCtx "go-socket/core/context"
	paymentprocessor "go-socket/core/modules/payment/application/projection"
	"go-socket/core/shared/config"
	"go-socket/core/shared/pkg/stackErr"
)

func BuildProcessors(cfg *config.Config, appCtx *appCtx.AppContext) (paymentprocessor.Processor, error) {
	processor, err := paymentprocessor.NewProcessor(cfg, appCtx)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return processor, nil
}
