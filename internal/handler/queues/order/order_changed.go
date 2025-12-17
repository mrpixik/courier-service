package order

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"service-order-avito/internal/domain/dto/kafka/order"
	"service-order-avito/internal/domain/errors/service"
)

type logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

type usecase interface {
	Process(context.Context, *order.Event) (*order.ProcessedEvent, error)
}

type handler struct {
	l  logger
	uc usecase
}

func NewOrderChangedHandler(l logger, uc usecase) *handler {
	return &handler{l: l, uc: uc}
}

func (h *handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	op := "order.changed.handler: "
	for dtoMsg := range claim.Messages() {
		ctx := sess.Context()
		h.l.Info("order.changed handler: received message",
			"key", string(dtoMsg.Key),
			"value", string(dtoMsg.Value),
			"partition", int(dtoMsg.Partition),
			"offset", dtoMsg.Offset,
		)

		var event order.Event
		err := json.Unmarshal(dtoMsg.Value, &event)
		if err != nil {
			h.l.Error(op+"received bad message",
				"error", err.Error(),
			)
			sess.MarkMessage(dtoMsg, "")
			continue
		}

		//TODO: add gRPC to check if order status is still actual

		res, err := h.uc.Process(ctx, &event)
		if err != nil {
			if errors.Is(err, service.ErrUnknownOrderStatus) {
				continue
			}
			h.l.Error(op+"failed process order",
				"error", err.Error(),
			)
			continue
		}

		h.l.Info(op+" message processed",
			"order_id", res.OrderId,
			"status", res.Status,
			"courier_id", res.CourierId,
		)

		sess.MarkMessage(dtoMsg, "")

	}
	return nil
}
