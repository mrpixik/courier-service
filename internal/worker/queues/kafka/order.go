package kafka

import (
	"context"
	"github.com/IBM/sarama"
)

type logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

type orderConsumerWorker struct {
	l       logger
	client  sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topic   string
}

func NewOrderConsumerWorker(l logger, client sarama.ConsumerGroup, handler sarama.ConsumerGroupHandler, topic string) *orderConsumerWorker {
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
