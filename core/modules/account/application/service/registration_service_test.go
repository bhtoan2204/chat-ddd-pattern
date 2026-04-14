package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-socket/core/modules/account/domain/aggregate"
	"go-socket/core/modules/account/domain/entity"
	"go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/hasher"

	"go.uber.org/mock/gomock"
)

func TestRegistrationService_Register_IssuesAccessAndRefreshTokens(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	baseRepo := repos.NewMockRepos(ctrl)
	txRepos := repos.NewMockRepos(ctrl)
	accountRepo := repos.NewMockAccountRepository(ctrl)
	accountAggregateRepo := repos.NewMockAccountAggregateRepository(ctrl)
	hasherMock := hasher.NewMockHasher(ctrl)
	pasetoMock := xpaseto.NewMockPasetoService(ctrl)

	accessExpiresAt := time.Date(2026, time.April, 14, 10, 0, 0, 0, time.UTC)
	refreshExpiresAt := time.Date(2026, time.April, 21, 10, 0, 0, 0, time.UTC)

	baseRepo.EXPECT().AccountRepository().Return(accountRepo)
	accountRepo.EXPECT().IsEmailExists(gomock.Any(), "alice@example.com").Return(false, nil)
	hasherMock.EXPECT().Hash(gomock.Any(), "password123").Return("hashed-password", nil)
	baseRepo.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(repos.Repos) error) error {
			return fn(txRepos)
		})
	txRepos.EXPECT().AccountAggregateRepository().Return(accountAggregateRepo)
	accountAggregateRepo.EXPECT().
		Save(gomock.Any(), gomock.AssignableToTypeOf(&aggregate.AccountAggregate{})).
		Return(nil)
	pasetoMock.EXPECT().
		GenerateAccessToken(gomock.Any(), gomock.AssignableToTypeOf(&entity.Account{})).
		DoAndReturn(func(_ context.Context, account *entity.Account) (string, time.Time, error) {
			if account == nil || account.ID == "" {
				t.Fatalf("expected non-nil account snapshot, got %+v", account)
			}
			return "access-token", accessExpiresAt, nil
		})
	pasetoMock.EXPECT().
		GenerateRefreshToken(gomock.Any(), gomock.AssignableToTypeOf(&entity.Account{})).
		DoAndReturn(func(_ context.Context, account *entity.Account) (string, time.Time, error) {
			if account == nil || account.ID == "" {
				t.Fatalf("expected non-nil account snapshot, got %+v", account)
			}
			return "refresh-token", refreshExpiresAt, nil
		})

	service := &registrationService{
		baseRepo: baseRepo,
		hasher:   hasherMock,
		paseto:   pasetoMock,
	}

	result, err := service.Register(context.Background(), RegisterAccountCommand{
		Email:       "alice@example.com",
		Password:    "password123",
		DisplayName: "Alice",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if result.AccessToken != "access-token" {
		t.Fatalf("expected access-token, got %q", result.AccessToken)
	}
	if result.RefreshToken != "refresh-token" {
		t.Fatalf("expected refresh-token, got %q", result.RefreshToken)
	}
	if !result.AccessExpiresAt.Equal(accessExpiresAt) {
		t.Fatalf("expected access expiry %v, got %v", accessExpiresAt, result.AccessExpiresAt)
	}
	if !result.RefreshExpiresAt.Equal(refreshExpiresAt) {
		t.Fatalf("expected refresh expiry %v, got %v", refreshExpiresAt, result.RefreshExpiresAt)
	}
}

func TestRegistrationService_Register_ReturnsAccountExists(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	baseRepo := repos.NewMockRepos(ctrl)
	accountRepo := repos.NewMockAccountRepository(ctrl)
	hasherMock := hasher.NewMockHasher(ctrl)
	pasetoMock := xpaseto.NewMockPasetoService(ctrl)

	baseRepo.EXPECT().AccountRepository().Return(accountRepo)
	accountRepo.EXPECT().IsEmailExists(gomock.Any(), "alice@example.com").Return(true, nil)

	service := &registrationService{
		baseRepo: baseRepo,
		hasher:   hasherMock,
		paseto:   pasetoMock,
	}

	_, err := service.Register(context.Background(), RegisterAccountCommand{
		Email:       "alice@example.com",
		Password:    "password123",
		DisplayName: "Alice",
	})
	if !errors.Is(err, ErrRegistrationAccountExists) {
		t.Fatalf("expected ErrRegistrationAccountExists, got %v", err)
	}
}
