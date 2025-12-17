package model

import "time"

const (
	StatusAvailable = "available"
	StatusBusy      = "busy"
	StatusPaused    = "paused"

	TransportTypeFoot    = "on_foot"
	TransportTypeScooter = "scooter"
	TransportTypeCar     = "car"
)

// Courier сущность из таблицы couriers
type Courier struct {
	Id              int
	Name            string
	Phone           string
	Status          string // available | busy | paused
	TransportType   string // on_foot | scooter | car
	TotalDeliveries int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
