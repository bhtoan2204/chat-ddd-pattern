package assembly

import (
	"context"

	appCtx "wechat-clone/core/context"
	ledgerapp "wechat-clone/core/modules/ledger/application/service"
	ledgerrepo "wechat-clone/core/modules/ledger/infra/persistent/repository"
	ledgergrpc "wechat-clone/core/modules/ledger/transport/grpc"
	infragrpc "wechat-clone/core/shared/transport/grpc"
	ledgerv1 "wechat-clone/core/shared/transport/grpc/gen/ledger/v1"

	"google.golang.org/grpc"
)

type ledgerGRPCRegistrar struct {
	server ledgerv1.LedgerPaymentServiceServer
}

func buildGRPCServer(_ context.Context, appContext *appCtx.AppContext) (infragrpc.GRPCServer, error) {
	ledgerRepos := ledgerrepo.NewRepoImpl(appContext)
	paymentEventService := ledgerapp.NewPaymentEventService(ledgerRepos, appContext.GetConfig().LedgerConfig.Stripe.FeeAccountID)

	return &ledgerGRPCRegistrar{
		server: ledgergrpc.NewPaymentServer(
			appContext.Locker(),
			paymentEventService,
			appContext.GetConfig().LedgerConfig.Stripe.FeeAccountID,
		),
	}, nil
}

func (r *ledgerGRPCRegistrar) Register(server grpc.ServiceRegistrar) {
	ledgerv1.RegisterLedgerPaymentServiceServer(server, r.server)
}
