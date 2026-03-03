package pkg

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisCacheManager_SetGetDelete(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	cfg := &config.Config{
		RedisHost: mr.Host(),
		RedisPort: mr.Port(),
	}

	cache := NewRedisCacheManager(cfg)
	defer func() { _ = cache.Close() }()

	ctx := context.Background()

	// Test Set
	err = cache.Set(ctx, "test-key", map[string]string{"foo": "bar"}, 10*time.Minute)
	require.NoError(t, err)

	// Test Get
	var dest map[string]string
	err = cache.Get(ctx, "test-key", &dest)
	require.NoError(t, err)
	assert.Equal(t, "bar", dest["foo"])

	// Test Delete
	err = cache.Delete(ctx, "test-key")
	require.NoError(t, err)

	// Ensure deleted
	var dest2 map[string]string
	err = cache.Get(ctx, "test-key", &dest2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "redis: nil")
}

func TestRedisCacheManager_DeletePrefix(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	cfg := &config.Config{
		RedisHost: mr.Host(),
		RedisPort: mr.Port(),
	}

	cache := NewRedisCacheManager(cfg)
	defer func() { _ = cache.Close() }()

	ctx := context.Background()

	_ = cache.Set(ctx, "prefix-1", "val1", 0)
	_ = cache.Set(ctx, "prefix-2", "val2", 0)
	_ = cache.Set(ctx, "other-3", "val3", 0)

	err = cache.DeletePrefix(ctx, "prefix-")
	require.NoError(t, err)

	var dest string
	err = cache.Get(ctx, "prefix-1", &dest)
	require.Error(t, err)

	err = cache.Get(ctx, "other-3", &dest)
	require.NoError(t, err)
	assert.Equal(t, "val3", dest)
}

func TestRedisCacheManager_SetInvalidJSON(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	cfg := &config.Config{
		RedisHost: mr.Host(),
		RedisPort: mr.Port(),
	}

	cache := NewRedisCacheManager(cfg)
	defer func() { _ = cache.Close() }()

	// Channels cannot be JSON-marshaled
	invalidValue := make(chan int)

	err = cache.Set(context.Background(), "invalid-key", invalidValue, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type: chan int")
}
