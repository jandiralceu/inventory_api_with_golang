package platform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTracer_Success(t *testing.T) {
	// grpc.NewClient is lazy and won't throw error immediately even if address doesn't exist
	shutdown := InitTracer(context.Background(), "test-service", "test", "127.0.0.1:4317")

	assert.NotNil(t, shutdown)

	// Invoke the returned shutdown hook to ensure no panics occur
	shutdown(context.Background())
}
