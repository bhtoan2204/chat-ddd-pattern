package main

import (
	"context"
	"errors"
	"gateway/config"
	"gateway/pkg/logging"
	"gateway/transport"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := logging.FromContext(ctx)
	logger.Infow("Starting application")

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		logger.Errorw("Failed to load config", zap.Error(err))
		return
	}

	httpTransport, err := transport.NewHTTPTransport(ctx, cfg)
	if err != nil {
		logger.Errorw("Failed to create HTTP transport", zap.Error(err))
		return
	}

	startErrCh := make(chan error, 1)
	go func() {
		startErrCh <- httpTransport.Start()
	}()

	select {
	case <-ctx.Done():
	case err := <-startErrCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorw("HTTP server stopped", zap.Error(err))
		}
		return
	}

	logger.Infow("Shutting down server")

	if err := httpTransport.Stop(); err != nil {
		logger.Errorw("Failed to stop HTTP transport", zap.Error(err))
	}

	logger.Infow("Application stopped")
}
