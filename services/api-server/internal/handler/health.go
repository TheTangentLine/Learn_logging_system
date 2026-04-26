package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Health pings both the write and read Postgres pools and reports their status.
func Health(write, read *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		status := gin.H{"write_db": "ok", "read_db": "ok"}
		code := http.StatusOK

		if err := write.Ping(ctx); err != nil {
			status["write_db"] = err.Error()
			code = http.StatusServiceUnavailable
		}
		if err := read.Ping(ctx); err != nil {
			status["read_db"] = err.Error()
			code = http.StatusServiceUnavailable
		}

		c.JSON(code, status)
	}
}
