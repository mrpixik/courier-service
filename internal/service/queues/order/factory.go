package order

import (
	"context"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/dto/kafka/order"
	"service-order-avito/internal/service/queues/order/strategies"
)

type delService interface {
	Assign(context.Context, *dto.AssignDeliveryRequest) (*dto.AssignDeliveryResponse, error)
	Unassign(context.Context, *dto.UnassignDeliveryRequest) (*dto.UnassignDeliveryResponse, error)
	Complete(context.Context, *dto.CompleteDeliveryRequest) (*dto.CompleteDeliveryResponse, error)
}

type orderChangedStrategyFabric interface {
	Process(context.Context, string) (*order.ProcessedEvent, error)
}

type orderStrategyFactory struct {
	service delService
}

func NewOrderStrategyFactory(service delService) *orderStrategyFactory {
	return &orderStrategyFactory{service: service}
}

func (of *orderStrategyFactory) SelectStrategy(status string) orderChangedStrategyFabric {

	switch status {
	case order.StatusCreated:
		return strategies.NewCreateStrategy(of.service)
	case order.StatusCancelled:
		return strategies.NewCancelStrategy(of.service)
	case order.StatusCompleted:
		return strategies.NewCompleteStrategy(of.service)
	default:
		return nil
	}
}
