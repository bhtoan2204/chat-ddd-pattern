package google

import (
	"context"
	"go-socket/core/shared/config"
	"testing"
)

func TestUrlLogin(t *testing.T) {
	googleProvider := NewGoogleProvider(context.Background(), &config.Config{
		AuthConfig: config.AuthConfig{
			GoogleConfig: config.GoogleConfig{
				GoogleClientID:     "",
				GoogleClientSecret: "",
				GoogleRedirectURL:  "http://localhost:35000/api/v1/auth/google/callback",
			},
		},
	})
	t.Log(googleProvider.Login())
}
