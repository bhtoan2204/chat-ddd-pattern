// CODE_GENERATOR: assembly
package assembly

import (
	"context"
	appCtx "go-socket/core/context"
	infrahttp "go-socket/core/shared/transport/http"
)

func BuildHTTPServer(ctx context.Context, appContext *appCtx.AppContext) (infrahttp.HTTPServer, error) {
	return buildHTTPServer(ctx, appContext)
}
