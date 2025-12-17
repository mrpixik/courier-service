package dto

// GetCourierRequest запрос за получение данных о курьере
type GetCourierRequest struct {
	Id int `json:"id"`
}

// CreateCourierRequest запрос на создание профиля курьера
type CreateCourierRequest struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	TransportType string `json:"transport_type"`
}

// UpdateCourierRequest запрос на обновление данных курьера
type UpdateCourierRequest struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	TransportType string `json:"transport_type"`
}

// DeleteCourierRequest запрос за удаление данных о курьере
type DeleteCourierRequest struct {
	Id int `json:"id"`
}

// AssignDeliveryRequest запрос на назначение доставки
type AssignDeliveryRequest struct {
	OrderId string `json:"order_id"`
}

// UnassignDeliveryRequest запрос на завершение доставки
type UnassignDeliveryRequest struct {
	OrderId string `json:"order_id"`
}

// CompleteDeliveryRequest запрос на завершение доставки (без удаления из таблицы)
type CompleteDeliveryRequest struct {
	OrderId string `json:"order_id"`
}
