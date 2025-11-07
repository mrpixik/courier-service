package dto

import "time"

type PingResponse struct {
	Message string `json:"message"`
}

// Courier модель данных для получения профиля
type Courier struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// CourierCreateRequest запрос на создание профиля
type CourierCreateRequest struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

// CourierCreateResponse ответ на создание профиля
type CourierCreateResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

// CourierUpdateRequest запрос на обновление профиля
type CourierUpdateRequest struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

// CourierUpdateResponse запрос на обновление профиля
type CourierUpdateResponse struct {
	Message string `json:"message"`
}
