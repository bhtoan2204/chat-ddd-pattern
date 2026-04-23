package assembly

import (
	"time"

	appCtx "wechat-clone/core/context"
	paymentprocessor "wechat-clone/core/modules/payment/application/processor"
	paymentservice "wechat-clone/core/modules/payment/application/service"
	paymentrepo "wechat-clone/core/modules/payment/infra/persistent/repository"
	provideradapter "wechat-clone/core/modules/payment/infra/provider"
	"wechat-clone/core/modules/payment/providers"
	mockprovider "wechat-clone/core/modules/payment/providers/mock"
	stripeprovider "wechat-clone/core/modules/payment/providers/stripe"
	"wechat-clone/core/shared/config"
	modruntime "wechat-clone/core/shared/runtime"
)

func buildMessagingRuntime(cfg *config.Config, appContext *appCtx.AppContext) (modruntime.Module, error) {
	paymentRepos := paymentrepo.NewRepoImpl(appContext)
	providerRegistry := providers.NewProviderRegistry()
	providerRegistry.Register(mockprovider.NewProvider(appContext.GetConfig().LedgerConfig.MockWebhookSecret))
	if stripe := stripeprovider.NewProvider(appContext.GetConfig().LedgerConfig.Stripe); stripe.Enabled() {
		providerRegistry.Register(stripe)
	}

	commandService := paymentservice.NewPaymentCommandService(appContext, paymentRepos, provideradapter.NewPaymentProviderRegistry(providerRegistry))
	interval := time.Duration(cfg.LedgerConfig.Stripe.WithdrawalPollIntervalSecond) * time.Second
	return paymentprocessor.NewProcessor(commandService, interval), nil
}
