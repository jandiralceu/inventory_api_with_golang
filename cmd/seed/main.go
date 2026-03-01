package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/jandiralceu/inventory_api_with_golang/internal/database"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize structured logger
	pkg.InitLogger(cfg.Env)

	slog.Info("Starting database seed process...")

	// Initialize database connection
	db, err := database.Init(ctx, cfg)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get underlying DB connection", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	// Run SeedRoles
	if err := database.SeedRoles(ctx, db); err != nil {
		slog.Error("Failed to seed roles", "error", err)
		os.Exit(1)
	}

	slog.Info("Database seeding completed successfully.")
}
