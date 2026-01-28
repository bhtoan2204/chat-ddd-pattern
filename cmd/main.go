package main

import (
	"context"
	"flag"
	"fmt"
	"go-socket/config"
	appCtx "go-socket/core/context"
	"go-socket/core/delivery/http"
	"go-socket/core/shared/infra/db"
	"go-socket/core/shared/pkg/logging"
	"os/signal"
	"syscall"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	logger := logging.FromContext(ctx)
	logger.Infow("Starting application")
	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Errorw("Recovered from panic", "error", r)
		}
	}()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		logger.Errorw("Failed to load config", "error", err)
		return
	}
	appCtx, err := appCtx.LoadAppCtx(ctx, cfg)
	if err != nil {
		logger.Errorw("Failed to create app context", "error", err)
		return
	}

	migrateTool := db.NewMigrateTool()
	pathMigration := flag.String("path", "migration/", "path to migrations folder")
	if err := migrateTool.Migrate(fmt.Sprintf("file://%s", *pathMigration), cfg.DBConfig.ConnectionURL); err != nil {
		logger.Errorw("Failed to migrate database", "error", err)
		return
	}

	appHttp := http.NewServer(cfg)
	if err := appHttp.Start(ctx, appCtx); err != nil {
		logger.Errorw("Failed to start app http", "error", err)
		return
	}

	<-ctx.Done()
	appCtx.Close()
}
