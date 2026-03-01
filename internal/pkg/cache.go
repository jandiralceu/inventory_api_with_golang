package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// CacheManager abstracts a key-value store used for session management and performance.
type CacheManager interface {
	// Set stores a value with an expiration duration.
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	// Get retrieves a value and unmarshals it into the destination pointer.
	Get(ctx context.Context, key string, dest any) error
	// Delete removes a specific key.
	Delete(ctx context.Context, key string) error
	// DeletePrefix removes all keys matching the specified prefix using SCAN.
	DeletePrefix(ctx context.Context, prefix string) error
	// Close releases the underlying connection resources.
	Close() error
	// GetClient returns the underlying redis client.
	GetClient() *redis.Client
}

type redisCacheManager struct {
	client *redis.Client
}

// NewRedisCacheManager creates a CacheManager backed by Redis and enables
// OpenTelemetry instrumentation for command tracing.
func NewRedisCacheManager(cfg *config.Config) CacheManager {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
	} else {
		slog.Info("Connected to Redis", "addr", addr)
	}

	// Enable Redis OpenTelemetry tracing
	if err := redisotel.InstrumentTracing(client); err != nil {
		slog.Error("Failed to register Redis tracing", "error", err)
	}

	return &redisCacheManager{client: client}
}

func (m *redisCacheManager) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(ctx, key, bytes, expiration).Err()
}

func (m *redisCacheManager) Get(ctx context.Context, key string, dest any) error {
	val, err := m.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (m *redisCacheManager) Delete(ctx context.Context, key string) error {
	return m.client.Del(ctx, key).Err()
}

func (m *redisCacheManager) DeletePrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var err error
	var keys []string

	for {
		keys, cursor, err = m.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := m.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}

func (m *redisCacheManager) Close() error {
	return m.client.Close()
}

func (m *redisCacheManager) GetClient() *redis.Client {
	return m.client
}
