package service

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/domain/repos"
)

type Services interface {
	AccountService() AccountService
	AuthenticationService() AuthenticationService
	EmailVerificationService() EmailVerificationService
	RegistrationService() RegistrationService
}

type services struct {
	accountService           AccountService
	authenticationService    AuthenticationService
	emailVerificationService EmailVerificationService
	registrationService      RegistrationService
}

func NewServices(appCtx *appCtx.AppContext, repos repos.Repos) Services {
	accountService := NewAccountService(repos)
	authenticationService := NewAuthenticationService(appCtx, repos)
	emailVerificationService := NewEmailVerificationService(appCtx)
	registrationService := NewRegistrationService(appCtx, repos)

	return &services{
		accountService:           accountService,
		authenticationService:    authenticationService,
		emailVerificationService: emailVerificationService,
		registrationService:      registrationService,
	}
}

func (s *services) AccountService() AccountService {
	return s.accountService
}

func (s *services) AuthenticationService() AuthenticationService {
	return s.authenticationService
}

func (s *services) EmailVerificationService() EmailVerificationService {
	return s.emailVerificationService
}

func (s *services) RegistrationService() RegistrationService {
	return s.registrationService
}
