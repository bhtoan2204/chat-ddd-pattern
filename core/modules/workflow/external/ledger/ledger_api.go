package ledger

import (
	"context"
	"time"

	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/infra/tracing"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"
	ledgerv1 "wechat-clone/core/shared/transport/grpc/gen/ledger/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type LedgerGrpc interface {
	ledgerv1.LedgerPaymentServiceClient
	Close() error
}

type ledgerGrpc struct {
	cc         *grpc.ClientConn
	grpcClient ledgerv1.LedgerPaymentServiceClient
}

func New(ctx context.Context, cfg *config.Config) (LedgerGrpc, error) {
	log := logging.FromContext(ctx)
	cc, err := grpc.NewClient(
		cfg.GrpcConfig.LedgerGRPCBaseURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(),
		grpc.WithStatsHandler(tracing.GrpcStatsHandler()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		log.Warnw("Could not setup Ledger Service GRPC connection", zap.String("connection string", cfg.GrpcConfig.LedgerGRPCBaseURL), zap.Error(err))
		return nil, stackErr.Error(err)
	}

	return ledgerGrpc{
		cc:         cc,
		grpcClient: ledgerv1.NewLedgerPaymentServiceClient(cc),
	}, nil
}

func (l ledgerGrpc) ApplyPaymentEvent(ctx context.Context, req *ledgerv1.ApplyPaymentEventRequest, opts ...grpc.CallOption) (*ledgerv1.ApplyPaymentEventResponse, error) {
	response, err := l.grpcClient.ApplyPaymentEvent(ctx, req, opts...)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return response, nil
}

func (l ledgerGrpc) Close() error {
	if l.cc == nil {
		return nil
	}
	return stackErr.Error(l.cc.Close())
}
