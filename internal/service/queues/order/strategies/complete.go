package strategies

import (
	"context"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/dto/kafka/order"
	"service-order-avito/internal/domain/model"
)

type completeService interface {
	Complete(context.Context, *dto.CompleteDeliveryRequest) (*dto.CompleteDeliveryResponse, error)
}

type completeStrategy struct {
	service completeService
}

func NewCompleteStrategy(service completeService) *completeStrategy {
	return &completeStrategy{service: service}
}

func (cs *completeStrategy) Process(ctx context.Context, orderId string) (*order.ProcessedEvent, error) {
	req := &dto.CompleteDeliveryRequest{OrderId: orderId}

	res, err := cs.service.Complete(ctx, req)
	if err != nil {
		return nil, err
	}

	return &order.ProcessedEvent{
		OrderId:   res.OrderId,
		Status:    model.StatusCompleted,
		CourierId: res.CourierId,
	}, nil
}
