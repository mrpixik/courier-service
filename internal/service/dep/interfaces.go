package dep

import (
	"context"
	"service-order-avito/internal/domain"
	"time"
)

// mockgen -source="internal/service/dep/interfaces.go" -destination="internal/service/dep/mocks/mock_repositories.go"
type TransactionManager interface {
	Begin(context.Context, func(context.Context) error) error
}

type CourierRepository interface {
	Create(context.Context, domain.Courier) (int, error)
	GetById(context.Context, int) (domain.Courier, error)
	GetAll(context.Context) ([]domain.Courier, error)
	Update(context.Context, domain.Courier) error
	UpdateStatusManyById(context.Context, ...int) error
	DeleteById(context.Context, int) error
	GetAvailable(context.Context) (domain.Courier, error)
}

type DeliveryRepository interface {
	Create(context.Context, domain.Delivery) (int, error)
	GetByOrderId(context.Context, string) (domain.Delivery, error)
	GetAllCompleted(context.Context) ([]domain.Delivery, error)
	DeleteByOrderId(context.Context, string) error
	DeleteManyById(context.Context, ...int) error
}

type DeliveryTimeCalculator interface {
	Calculate(transportType string) time.Time
}
