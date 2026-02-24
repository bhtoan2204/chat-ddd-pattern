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
	loginUseCase := command.NewLoginUseCase(appCtx, accountRepos)
	registerUseCase := command.NewRegisterUseCase(appCtx, accountRepos)
	logoutUseCase := command.NewLogoutUseCase()
	getProfileUseCase := query.NewGetProfileUseCase(accountRepos)
	return Buses{
		Command: command.NewBus(loginUseCase, registerUseCase, logoutUseCase),
		Query:   query.NewBus(getProfileUseCase),
	}
}
