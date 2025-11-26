package repository

import "errors"

var (
	// Courier
	ErrCourierExists       = errors.New("courier already exists")
	ErrNoAvailableCouriers = errors.New("no available couriers")
	ErrCourierNotFound     = errors.New("courier not found")
	// Delivery
	ErrDeliveryExists   = errors.New("delivery already exists")
	ErrDeliveryNotFound = errors.New("delivery not found")
	// Default
	ErrInternalError = errors.New("internal error")
)
