package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"api-server/internal/config"
	"api-server/internal/db"
	"api-server/internal/elastic"
	"api-server/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()

	writePool, err := db.NewWritePool(ctx, cfg.DBWriteURL)
	if err != nil {
		log.Fatalf("connect to postgres write pool: %v", err)
	}
	defer writePool.Close()

	readPool, err := db.NewReadPool(ctx, cfg.DBReadURL)
	if err != nil {
		log.Fatalf("connect to postgres read pool: %v", err)
	}
	defer readPool.Close()

	esClient, err := elastic.NewClient(cfg.ESURL)
	if err != nil {
		log.Fatalf("connect to elasticsearch: %v", err)
	}

	r := gin.Default()
	r.POST("/logs", handler.Ingest(writePool))
	r.GET("/logs", handler.Search(esClient))
	r.GET("/health", handler.Health(writePool, readPool))

	log.Printf("api-server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
