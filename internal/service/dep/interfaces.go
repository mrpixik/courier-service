package dep

import (
	"context"
	"service-order-avito/internal/domain/model"
	"time"
)

// mockgen -source="internal/service/dep/interfaces.go" -destination="internal/service/dep/mocks/mock_repositories.go"
type TransactionManager interface {
	Begin(context.Context, func(context.Context) error) error
}

type CourierRepository interface {
	Create(context.Context, model.Courier) (int, error)
	GetById(context.Context, int) (model.Courier, error)
	GetAll(context.Context) ([]model.Courier, error)
	Update(context.Context, model.Courier) error
	UpdateStatusManyById(context.Context, ...int) error
	DeleteById(context.Context, int) error
	GetAvailable(context.Context) (model.Courier, error)
}

type DeliveryRepository interface {
	Create(context.Context, model.Delivery) (int, error)
	GetByOrderId(context.Context, string) (model.Delivery, error)
	GetAllCompleted(context.Context) ([]model.Delivery, error)
	DeleteByOrderId(context.Context, string) error
	DeleteManyById(context.Context, ...int) error
}

type DeliveryTimeCalculator interface {
	Calculate(transportType string) time.Time
}
