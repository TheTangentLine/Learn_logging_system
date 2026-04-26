package config

import (
	"os"
	"testing"
)

func clearEnv() {
	os.Unsetenv("DATABASE_WRITE_URL")
	os.Unsetenv("DATABASE_READ_URL")
	os.Unsetenv("PORT")
	os.Unsetenv("ELASTICSEARCH_URL")
}

func TestLoad_MissingWriteURL(t *testing.T) {
	clearEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when DATABASE_WRITE_URL is not set, got nil")
	}
}

func TestLoad_Defaults(t *testing.T) {
	clearEnv()
	os.Setenv("DATABASE_WRITE_URL", "postgres://localhost/test")
	defer clearEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected default Port=8080, got %q", cfg.Port)
	}
	if cfg.ESURL != "http://localhost:9200" {
		t.Errorf("expected default ESURL=http://localhost:9200, got %q", cfg.ESURL)
	}
	if cfg.DBReadURL != cfg.DBWriteURL {
		t.Errorf("expected DBReadURL to fall back to DBWriteURL, got %q", cfg.DBReadURL)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	clearEnv()
	os.Setenv("DATABASE_WRITE_URL", "postgres://primary/db")
	os.Setenv("DATABASE_READ_URL", "postgres://replica/db")
	os.Setenv("PORT", "9090")
	os.Setenv("ELASTICSEARCH_URL", "http://es:9200")
	defer clearEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DBWriteURL != "postgres://primary/db" {
		t.Errorf("unexpected DBWriteURL: %q", cfg.DBWriteURL)
	}
	if cfg.DBReadURL != "postgres://replica/db" {
		t.Errorf("unexpected DBReadURL: %q", cfg.DBReadURL)
	}
	if cfg.Port != "9090" {
		t.Errorf("unexpected Port: %q", cfg.Port)
	}
	if cfg.ESURL != "http://es:9200" {
		t.Errorf("unexpected ESURL: %q", cfg.ESURL)
	}
}
