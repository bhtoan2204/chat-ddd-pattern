package grpc

import (
	"context"
	"fmt"

	appCtx "wechat-clone/core/context"
	"wechat-clone/core/shared/pkg/stackErr"
)

type ModuleBuilder func(ctx context.Context, appCtx *appCtx.AppContext) (GRPCServer, error)

func BuildModuleServers(ctx context.Context, appCtx *appCtx.AppContext, builders ...ModuleBuilder) ([]GRPCServer, error) {
	servers := make([]GRPCServer, 0, len(builders))
	for idx, builder := range builders {
		if builder == nil {
			return nil, stackErr.Error(fmt.Errorf("grpc module builder %d is nil", idx))
		}
		server, err := builder(ctx, appCtx)
		if err != nil {
			return nil, stackErr.Error(fmt.Errorf("build grpc module server %d failed: %w", idx, err))
		}
		if server == nil {
			continue
		}
		servers = append(servers, server)
	}
	return servers, nil
}
