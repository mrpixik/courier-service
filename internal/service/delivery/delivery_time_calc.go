package delivery

import (
	"service-order-avito/internal/domain/model"
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
	case model.TransportTypeFoot:
		return time.Now().Add(footDeliveryTime)
	case model.TransportTypeScooter:
		return time.Now().Add(scooterDeliveryTime)
	case model.TransportTypeCar:
		return time.Now().Add(carDeliveryTime)
	default: // такого быть не может
		return time.Now().Add(footDeliveryTime)
	}
}
