package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App       App
		CORS      CORS
		Cache     Cache
		HTTP      HTTP
		JWT       JWT
		Log       Log
		Metrics   Metrics
		PG        PG
		RateLimit RateLimit
		Swagger   Swagger
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// CORS -.
	CORS struct {
		AllowOrigins     string `env:"CORS_ALLOW_ORIGINS,required"`
		AllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS,required"`
		AllowHeaders     string `env:"CORS_ALLOW_HEADERS,required"`
		AllowMethods     string `env:"CORS_ALLOW_METHODS,required"`
	}

	// Cache -.
	Cache struct {
		MaxCost int           `env:"CACHE_MAX_COST" envDefault:"10000"`
		TTL     time.Duration `env:"CACHE_TTL"      envDefault:"5m"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	// JWT -.
	JWT struct {
		Issuer            string        `env:"JWT_ISSUER,required"`
		Secret            string        `env:"JWT_SECRET,required"` //nolint:gosec // intended
		Expiration        time.Duration `env:"JWT_EXPIRATION,required"`
		RefreshIssuer     string        `env:"JWT_REFRESH_ISSUER,required"`
		RefreshSecret     string        `env:"JWT_REFRESH_SECRET,required"` //nolint:gosec // intended
		RefreshExpiration time.Duration `env:"JWT_REFRESH_EXPIRATION,required"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// RateLimit -.
	RateLimit struct {
		RequestsPerSecond float64 `env:"RATE_LIMIT_RPS"   envDefault:"10"`
		Burst             int     `env:"RATE_LIMIT_BURST" envDefault:"20"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
