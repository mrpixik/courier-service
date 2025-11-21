package dto

import "time"

type PingResponse struct {
	Message string `json:"message"`
}

// GetCourierResponse модель данных для получения профиля
type GetCourierResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// CreateCourierResponse ответ на создание профиля
type CreateCourierResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

// UpdateCourierResponse запрос на обновление профиля
type UpdateCourierResponse struct {
	Message string `json:"message"`
}

// DeleteCourierResponse запрос на обновление профиля
type DeleteCourierResponse struct {
	Message string `json:"message"`
}
