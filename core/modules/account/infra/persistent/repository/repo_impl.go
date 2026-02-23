package repos

import (
	appCtx "go-socket/core/context"
	accountrepos "go-socket/core/modules/account/domain/repos"
)

type repoImpl struct {
	accountRepo accountrepos.AccountRepository
}

func NewRepoImpl(appCtx *appCtx.AppContext) accountrepos.Repos {
	accountRepo := NewAccountRepoImpl(appCtx.GetDB(), appCtx.GetCache())
	return &repoImpl{accountRepo: accountRepo}
}

func (r *repoImpl) AccountRepository() accountrepos.AccountRepository {
	return r.accountRepo
}
