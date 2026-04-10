package assembly

import (
	appCtx "go-socket/core/context"
	roomprojection "go-socket/core/modules/room/application/messaging"
	roomrepo "go-socket/core/modules/room/infra/persistent/repository"
	"go-socket/core/shared/config"
)

func buildProjectionHandler(cfg *config.Config, appCtx *appCtx.AppContext) (roomprojection.MessageHandler, error) {
	repos := roomrepo.NewRepoImpl(appCtx)
	return roomprojection.NewMessageHandler(cfg, repos.RoomAccountProjectionRepository())
}
