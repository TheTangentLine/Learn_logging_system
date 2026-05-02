package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBWriteURL   string        // DATABASE_WRITE_URL (required) — primary Postgres
	RabbitMQURL  string        // RABBITMQ_URL       (required)
	BatchSize    int           // BATCH_SIZE          (default 10)
	PollInterval time.Duration // POLL_INTERVAL       (default 30s) — fallback if no LISTEN notification
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbWriteURL := os.Getenv("DATABASE_WRITE_URL")
	if dbWriteURL == "" {
		return nil, errors.New("DATABASE_WRITE_URL environment variable is required")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		return nil, errors.New("RABBITMQ_URL environment variable is required")
	}

	return &Config{
		DBWriteURL:   dbWriteURL,
		RabbitMQURL:  rabbitMQURL,
		BatchSize:    getEnvInt("BATCH_SIZE", 10),
		PollInterval: getEnvDuration("POLL_INTERVAL", 30*time.Second),
	}, nil
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
