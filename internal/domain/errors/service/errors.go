package service

import "errors"

var (
	// Couriers
	ErrInvalidName          = errors.New("invalid name")
	ErrInvalidPhone         = errors.New("invalid phone")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrInvalidTransportType = errors.New("invalid transport type")
	ErrCourierExists        = errors.New("courier already exists")
	ErrCourierNotFound      = errors.New("courier not found")
	ErrNoAvailableCouriers  = errors.New("no available couriers")
	// Delivery
	ErrDeliveryExists   = errors.New("delivery already exists")
	ErrDeliveryNotFound = errors.New("delivery not found")
	// Default
	ErrInternalError = errors.New("internal error")
)
