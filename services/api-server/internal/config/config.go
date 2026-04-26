package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBWriteURL string // DATABASE_WRITE_URL — primary Postgres, used for writes
	DBReadURL  string // DATABASE_READ_URL  — read replica, used by relay worker and read queries
	ESURL      string
}

func Load() (*Config, error) {
	// Load .env if present; ignore error when file doesn't exist
	_ = godotenv.Load()

	writeURL := os.Getenv("DATABASE_WRITE_URL")
	if writeURL == "" {
		return nil, errors.New("DATABASE_WRITE_URL environment variable is required")
	}

	// Read URL is optional; falls back to the write URL for single-node setups
	readURL := getEnv("DATABASE_READ_URL", writeURL)

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBWriteURL: writeURL,
		DBReadURL:  readURL,
		ESURL:      getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
