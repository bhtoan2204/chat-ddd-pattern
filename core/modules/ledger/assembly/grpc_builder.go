// CODE_GENERATOR: assembly
package assembly

import (
	"context"
	appCtx "wechat-clone/core/context"
	infragrpc "wechat-clone/core/shared/transport/grpc"
)

func BuildGRPCServer(ctx context.Context, appContext *appCtx.AppContext) (infragrpc.GRPCServer, error) {
	return buildGRPCServer(ctx, appContext)
}
