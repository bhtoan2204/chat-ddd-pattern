package persistent

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-socket/config"

	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
)

func NewConnection(ctx context.Context, cfg *config.Config) (*gorm.DB, *sql.DB, error) {
	dialector := oracle.New(oracle.Config{
		DSN: cfg.DBConfig.ConnectionURL,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open gorm oracle failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get sql db failed: %w", err)
	}

	// Pool config
	sqlDB.SetMaxOpenConns(cfg.DBConfig.MaxOpenConnNumber)
	sqlDB.SetMaxIdleConns(cfg.DBConfig.MaxIdleConnNumber)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConfig.ConnMaxLifeTimeSeconds) * time.Second)

	// Health check
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("ping db failed: %w", err)
	}

	return db, sqlDB, nil
}
