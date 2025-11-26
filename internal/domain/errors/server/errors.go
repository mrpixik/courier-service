package server

const (
	// Courier
	ErrInvalidCourierId     = "invalid courier's id"
	ErrInvalidCourierName   = "invalid courier's name"
	ErrInvalidCourierStatus = "invalid courier's status"
	ErrInvalidCourierPhone  = "invalid courier's phone"
	ErrInvalidTransportType = "invalid transport type"
	ErrCourierExists        = "courier with this parameters already exists"
	ErrCourierNotFound      = "courier not found"
	ErrNoAvailableCouriers  = "no available couriers"
	// Delivery
	ErrDeliveryExists   = "this delivery already exists"
	ErrDeliveryNotFound = "delivery not found"
	// Default
	ErrRequestCanceled = "request canceled"
	ErrInvalidJSON     = "invalid JSON"
	ErrInternalError   = "internal error"
)
