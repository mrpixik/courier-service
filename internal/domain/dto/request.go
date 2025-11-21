package dto

// GetCourierRequest запрос за получение данных о курьере
type GetCourierRequest struct {
	Id int `json:"id"`
}

// CreateCourierRequest запрос на создание профиля
type CreateCourierRequest struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

// UpdateCourierRequest запрос на обновление профиля
type UpdateCourierRequest struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

// DeleteCourierRequest запрос за получение данных о курьере
type DeleteCourierRequest struct {
	Id int `json:"id"`
}
