package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"

	"rmq-consumer/internal/model"
)

const logsIndex = "logs"

// ErrInvalidPayload means the body cannot be indexed and must not be requeued.
var ErrInvalidPayload = errors.New("invalid log payload")

type Syncer struct {
	es *es8.Client
}

func NewSyncer(url string) (*Syncer, error) {
	c, err := es8.NewClient(es8.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}
	return &Syncer{es: c}, nil
}

// ParseAndEncodeLog unmarshals ingest JSON and returns canonical UTF-8 JSON for indexing alongside the document id.
func ParseAndEncodeLog(body []byte) (payload []byte, docID string, err error) {
	var doc model.Log
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrInvalidPayload, err)
	}
	if doc.ID == uuid.Nil || doc.Level == "" || doc.Message == "" || doc.ServiceName == "" || doc.Timestamp.IsZero() || doc.CreatedAt.IsZero() {
		return nil, "", fmt.Errorf("%w: missing required fields", ErrInvalidPayload)
	}
	payload, err = json.Marshal(doc)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrInvalidPayload, err)
	}
	return payload, doc.ID.String(), nil
}

// Upsert validates body as Log JSON and indexes it with document id equal to log id (idempotent).
func (s *Syncer) Upsert(ctx context.Context, body []byte) error {
	payload, docID, err := ParseAndEncodeLog(body)
	if err != nil {
		return err
	}

	res, err := s.es.Index(
		logsIndex,
		bytes.NewReader(payload),
		s.es.Index.WithDocumentID(docID),
		s.es.Index.WithRefresh("false"),
		s.es.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("elasticsearch index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		raw, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch error response: %s", raw)
	}
	return nil
}
