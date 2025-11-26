package delivery

import (
	"service-order-avito/internal/domain"
	"time"
)

const (
	footDeliveryTime    = 30 * time.Minute
	scooterDeliveryTime = 15 * time.Minute
	carDeliveryTime     = 5 * time.Minute
)

type deliveryTimeFactory struct{}

func NewDeliveryTimeFactory() *deliveryTimeFactory {
	return &deliveryTimeFactory{}
}

func (dtf *deliveryTimeFactory) Calculate(transportType string) time.Time {
	switch transportType {
	case domain.TransportTypeFoot:
		return time.Now().Add(footDeliveryTime)
	case domain.TransportTypeScooter:
		return time.Now().Add(scooterDeliveryTime)
	case domain.TransportTypeCar:
		return time.Now().Add(carDeliveryTime)
	default: // такого быть не может
		return time.Now().Add(footDeliveryTime)
	}
}
