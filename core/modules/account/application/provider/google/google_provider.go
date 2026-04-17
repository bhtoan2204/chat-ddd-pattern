package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"wechat-clone/core/modules/account/application/provider"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleProvider struct {
	oauthClient *oauth2.Config
}

func NewGoogleProvider(ctx context.Context, cfg *config.Config) provider.AuthProvider {
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  cfg.AuthConfig.GoogleConfig.GoogleRedirectURL,
		ClientID:     cfg.AuthConfig.GoogleConfig.GoogleClientID,
		ClientSecret: cfg.AuthConfig.GoogleConfig.GoogleClientSecret,
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		Endpoint: google.Endpoint,
	}

	return &googleProvider{
		oauthClient: googleOauthConfig,
	}
}

func (g *googleProvider) Login() string {
	return g.oauthClient.AuthCodeURL(uuid.NewString(), oauth2.AccessTypeOffline)
}

func (g *googleProvider) Callback(ctx context.Context, code string) (*provider.AuthResult, error) {
	token, err := g.oauthClient.Exchange(ctx, code)
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("exchange oauth code: %w", err))
	}

	result := &provider.AuthResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiryUnix:   token.Expiry.Unix(),
	}

	if idToken, ok := token.Extra("id_token").(string); ok {
		result.IDToken = idToken
	}

	return result, nil
}

func (g *googleProvider) UserInfo(ctx context.Context, accessToken string) (*provider.UserInfo, error) {
	client := g.oauthClient.Client(ctx, &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	})

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo",
		nil,
	)
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("create userinfo request: %w", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("request userinfo: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, stackErr.Error(fmt.Errorf("userinfo request failed with status: %d", resp.StatusCode))
	}

	var userInfo provider.UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, stackErr.Error(fmt.Errorf("decode userinfo response: %w", err))
	}

	return &userInfo, nil
}

func (g *googleProvider) Name() string {
	return "google"
}
