package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-server/internal/db"
	"api-server/internal/model"
)

type ingestRequest struct {
	Level       string                 `json:"level"        binding:"required,oneof=DEBUG INFO WARN ERROR FATAL"`
	Message     string                 `json:"message"      binding:"required"`
	ServiceName string                 `json:"service_name" binding:"required"`
	Timestamp   time.Time              `json:"timestamp"    binding:"required"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func Ingest(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ingestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l := &model.Log{
			ID:          uuid.New(),
			Level:       req.Level,
			Message:     req.Message,
			ServiceName: req.ServiceName,
			Timestamp:   req.Timestamp,
			Metadata:    req.Metadata,
			CreatedAt:   time.Now().UTC(),
		}

		if err := db.InsertLogWithOutbox(c.Request.Context(), pool, l); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist log"})
			return
		}

		c.JSON(http.StatusCreated, l)
	}
}
