package strategies

import (
	"context"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/dto/kafka/order"
	"service-order-avito/internal/domain/model"
)

type createService interface {
	Assign(context.Context, *dto.AssignDeliveryRequest) (*dto.AssignDeliveryResponse, error)
}

type createStrategy struct {
	service createService
}

func NewCreateStrategy(service createService) *createStrategy {
	return &createStrategy{service: service}
}

func (cs *createStrategy) Process(ctx context.Context, orderId string) (*order.ProcessedEvent, error) {
	req := &dto.AssignDeliveryRequest{OrderId: orderId}

	res, err := cs.service.Assign(ctx, req)
	if err != nil {
		return nil, err
	}

	return &order.ProcessedEvent{
		OrderId:   res.OrderId,
		Status:    model.StatusAssigned,
		CourierId: res.CourierId,
	}, nil
}
