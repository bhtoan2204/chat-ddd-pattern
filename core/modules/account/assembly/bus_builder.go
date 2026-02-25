package assembly

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/command"
	"go-socket/core/modules/account/application/query"
	accountrepo "go-socket/core/modules/account/infra/persistent/repository"
)

type Buses struct {
	Command command.Bus
	Query   query.Bus
}

func BuildBuses(appCtx *appCtx.AppContext) Buses {
	accountRepos := accountrepo.NewRepoImpl(appCtx)
	loginHandler := command.NewLoginHandler(appCtx, accountRepos)
	registerHandler := command.NewRegisterHandler(appCtx, accountRepos)
	logoutHandler := command.NewLogoutHandler()
	getProfileHandler := query.NewGetProfileHandler(accountRepos)
	return Buses{
		Command: command.NewBus(loginHandler, registerHandler, logoutHandler),
		Query:   query.NewBus(getProfileHandler),
	}
}
