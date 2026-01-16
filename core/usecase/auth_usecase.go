package usecase

import "context"

type AuthUsecase interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) (string, error)
	Logout(ctx context.Context, token string) error
}
