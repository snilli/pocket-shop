package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-simpler.org/env"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" default:"info"`

	PostgresHost     string `env:"POSTGRES_HOST" default:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT" default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" default:"ez" validate:"required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" default:"ez" validate:"required"`
	PostgresDB       string `env:"POSTGRES_DB" default:"ez" validate:"required"`
	PostgresSSLMode  string `env:"POSTGRES_SSLMODE" default:"disable"`

	ServerPort    string `env:"SERVER_PORT" default:"8080"`
	ServerHost    string `env:"SERVER_HOST" default:"0.0.0.0"`
	ServerMode    string `env:"SERVER_MODE" default:"debug"`
	EnableSwagger bool   `env:"ENABLE_SWAGGER" default:"true"`

	EZBaseURL   string `env:"EZ_BASE_URL" default:"https://api.ezcards.io"`
	EZAPIKey    string `env:"EZ_API_KEY" validate:"required"`
	EZAuthToken string `env:"EZ_ACCESS_TOKEN" validate:"required"`

	EZSKU string `env:"EZ_SKU" validate:"required"`

	EZRetryMaxAttempts int `env:"EZ_RETRY_MAX_ATTEMPTS" default:"3"`
	EZRetryBackoffSec  int `env:"EZ_RETRY_BACKOFF_SECONDS" default:"1"`

	RefSource string `env:"REF_SOURCE" default:"ez"`

	OrderFulfillmentTimeoutSec int `env:"ORDER_FULFILLMENT_TIMEOUT_SECONDS" default:"60"`
	PollIntervalSec            int `env:"POLL_INTERVAL_SECONDS" default:"2"`

	DiscoverIntervalSec int `env:"DISCOVER_INTERVAL_SECONDS" default:"5"`
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PostgresSSLMode)
}

var (
	validate = validator.New()
	instance *Config
	once     sync.Once
)

func Load() (*Config, error) {
	var loadErr error

	once.Do(func() {
		_ = godotenv.Load()

		cfg := &Config{}

		if err := env.Load(cfg, nil); err != nil {
			loadErr = fmt.Errorf("failed to load environment variables: %w", err)
			return
		}
		if err := validate.Struct(cfg); err != nil {
			loadErr = fmt.Errorf("config validation failed: %w", err)
			return
		}
		instance = cfg
	})

	if loadErr != nil {
		return nil, loadErr
	}

	return instance, nil
}

func Get() *Config {
	if instance == nil {
		panic("config not loaded, call Load() first")
	}
	return instance
}

func (c *Config) SetupLogger() {
	switch c.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if c.ServerMode == "production" {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		})
	}
}
