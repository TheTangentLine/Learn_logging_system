package consumer

import (
	"context"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ErrPoison indicates the message should not be requeued (e.g. malformed JSON).
var ErrPoison = errors.New("poison message")

const (
	exchangeName = "logs"
	exchangeType = "direct"
	routingKey   = "logs"
	queueName    = "logs"

	prefetchCount = 50
	consumerTag   = "logs-es-sync"
)

type Handler func(ctx context.Context, body []byte) error

type Consumer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func New(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("queue declare: %w", err)
	}

	if err := ch.QueueBind(q.Name, routingKey, exchangeName, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("queue bind: %w", err)
	}

	if err := ch.Qos(prefetchCount, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("qos: %w", err)
	}

	return &Consumer{conn: conn, ch: ch, queue: q}, nil
}

// Consume registers a consumer and processes deliveries until ctx is cancelled or the channel closes.
func (c *Consumer) Consume(ctx context.Context, h Handler) error {
	msgs, err := c.ch.Consume(
		c.queue.Name,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			_ = c.ch.Cancel(consumerTag, false)
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}

			procCtx := ctx
			if procCtx.Err() != nil {
				_ = d.Nack(false, true)
				continue
			}

			if err := h(procCtx, d.Body); err != nil {
				if errors.Is(err, ErrPoison) {
					_ = d.Nack(false, false)
				} else {
					_ = d.Nack(false, true)
				}
				continue
			}
			_ = d.Ack(false)
		}
	}
}

func (c *Consumer) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
