package service

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/ledger/domain/repos"
)

type Services interface {
}

type services struct {
}

func NewServices(appCtx *appCtx.AppContext, repos repos.Repos) Services {
	return &services{}
}
