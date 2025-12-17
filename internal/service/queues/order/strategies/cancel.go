package strategies

import (
	"context"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/dto/kafka/order"
	"service-order-avito/internal/domain/model"
)

type cancelService interface {
	Unassign(context.Context, *dto.UnassignDeliveryRequest) (*dto.UnassignDeliveryResponse, error)
}

type cancelStrategy struct {
	service cancelService
}

func NewCancelStrategy(service cancelService) *cancelStrategy {
	return &cancelStrategy{service: service}
}

func (cs *cancelStrategy) Process(ctx context.Context, orderId string) (*order.ProcessedEvent, error) {
	req := &dto.UnassignDeliveryRequest{OrderId: orderId}

	res, err := cs.service.Unassign(ctx, req)
	if err != nil {
		return nil, err
	}

	return &order.ProcessedEvent{
		OrderId:   res.OrderId,
		Status:    model.StatusUnassigned,
		CourierId: res.CourierId,
	}, nil
}
