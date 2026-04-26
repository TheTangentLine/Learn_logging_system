package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupIngestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery()) // catches nil-pool panic on the happy path in tests
	// nil pool is safe for validation-failure tests; Recovery handles the nil-pool panic on valid input
	r.POST("/logs", Ingest(nil))
	return r
}

func TestIngest_BadJSON(t *testing.T) {
	r := setupIngestRouter()

	req := httptest.NewRequest(http.MethodPost, "/logs", strings.NewReader("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestIngest_MissingRequiredFields(t *testing.T) {
	r := setupIngestRouter()

	// Empty JSON object — all required fields absent
	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/logs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for missing fields, got %d", w.Code)
	}
}

func TestIngest_InvalidLevel(t *testing.T) {
	r := setupIngestRouter()

	body := `{
		"level":        "VERBOSE",
		"message":      "test message",
		"service_name": "my-service",
		"timestamp":    "2026-01-01T00:00:00Z"
	}`
	req := httptest.NewRequest(http.MethodPost, "/logs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid level, got %d", w.Code)
	}
}

func TestIngest_ValidLevels(t *testing.T) {
	levels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			r := setupIngestRouter()

			body := `{
				"level":        "` + level + `",
				"message":      "test message",
				"service_name": "my-service",
				"timestamp":    "2026-01-01T00:00:00Z"
			}`
			req := httptest.NewRequest(http.MethodPost, "/logs", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Validation passes (201 or 500 from nil pool — not 400)
			if w.Code == http.StatusBadRequest {
				t.Errorf("level %q should be valid but got 400", level)
			}
		})
	}
}
