package elastic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

func TestParseAndEncodeLog_invalidJSON(t *testing.T) {
	_, _, err := ParseAndEncodeLog([]byte("{not-json"))
	if err == nil || !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload wrap, got %v", err)
	}
}

func TestParseAndEncodeLog_missingRequiredFields(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{name: "empty_level", body: `{"id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee","level":"","message":"m","service_name":"x","timestamp":"2026-05-03T12:00:00Z","created_at":"2026-05-03T12:00:00Z"}`},
		{name: "nil_uuid", body: `{"id":"00000000-0000-0000-0000-000000000000","level":"INFO","message":"m","service_name":"x","timestamp":"2026-05-03T12:00:00Z","created_at":"2026-05-03T12:00:00Z"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ParseAndEncodeLog([]byte(tc.body))
			if err == nil || !errors.Is(err, ErrInvalidPayload) {
				t.Fatalf("expected ErrInvalidPayload, got %v", err)
			}
		})
	}
}

func TestParseAndEncodeLog_success(t *testing.T) {
	const raw = `{"id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee","level":"INFO","message":"hello","service_name":"svc","timestamp":"2026-05-03T12:00:00Z","created_at":"2026-05-03T12:01:00Z","metadata":{"k":1}}`
	payload, docID, err := ParseAndEncodeLog([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if docID != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Fatalf("docID = %s", docID)
	}
	if len(payload) == 0 || !strings.Contains(string(payload), "hello") {
		t.Fatalf("payload = %s", payload)
	}
}

func TestUpsert_success_mockServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"result":"created"}`)
	}))
	t.Cleanup(srv.Close)

	c, err := es8.NewClient(es8.Config{Addresses: []string{srv.URL}})
	if err != nil {
		t.Fatal(err)
	}
	syn := &Syncer{es: c}

	body := []byte(`{"id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee","level":"INFO","message":"hi","service_name":"svc","timestamp":"2026-05-03T12:00:00Z","created_at":"2026-05-03T12:00:01Z"}`)
	if err := syn.Upsert(context.Background(), body); err != nil {
		t.Fatal(err)
	}
}
