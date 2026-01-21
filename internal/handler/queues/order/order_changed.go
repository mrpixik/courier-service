package order

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"service-order-avito/internal/adapters/logger"
	"service-order-avito/internal/domain/dto/kafka/order"
	"service-order-avito/internal/domain/errors/service"
	"time"

	"github.com/IBM/sarama"
)

const (
	BASE_DELAY  = 100
	MAX_RETRIES = 3
)

type usecase interface {
	Process(context.Context, *order.Event) (*order.ProcessedEvent, error)
}

type orderServiceGRPCGateway interface {
	GetOrderStatusById(ctx context.Context, id string) (string, error)
}

type totalGatewayRetriesCounter interface {
	IncTotalGatewayRetries()
}

type handler struct {
	l           logger.LoggerAdapter
	og          orderServiceGRPCGateway
	uc          usecase
	retrCounter totalGatewayRetriesCounter
}

func NewOrderChangedHandler(l logger.LoggerAdapter, og orderServiceGRPCGateway, uc usecase, retrCounter totalGatewayRetriesCounter) *handler {
	return &handler{l: l, og: og, uc: uc, retrCounter: retrCounter}
}

func (h *handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	op := "order.changed.handler: "
CLAIM_MSSG_LOOP:
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

		// check if order status is still actual
		for i := range MAX_RETRIES {

			actualStatus, err := h.og.GetOrderStatusById(ctx, event.OrderID)
			if err != nil {
				h.retrCounter.IncTotalGatewayRetries()
				// exponential backoff + jitter
				if i != MAX_RETRIES-1 {
					delay := i*i*BASE_DELAY + rand.Intn(20)
					time.Sleep(time.Duration(delay) * time.Millisecond)
				}

			} else if actualStatus != event.Status {
				h.l.Info("order.changed handler: order's status changed",
					"id", event.OrderID,
					"prev_status", event.Status,
					"actual_status", actualStatus,
				)
				continue CLAIM_MSSG_LOOP
			}
		}

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
