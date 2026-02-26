package testhelpers

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer wraps the testcontainers postgres container with connection info.
type PostgresContainer struct {
	*postgres.PostgresContainer
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// RedisContainer wraps the testcontainers redis container with connection info.
type RedisContainer struct {
	*redis.RedisContainer
	Host string
	Port string
}

// SetupPostgres creates a PostgreSQL container with the project migrations applied.
func SetupPostgres(ctx context.Context, migrationsDir string) (*PostgresContainer, error) {
	dbName := "inventory_test"
	dbUser := "test"
	dbPassword := "test"

	// Collect migration files in order.
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to find migration files: %w", err)
	}

	initScripts := make([]testcontainers.ContainerFile, len(migrationFiles))
	for i, f := range migrationFiles {
		initScripts[i] = testcontainers.ContainerFile{
			HostFilePath:      f,
			ContainerFilePath: fmt.Sprintf("/docker-entrypoint-initdb.d/%03d_%s", i, filepath.Base(f)),
			FileMode:          0o644,
		}
	}

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithFiles(initScripts...),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: container,
		Host:              host,
		Port:              mappedPort.Port(),
		User:              dbUser,
		Password:          dbPassword,
		DBName:            dbName,
	}, nil
}

// SetupRedis creates a Redis container ready for testing.
func SetupRedis(ctx context.Context) (*RedisContainer, error) {
	container, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(15*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	return &RedisContainer{
		RedisContainer: container,
		Host:           host,
		Port:           mappedPort.Port(),
	}, nil
}
