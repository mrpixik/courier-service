package dto

import "time"

type PingResponse struct {
	Message string `json:"message"`
}

// GetCourierResponse модель данных для получения профиля
type GetCourierResponse struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	TransportType string    `json:"transport-type"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// CreateCourierResponse ответ на создание профиля
type CreateCourierResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

// UpdateCourierResponse ответ на обновление профиля
type UpdateCourierResponse struct {
	Message string `json:"message"`
}

// DeleteCourierResponse ответ на удаление профиля
type DeleteCourierResponse struct {
	Message string `json:"message"`
}

// AssignDeliveryResponse запрос на назначение заказа
type AssignDeliveryResponse struct {
	CourierId        int       `json:"courier_id"`
	OrderId          string    `json:"order_id"`
	TransportType    string    `json:"transport_type"`
	DeliveryDeadline time.Time `json:"delivery_deadline"`
}

// UnassignDeliveryResponse запрос на снятие заказа
type UnassignDeliveryResponse struct {
	OrderId   string `json:"order_id"`
	Status    string `json:"status"`
	CourierId int    `json:"courier_id"`
}

// CompleteDeliveryResponse запрос на снятие заказа
type CompleteDeliveryResponse struct {
	OrderId   string `json:"order_id"`
	Status    string `json:"status"`
	CourierId int    `json:"courier_id"`
}
