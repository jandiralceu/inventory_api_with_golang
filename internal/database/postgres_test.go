package database_test

import (
	"context"
	"testing"

	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/jandiralceu/inventory_api_with_golang/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestInit_DatabaseConnectionFailure(t *testing.T) {
	// Provide a config that will definitely fail to connect (invalid port/host).
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "0", // Invalid port shouldn't be accessible
		DBUser:     "fakeuser",
		DBPassword: "fakepassword",
		DBName:     "fakedb",
	}

	db, err := database.Init(context.Background(), cfg)

	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to connect to database", "Error message should mention connection failure")
}
