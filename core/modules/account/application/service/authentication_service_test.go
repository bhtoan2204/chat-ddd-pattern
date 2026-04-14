package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-socket/core/modules/account/domain/aggregate"
	"go-socket/core/modules/account/domain/entity"
	"go-socket/core/modules/account/domain/repos"
	valueobject "go-socket/core/modules/account/domain/value_object"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/hasher"

	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestAuthenticationService_Authenticate_IssuesAccessAndRefreshTokens(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	accountAggregate := newRegisteredAccountAggregate(t, "acc-1", "alice@example.com", "hashed-password")
	accountAggregateRepo := repos.NewMockAccountAggregateRepository(ctrl)
	baseRepo := repos.NewMockRepos(ctrl)
	hasherMock := hasher.NewMockHasher(ctrl)
	pasetoMock := xpaseto.NewMockPasetoService(ctrl)

	accessExpiresAt := time.Date(2026, time.April, 14, 11, 0, 0, 0, time.UTC)
	refreshExpiresAt := time.Date(2026, time.April, 21, 11, 0, 0, 0, time.UTC)

	baseRepo.EXPECT().AccountAggregateRepository().Return(accountAggregateRepo)
	accountAggregateRepo.EXPECT().LoadByEmail(gomock.Any(), "alice@example.com").Return(accountAggregate, nil)
	hasherMock.EXPECT().Verify(gomock.Any(), "password123", "hashed-password").Return(true, nil)
	pasetoMock.EXPECT().
		GenerateAccessToken(gomock.Any(), gomock.AssignableToTypeOf(&entity.Account{})).
		DoAndReturn(func(_ context.Context, account *entity.Account) (string, time.Time, error) {
			if account == nil || account.ID != "acc-1" {
				t.Fatalf("expected account snapshot from aggregate, got %+v", account)
			}
			return "access-token", accessExpiresAt, nil
		})
	pasetoMock.EXPECT().
		GenerateRefreshToken(gomock.Any(), gomock.AssignableToTypeOf(&entity.Account{})).
		DoAndReturn(func(_ context.Context, account *entity.Account) (string, time.Time, error) {
			if account == nil || account.ID != "acc-1" {
				t.Fatalf("expected account snapshot from aggregate, got %+v", account)
			}
			return "refresh-token", refreshExpiresAt, nil
		})

	service := &authenticationService{
		baseRepo: baseRepo,
		hasher:   hasherMock,
		paseto:   pasetoMock,
	}

	result, err := service.Authenticate(context.Background(), AuthenticateAccountCommand{
		Email:    "alice@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
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

func TestAuthenticationService_RefreshAuthenticate_RotatesRefreshToken(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	accountAggregate := newRegisteredAccountAggregate(t, "acc-1", "alice@example.com", "hashed-password")
	accountAggregateRepo := repos.NewMockAccountAggregateRepository(ctrl)
	baseRepo := repos.NewMockRepos(ctrl)
	pasetoMock := xpaseto.NewMockPasetoService(ctrl)

	accessExpiresAt := time.Date(2026, time.April, 14, 12, 0, 0, 0, time.UTC)
	refreshExpiresAt := time.Date(2026, time.April, 21, 12, 0, 0, 0, time.UTC)

	pasetoMock.EXPECT().
		ParseRefreshToken(gomock.Any(), "incoming-refresh-token").
		Return(&xpaseto.PasetoPayload{
			AccountID: "acc-1",
			TokenType: xpaseto.TokenTypeRefresh,
		}, nil)
	baseRepo.EXPECT().AccountAggregateRepository().Return(accountAggregateRepo)
	accountAggregateRepo.EXPECT().Load(gomock.Any(), "acc-1").Return(accountAggregate, nil)
	pasetoMock.EXPECT().
		GenerateAccessToken(gomock.Any(), gomock.AssignableToTypeOf(&entity.Account{})).
		Return("rotated-access-token", accessExpiresAt, nil)
	pasetoMock.EXPECT().
		GenerateRefreshToken(gomock.Any(), gomock.AssignableToTypeOf(&entity.Account{})).
		Return("rotated-refresh-token", refreshExpiresAt, nil)

	service := &authenticationService{
		baseRepo: baseRepo,
		paseto:   pasetoMock,
	}

	result, err := service.RefreshAuthenticate(context.Background(), RefreshTokenCommand{
		RefreshToken: "incoming-refresh-token",
	})
	if err != nil {
		t.Fatalf("RefreshAuthenticate() error = %v", err)
	}

	if result.AccessToken != "rotated-access-token" {
		t.Fatalf("expected rotated-access-token, got %q", result.AccessToken)
	}
	if result.RefreshToken != "rotated-refresh-token" {
		t.Fatalf("expected rotated-refresh-token, got %q", result.RefreshToken)
	}
	if !result.AccessExpiresAt.Equal(accessExpiresAt) {
		t.Fatalf("expected access expiry %v, got %v", accessExpiresAt, result.AccessExpiresAt)
	}
	if !result.RefreshExpiresAt.Equal(refreshExpiresAt) {
		t.Fatalf("expected refresh expiry %v, got %v", refreshExpiresAt, result.RefreshExpiresAt)
	}
}

func TestAuthenticationService_Authenticate_MapsNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	accountAggregateRepo := repos.NewMockAccountAggregateRepository(ctrl)
	baseRepo := repos.NewMockRepos(ctrl)
	hasherMock := hasher.NewMockHasher(ctrl)
	pasetoMock := xpaseto.NewMockPasetoService(ctrl)

	baseRepo.EXPECT().AccountAggregateRepository().Return(accountAggregateRepo)
	accountAggregateRepo.EXPECT().LoadByEmail(gomock.Any(), "missing@example.com").Return(nil, gorm.ErrRecordNotFound)

	service := &authenticationService{
		baseRepo: baseRepo,
		hasher:   hasherMock,
		paseto:   pasetoMock,
	}

	_, err := service.Authenticate(context.Background(), AuthenticateAccountCommand{
		Email:    "missing@example.com",
		Password: "password123",
	})

	if !errors.Is(err, ErrAuthenticationAccountNotFound) {
		t.Fatalf("expected ErrAuthenticationAccountNotFound, got %v", err)
	}
}

func TestAuthenticationService_Authenticate_MapsInvalidPassword(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	accountAggregate := newRegisteredAccountAggregate(t, "acc-1", "alice@example.com", "hashed-password")
	accountAggregateRepo := repos.NewMockAccountAggregateRepository(ctrl)
	baseRepo := repos.NewMockRepos(ctrl)
	hasherMock := hasher.NewMockHasher(ctrl)
	pasetoMock := xpaseto.NewMockPasetoService(ctrl)

	baseRepo.EXPECT().AccountAggregateRepository().Return(accountAggregateRepo)
	accountAggregateRepo.EXPECT().LoadByEmail(gomock.Any(), "alice@example.com").Return(accountAggregate, nil)
	hasherMock.EXPECT().Verify(gomock.Any(), "password123", "hashed-password").Return(false, nil)

	service := &authenticationService{
		baseRepo: baseRepo,
		hasher:   hasherMock,
		paseto:   pasetoMock,
	}

	_, err := service.Authenticate(context.Background(), AuthenticateAccountCommand{
		Email:    "alice@example.com",
		Password: "password123",
	})

	if !errors.Is(err, ErrAuthenticationInvalidPassword) {
		t.Fatalf("expected ErrAuthenticationInvalidPassword, got %v", err)
	}
}

func newRegisteredAccountAggregate(t *testing.T, accountID, emailValue, passwordHashValue string) *aggregate.AccountAggregate {
	t.Helper()

	email, err := valueobject.NewEmail(emailValue)
	if err != nil {
		t.Fatalf("NewEmail() error = %v", err)
	}

	passwordHash, err := valueobject.NewHashedPassword(passwordHashValue)
	if err != nil {
		t.Fatalf("NewHashedPassword() error = %v", err)
	}

	accountAggregate, err := aggregate.NewAccountAggregate(accountID)
	if err != nil {
		t.Fatalf("NewAccountAggregate() error = %v", err)
	}

	now := time.Date(2026, time.April, 14, 8, 0, 0, 0, time.UTC)
	if err := accountAggregate.Register(email, passwordHash, "Alice", now); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	return accountAggregate
}
