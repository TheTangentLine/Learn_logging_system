package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"relay-worker/internal/config"
	"relay-worker/internal/db"
	"relay-worker/internal/poller"
	"relay-worker/internal/producer"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	writePool, err := db.NewWritePool(ctx, cfg.DBWriteURL)
	if err != nil {
		log.Fatalf("connect write pool: %v", err)
	}
	defer writePool.Close()

	listenConn, err := db.NewListenConn(ctx, cfg.DBWriteURL)
	if err != nil {
		log.Fatalf("connect listen conn: %v", err)
	}
	defer listenConn.Close(ctx)

	prod, err := producer.New(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("connect rabbitmq: %v", err)
	}
	defer prod.Close()

	p := poller.New(writePool, listenConn, prod, cfg.BatchSize, cfg.PollInterval)

	log.Println("relay-worker: started")
	if err := p.Run(ctx); err != nil {
		log.Fatalf("poller error: %v", err)
	}
	log.Println("relay-worker: stopped")
}
