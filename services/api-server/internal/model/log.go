package model

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID          uuid.UUID              `json:"id"`
	Level       string                 `json:"level"        binding:"required,oneof=DEBUG INFO WARN ERROR FATAL"`
	Message     string                 `json:"message"      binding:"required"`
	ServiceName string                 `json:"service_name" binding:"required"`
	Timestamp   time.Time              `json:"timestamp"    binding:"required"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}
