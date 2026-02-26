package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init opens a database connection using the settings provided in [config.Config].
// It verifies the connection is alive using the context and registers an
// OpenTelemetry plugin for query tracing.
func Init(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify the connection is actually alive within the given context.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Connected to Postgres", "host", cfg.DBHost, "dbname", cfg.DBName)

	// Register OpenTelemetry plugin to trace all database queries.
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		slog.Error("Failed to register otel gorm plugin", "error", err)
		return nil, fmt.Errorf("failed to register otel gorm plugin: %w", err)
	}

	return db, nil
}
