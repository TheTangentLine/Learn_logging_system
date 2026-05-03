package config

import (
	"testing"
)

func TestLoad_requiredEnv(t *testing.T) {
	t.Chdir(t.TempDir())

	t.Run("missing RABBITMQ_URL", func(t *testing.T) {
		t.Setenv("RABBITMQ_URL", "")
		t.Setenv("ELASTICSEARCH_URL", "http://localhost:9200")
		_, err := Load()
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("missing ELASTICSEARCH_URL", func(t *testing.T) {
		t.Setenv("RABBITMQ_URL", "amqp://guest:guest@rabbit:5672/")
		t.Setenv("ELASTICSEARCH_URL", "")
		_, err := Load()
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestLoad_success(t *testing.T) {
	t.Chdir(t.TempDir())

	t.Setenv("RABBITMQ_URL", "amqp://guest:guest@rabbit:5672/")
	t.Setenv("ELASTICSEARCH_URL", "http://elasticsearch:9200")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.RabbitMQURL != "amqp://guest:guest@rabbit:5672/" || cfg.ESURL != "http://elasticsearch:9200" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
