package service

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/domain/repos"
)

type Services interface {
	AccountService() AccountService
	EmailVerificationService() EmailVerificationService
}

type services struct {
	accountService           AccountService
	emailVerificationService EmailVerificationService
}

func NewServices(appCtx *appCtx.AppContext, repos repos.Repos) Services {
	accountService := NewAccountService(repos)
	emailVerificationService := NewEmailVerificationService(appCtx)

	return &services{
		accountService:           accountService,
		emailVerificationService: emailVerificationService,
	}
}

func (s *services) AccountService() AccountService {
	return s.accountService
}

func (s *services) EmailVerificationService() EmailVerificationService {
	return s.emailVerificationService
}
