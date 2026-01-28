package xpaseto

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"go-socket/core/domain/account/entity"
	"go-socket/core/shared/pkg/logging"

	"github.com/o1egl/paseto"
	"go.uber.org/zap"
)

type PasetoPayload struct {
	AccountID string
	Email     string
	ExpiresAt time.Time
}

type PasetoService interface {
	GenerateToken(ctx context.Context, account *entity.Account) (string, time.Time, error)
	ParseToken(ctx context.Context, token string) (*PasetoPayload, error)
}

type pasetoService struct {
	paseto       *paseto.V2
	symmetricKey []byte
	issuer       string
	ttl          time.Duration
}

func NewPaseto(symmetricKey string, issuer string, ttl time.Duration) (PasetoService, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(symmetricKey)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("paseto key must be 32 bytes")
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("token ttl must be positive")
	}
	return &pasetoService{
		paseto:       paseto.NewV2(),
		symmetricKey: keyBytes,
		issuer:       issuer,
		ttl:          ttl * time.Second,
	}, nil
}

func (p *pasetoService) GenerateToken(ctx context.Context, account *entity.Account) (string, time.Time, error) {
	if account == nil {
		return "", time.Time{}, fmt.Errorf("account is nil")
	}
	now := time.Now().UTC()
	exp := now.Add(p.ttl).UTC()
	payload := paseto.JSONToken{
		Issuer:     p.issuer,
		Subject:    account.ID,
		IssuedAt:   now.UTC(),
		Expiration: exp.UTC(),
	}
	payload.Set("email", account.Email)

	token, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, exp, nil
}

func (p *pasetoService) ParseToken(ctx context.Context, token string) (*PasetoPayload, error) {
	logger := logging.FromContext(ctx)
	var jsonToken paseto.JSONToken
	var custom map[string]interface{}
	if err := p.paseto.Decrypt(token, p.symmetricKey, &jsonToken, &custom); err != nil {
		logger.Errorw("Parse token failed", zap.Error(err))
		return nil, err
	}
	if !jsonToken.Expiration.IsZero() && time.Now().UTC().After(jsonToken.Expiration.UTC()) {
		return nil, errors.New("token expired")
	}
	email, _ := custom["email"].(string)
	return &PasetoPayload{
		AccountID: jsonToken.Subject,
		Email:     email,
		ExpiresAt: jsonToken.Expiration,
	}, nil
}
