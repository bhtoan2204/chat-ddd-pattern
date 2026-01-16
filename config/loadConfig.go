package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

func LoadConfig(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:   cfg,
		Lookuper: envconfig.OsLookuper(),
	}); err != nil {
		return nil, fmt.Errorf("envconfig.ProcessWith has err=%w", err)
	}
	return cfg, nil
}
