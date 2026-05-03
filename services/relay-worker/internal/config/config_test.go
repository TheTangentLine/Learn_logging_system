package config

import (
	"testing"
	"time"
)

func TestLoad_requiredEnv(t *testing.T) {
	t.Chdir(t.TempDir())

	t.Run("missing DATABASE_WRITE_URL", func(t *testing.T) {
		t.Setenv("DATABASE_WRITE_URL", "")
		t.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
		_, err := Load()
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("missing RABBITMQ_URL", func(t *testing.T) {
		t.Setenv("DATABASE_WRITE_URL", "postgres://u:p@localhost:5432/db")
		t.Setenv("RABBITMQ_URL", "")
		_, err := Load()
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestLoad_successAndOverrides(t *testing.T) {
	t.Chdir(t.TempDir())

	t.Setenv("DATABASE_WRITE_URL", "postgres://u:p@localhost:5432/logging")
	t.Setenv("RABBITMQ_URL", "amqp://guest:guest@rabbit:5672/")
	t.Setenv("BATCH_SIZE", "42")
	t.Setenv("POLL_INTERVAL", "2m")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DBWriteURL != "postgres://u:p@localhost:5432/logging" {
		t.Errorf("DBWriteURL = %q", cfg.DBWriteURL)
	}
	if cfg.RabbitMQURL != "amqp://guest:guest@rabbit:5672/" {
		t.Errorf("RabbitMQURL = %q", cfg.RabbitMQURL)
	}
	if cfg.BatchSize != 42 {
		t.Errorf("BatchSize = %d", cfg.BatchSize)
	}
	if cfg.PollInterval != 2*time.Minute {
		t.Errorf("PollInterval = %v", cfg.PollInterval)
	}
}

func TestLoad_invalidBatchSizeFallsBack(t *testing.T) {
	t.Chdir(t.TempDir())

	t.Setenv("DATABASE_WRITE_URL", "postgres://u:p@localhost:5432/db")
	t.Setenv("RABBITMQ_URL", "amqp://localhost/")
	t.Setenv("BATCH_SIZE", "not-a-number")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.BatchSize != 10 {
		t.Errorf("BatchSize fallback = %d, want 10", cfg.BatchSize)
	}
}
