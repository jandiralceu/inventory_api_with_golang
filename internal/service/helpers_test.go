package service

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockCacheManager is a shared mock for pkg.CacheManager used across service tests.
type MockCacheManager struct {
	mock.Mock
}

func (m *MockCacheManager) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheManager) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCacheManager) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheManager) DeletePrefix(ctx context.Context, prefix string) error {
	args := m.Called(ctx, prefix)
	return args.Error(0)
}

func (m *MockCacheManager) Close() error {
	args := m.Called()
	return args.Error(0)
}
