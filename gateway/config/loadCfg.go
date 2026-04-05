package config

import (
	"context"
	stackErr "gateway/pkg/stackErr"

	"github.com/sethvargo/go-envconfig"
)

func LoadConfig(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:   cfg,
		Lookuper: envconfig.OsLookuper(),
	}); err != nil {
		return nil, stackErr.Error(err)
	}
	return cfg, nil
}
