package main

// @title Inventory API
// @version 1.0
// @description REST API for an Inventory Management System built with Go and Gin.

// @contact.name Jandir A. Cutabiala
// @contact.url https://github.com/jandiralceu/inventory_api_with_golang

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/jandiralceu/inventory_api_with_golang/internal/database"
	"github.com/jandiralceu/inventory_api_with_golang/internal/handlers"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
	"github.com/jandiralceu/inventory_api_with_golang/internal/routes"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize structured logger based on environment.
	pkg.InitLogger(cfg.Env)

	// Initialize OpenTelemetry tracing.
	shutdownTracer := pkg.InitTracer(ctx, cfg.AppName, cfg.Env, cfg.OTLPEndpoint)
	defer shutdownTracer(ctx)

	// Load RSA keys and initialize JWT manager.
	privateKeyPEM, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		slog.Error("Failed to read private key file", "path", cfg.PrivateKeyPath, "error", err)
		os.Exit(1)
	}
	publicKeyPEM, err := os.ReadFile(cfg.PublicKeyPath)
	if err != nil {
		slog.Error("Failed to read public key file", "path", cfg.PublicKeyPath, "error", err)
		os.Exit(1)
	}
	jwtManager, err := pkg.NewJWTManager(string(privateKeyPEM), string(publicKeyPEM))
	if err != nil {
		slog.Error("Failed to initialize JWT manager", "error", err)
		os.Exit(1)
	}

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

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	cacheManager := pkg.NewRedisCacheManager(cfg)
	defer cacheManager.Close()

	hasher := pkg.NewArgon2PasswordHasher()

	roleRepository := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepository, cacheManager)
	roleHandler := handlers.NewRoleHandler(roleService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, hasher)
	authHandler := handlers.NewAuthHandler(userService, jwtManager, cacheManager, hasher)
	userHandler := handlers.NewUserHandler(userService)

	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepository, cacheManager)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	supplierRepository := repository.NewSupplierRepository(db)
	supplierService := service.NewSupplierService(supplierRepository, cacheManager)
	supplierHandler := handlers.NewSupplierHandler(supplierService)

	warehouseRepository := repository.NewWarehouseRepository(db)
	warehouseService := service.NewWarehouseService(warehouseRepository, cacheManager)
	warehouseHandler := handlers.NewWarehouseHandler(warehouseService)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository, cacheManager)
	productHandler := handlers.NewProductHandler(productService)

	inventoryRepository := repository.NewInventoryRepository(db)
	inventoryTransactionRepository := repository.NewInventoryTransactionRepository(db)
	inventoryService := service.NewInventoryService(inventoryRepository, productRepository, warehouseRepository, inventoryTransactionRepository, cacheManager)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	// Initialize Casbin Enforcer for RBAC.
	enforcer, err := casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		slog.Error("Failed to initialize Casbin enforcer", "error", err)
		os.Exit(1)
	}

	routeConfig := &routes.RouteConfig{
		AuthHandler:      authHandler,
		RoleHandler:      roleHandler,
		UserHandler:      userHandler,
		CategoryHandler:  categoryHandler,
		SupplierHandler:  supplierHandler,
		WarehouseHandler: warehouseHandler,
		ProductHandler:   productHandler,
		InventoryHandler: inventoryHandler,
	}

	r := routes.Setup(routeConfig, cfg, jwtManager, enforcer, cacheManager)

	// Health check endpoint for monitoring and load balancer probes.
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Create custom HTTP server for graceful shutdown control.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a separate goroutine so it doesn't block.
	go func() {
		slog.Info("Server starting", "port", cfg.AppPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Create a channel to listen for OS interrupt signals (SIGINT, SIGTERM).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-quit
	slog.Info("Received signal. Shutting down gracefully...", "signal", sig.String())

	// Give in-flight requests up to 10 seconds to complete.
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown of the HTTP server.
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited gracefully")
}
