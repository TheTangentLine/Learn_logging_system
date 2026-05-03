package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ESURL       string
	RabbitMQURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		return nil, errors.New("RABBITMQ_URL environment variable is required")
	}

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		return nil, errors.New("ELASTICSEARCH_URL environment variable is required")
	}

	return &Config{RabbitMQURL: rabbitMQURL, ESURL: esURL}, nil
}
