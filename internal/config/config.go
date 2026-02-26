// Package config loads application settings from environment variables,
// optionally falling back to a .env file for local development.
//
// Environment variables are mapped to [Config] fields via struct tags
// from [github.com/caarlos0/env]. Fields without a value and without an
// envDefault tag will cause [Load] to return an error.
package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds every environment-driven setting the application needs.
// Struct tags define the variable name and optional default value.
type Config struct {
	DBHost         string `env:"DB_HOST"`
	DBUser         string `env:"DB_USER"`
	DBPassword     string `env:"DB_PASSWORD"`
	DBName         string `env:"DB_NAME"`
	DBPort         string `env:"DB_PORT"`
	DBSSLMode      string `env:"DB_SSL_MODE" envDefault:"disable"`
	AppPort        string `env:"APP_PORT" envDefault:"8080"`
	Env            string `env:"ENV" envDefault:"development"`
	AppName        string `env:"APP_NAME" envDefault:"cms-api"`
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH" envDefault:"private.pem"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH" envDefault:"public.pem"`
	OTLPEndpoint   string `env:"OTLP_ENDPOINT" envDefault:"localhost:4317"`
	RedisHost      string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort      string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword  string `env:"REDIS_PASSWORD" envDefault:""`
}

// Load reads a .env file (if present) and parses environment variables into a [Config].
// In production the .env file is typically absent; values come from the runtime environment.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using system environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		slog.Error("Error parsing configuration", "error", err)
		return nil, err
	}
	return cfg, nil
}
