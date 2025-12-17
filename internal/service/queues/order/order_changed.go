package order

import (
	"context"
	"service-order-avito/internal/domain/dto/kafka/order"
	"service-order-avito/internal/domain/errors/service"
)

type orderChangedService struct {
	serv delService
	fab  *orderStrategyFactory
}

func NewOrderChangedService(serv delService) *orderChangedService {
	return &orderChangedService{serv: serv, fab: NewOrderStrategyFactory(serv)}
}

func (os *orderChangedService) Process(ctx context.Context, event *order.Event) (*order.ProcessedEvent, error) {
	strategy := os.fab.SelectStrategy(event.Status)
	if strategy == nil {
		return nil, service.ErrUnknownOrderStatus
	}

	return strategy.Process(ctx, event.OrderID)
}
