// CODE_GENERATOR: assembly
package assembly

import (
	appCtx "go-socket/core/context"
	"go-socket/core/shared/config"
	modruntime "go-socket/core/shared/runtime"
)

func BuildProjectionRuntime(cfg *config.Config, appContext *appCtx.AppContext) (modruntime.Module, error) {
	return buildProjectionRuntime(cfg, appContext)
}
