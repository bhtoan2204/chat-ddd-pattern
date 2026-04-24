package assembly

import (
	"context"

	appCtx "wechat-clone/core/context"
	paymentcommand "wechat-clone/core/modules/payment/application/command"
	paymentservice "wechat-clone/core/modules/payment/application/service"
	paymentrepo "wechat-clone/core/modules/payment/infra/persistent/repository"
	provideradapter "wechat-clone/core/modules/payment/infra/provider"
	"wechat-clone/core/modules/payment/providers"
	mockprovider "wechat-clone/core/modules/payment/providers/mock"
	stripeprovider "wechat-clone/core/modules/payment/providers/stripe"
	paymentgrpc "wechat-clone/core/modules/payment/transport/grpc"
	"wechat-clone/core/shared/pkg/cqrs"
	infragrpc "wechat-clone/core/shared/transport/grpc"
	paymentv1 "wechat-clone/core/shared/transport/grpc/gen/payment/v1"

	"google.golang.org/grpc"
)

type paymentGRPCRegistrar struct {
	server paymentv1.PaymentServiceServer
}

func buildGRPCServer(_ context.Context, appContext *appCtx.AppContext) (infragrpc.GRPCServer, error) {
	paymentRepos := paymentrepo.NewRepoImpl(appContext)
	providerRegistry := providers.NewProviderRegistry()
	providerRegistry.Register(mockprovider.NewProvider(appContext.GetConfig().LedgerConfig.MockWebhookSecret))
	if stripe := stripeprovider.NewProvider(appContext.GetConfig().LedgerConfig.Stripe); stripe.Enabled() {
		providerRegistry.Register(stripe)
	}
	paymentCommandService := paymentservice.NewPaymentCommandService(appContext, paymentRepos, provideradapter.NewPaymentProviderRegistry(providerRegistry))

	createPayment := cqrs.NewDispatcher(paymentcommand.NewCreatePayment(paymentCommandService))
	processWebhook := cqrs.NewDispatcher(paymentcommand.NewProcessWebhook(paymentCommandService))

	return &paymentGRPCRegistrar{
		server: paymentgrpc.NewServer(createPayment, processWebhook),
	}, nil
}

func (r *paymentGRPCRegistrar) Register(server grpc.ServiceRegistrar) {
	paymentv1.RegisterPaymentServiceServer(server, r.server)
}
