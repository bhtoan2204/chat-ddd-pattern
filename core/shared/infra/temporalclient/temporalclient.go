package temporalclient

import (
	"context"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/logging"

	"go.temporal.io/sdk/client"
)

func NewTemporalClient(ctx context.Context, cfg *config.Config) client.Client {
	log := logging.FromContext(ctx)
	c, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalConfig.Address,
		Namespace: cfg.TemporalConfig.Namespace,
	})
	if err != nil {
		log.Fatalf("failed to create Temporal client: %v", err)
	}

	return c
}
