package xpaseto

import (
	"context"
	"go-socket/core/modules/account/domain/entity"
	"go-socket/core/shared/config"
	"testing"
)

func TestPaseto(t *testing.T) {
	cfg := &config.Config{
		AuthConfig: config.AuthConfig{
			TokenIssuer:            "chat",
			AccessTokenTTLSeconds:  9000,
			RefreshTokenTTLSeconds: 26400,
			AccessPublicKey:        "g2NuXbGMgDnw04S8KmeKqJJ94WwABPoe/2HB66V1+QM=",
			AccessPrivateKey:       "OghFb8xO1EqyzKRc1/q7hgAkNzZfZJXOkczIoey2+ViDY25dsYyAOfDThLwqZ4qokn3hbAAE+h7/YcHrpXX5Aw==",
			RefreshPublicKey:       "g2NuXbGMgDnw04S8KmeKqJJ94WwABPoe/2HB66V1+QM=",
			RefreshPrivateKey:      "OghFb8xO1EqyzKRc1/q7hgAkNzZfZJXOkczIoey2+ViDY25dsYyAOfDThLwqZ4qokn3hbAAE+h7/YcHrpXX5Aw==",
		},
	}
	pasetoSvc, err := NewPaseto(cfg)
	if err != nil {
		t.Fatal(err)
	}
	str, _, _ := pasetoSvc.GenerateAccessToken(context.Background(), &entity.Account{
		ID: "test-abc-001",
	})
	claims, err := pasetoSvc.ParseAccessToken(context.Background(), str)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}

func TestParsePaseto(t *testing.T) {
	cfg := &config.Config{
		AuthConfig: config.AuthConfig{
			TokenIssuer:            "chat",
			AccessTokenTTLSeconds:  9000,
			RefreshTokenTTLSeconds: 26400,
			AccessPublicKey:        "vSKvNvjpCS3teuTBeXm9gHYSIGLaovZoM+vMnyNeFKk=",
			AccessPrivateKey:       "CncqpMFMEHuK1As2dIRECZ2qLZJAqgJKZmP9KdN+vLO9Iq82+OkJLe165MF5eb2AdhIgYtqi9mgz68yfI14UqQ==",
			RefreshPublicKey:       "g2NuXbGMgDnw04S8KmeKqJJ94WwABPoe/2HB66V1+QM=",
			RefreshPrivateKey:      "OghFb8xO1EqyzKRc1/q7hgAkNzZfZJXOkczIoey2+ViDY25dsYyAOfDThLwqZ4qokn3hbAAE+h7/YcHrpXX5Aw==",
		},
	}
	pasetoSvc, err := NewPaseto(cfg)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := pasetoSvc.ParseAccessToken(context.Background(), "v2.public.eyJlbWFpbCI6ImhvYW5ndHVkYWRlbkBnbWFpbC5jb20iLCJleHAiOiIyMDI2LTA0LTE1VDA2OjE0OjAxWiIsImlhdCI6IjIwMjYtMDQtMTVUMDM6NDQ6MDFaIiwiaXNzIjoiY2hhdCIsInN1YiI6IjQwOWM1MmVhLTcxYTMtNDUxMy05MzQwLWU5NDdhYTc0MTI3MyIsInRva2VuX3VzZSI6ImFjY2VzcyJ9JK56WWwIOAIKpO_8ES4dR7tu5pZDNkiJB2WUhKczKNUqUF5nwTyybeYAN3-y8fPQJdNnMf6zxGNjV16Mw6FmDA.bnVsbA")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
