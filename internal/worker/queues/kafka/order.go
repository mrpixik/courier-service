package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"service-order-avito/internal/adapters/logger"
)

type orderConsumerWorker struct {
	l       logger.LoggerAdapter
	client  sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topic   string
}

func NewOrderConsumerWorker(l logger.LoggerAdapter, client sarama.ConsumerGroup, handler sarama.ConsumerGroupHandler, topic string) *orderConsumerWorker {
	return &orderConsumerWorker{
		l:       l,
		client:  client,
		handler: handler,
		topic:   topic,
	}
}

func (w *orderConsumerWorker) Start(ctx context.Context) {
	for {
		err := w.client.Consume(ctx, []string{w.topic}, w.handler)
		if err != nil {
			w.l.Error("consume error",
				"error", err.Error(),
			)
		}

		select {
		case <-ctx.Done():
			w.l.Info("delivery monitor worker gracefully stopped")
			return
		}
	}
}
