package model

import (
	"time"

	"github.com/google/uuid"
)

// Log mirrors api-server ingest JSON / outbox payload for indexing.
type Log struct {
	ID          uuid.UUID              `json:"id"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	ServiceName string                 `json:"service_name"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}
