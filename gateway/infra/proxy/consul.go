package proxy

import (
	"context"
	"crypto/tls"
	"gateway/config"
	"gateway/pkg/logging"
	stackErr "gateway/pkg/stackErr"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func NewConsulClient(ctx context.Context, cfg *config.Config) (*api.Client, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.ConsulConfig.Address
	consulConfig.Scheme = cfg.ConsulConfig.Scheme
	consulConfig.Token = cfg.ConsulConfig.Token
	consulConfig.Datacenter = cfg.ConsulConfig.DataCenter

	if os.Getenv("ENVIRONMENT") == "production" {
		consulConfig.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
	}

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logging.FromContext(ctx).Error("Failed to connect to Consul:", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return consulClient, nil
}
