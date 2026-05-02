package producer

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "logs"
	exchangeType = "direct"
	routingKey   = "logs"
)

// Producer holds an AMQP connection and a publisher-confirm channel.
type Producer struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	confirm chan amqp.Confirmation
}

// New connects to RabbitMQ, declares the exchange, and enables publisher confirms.
func New(url string) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	// Declare a durable direct exchange (idempotent).
	if err := ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	// Enable publisher confirms on this channel.
	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("enable publisher confirms: %w", err)
	}

	confirmCh := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	return &Producer{conn: conn, ch: ch, confirm: confirmCh}, nil
}

// Publish sends payload to the exchange and blocks until a broker confirm is received.
// Returns an error if the broker nacks the message or if the channel is closed.
func (p *Producer) Publish(ctx context.Context, payload []byte) error {
	err := p.ch.PublishWithContext(ctx, exchangeName, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         payload,
	})
	if err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	select {
	case confirm, ok := <-p.confirm:
		if !ok {
			return fmt.Errorf("confirm channel closed")
		}
		if !confirm.Ack {
			return fmt.Errorf("broker nacked message")
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close shuts down the channel and connection gracefully.
func (p *Producer) Close() {
	p.ch.Close()
	p.conn.Close()
}
