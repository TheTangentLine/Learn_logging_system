package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-server/internal/elastic"
)

func Search(es *elastic.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		level := c.Query("level")
		service := c.Query("service")

		logs, err := es.Search(c.Request.Context(), q, level, service)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}

		c.JSON(http.StatusOK, logs)
	}
}
