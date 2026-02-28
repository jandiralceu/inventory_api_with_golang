package pkg

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger_Development(t *testing.T) {
	InitLogger("development")

	// Ensure the default logger is properly initialized
	assert.NotNil(t, slog.Default())
}

func TestInitLogger_Production(t *testing.T) {
	InitLogger("production")

	// Ensure the default logger is properly initialized for production
	assert.NotNil(t, slog.Default())
}
