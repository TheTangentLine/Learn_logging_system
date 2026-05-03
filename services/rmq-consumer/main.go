package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"rmq-consumer/internal/config"
	"rmq-consumer/internal/consumer"
	"rmq-consumer/internal/elastic"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	syncer, err := elastic.NewSyncer(cfg.ESURL)
	if err != nil {
		log.Fatalf("elasticsearch: %v", err)
	}

	cons, err := consumer.New(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("rabbitmq: %v", err)
	}
	defer cons.Close()

	handler := func(c context.Context, body []byte) error {
		err := syncer.Upsert(c, body)
		if err != nil && errors.Is(err, elastic.ErrInvalidPayload) {
			log.Printf("invalid message discarded: %v", err)
			return fmt.Errorf("%w", consumer.ErrPoison)
		}
		return err
	}

	log.Println("rmq-consumer: started")
	if err := cons.Consume(ctx, handler); err != nil && err != ctx.Err() {
		log.Fatalf("consume stopped: %v", err)
	}
	log.Println("rmq-consumer: stopped")
}
