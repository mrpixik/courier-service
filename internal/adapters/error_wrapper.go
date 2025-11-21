package adapters

import (
	"encoding/json"
	"net/http"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/errors/server"
	"service-order-avito/internal/domain/errors/service"
)

// Написал этот функционал, чтобы не передавать ошибки с уровня репозитория наверх к уровню контроллеров
// и чтобы избежать огромных структур if/else в сервисном слое при проверке соответствия ошибок.
// Я понимаю, что в целом, передавать можно и так делают в проде, и в высоконагруженных сервисах это оправданное решение,
// но так как это учебный проект, решил сделать все максимально правильно с точки зрения теории
// P.S. можно подумать побольше и, например, переделать это в фабрику
var repoToServiceMap = map[error]error{
	repository.ErrCourierExists:   service.ErrCourierExists,
	repository.ErrCourierNotFound: service.ErrCourierNotFound,
	repository.ErrInternalError:   service.ErrInternalError,
}

func ErrUnwrapRepoToService(err error) error {
	if err == nil {
		return nil
	}

	if mapped, ok := repoToServiceMap[err]; ok {
		return mapped
	}

	return service.ErrInternalError
}

type errorMeta struct {
	Message string
	Status  int
}

var serviceErrorMap = map[error]errorMeta{
	service.ErrInvalidName:     {server.ErrInvalidCourierName, http.StatusNotFound},
	service.ErrInvalidStatus:   {server.ErrInvalidCourierStatus, http.StatusNotFound},
	service.ErrInvalidPhone:    {server.ErrInvalidCourierPhone, http.StatusConflict},
	service.ErrCourierExists:   {server.ErrCourierExists, http.StatusBadRequest},
	service.ErrCourierNotFound: {server.ErrCourierNotFound, http.StatusNotFound},
	service.ErrInternalError:   {server.ErrInternalError, http.StatusInternalServerError},
}

// WriteServiceError принимает ошибку уровня service и пишет ошибку уровня контроллера в ResponseWriter
func WriteServiceError(w http.ResponseWriter, err error) {
	meta, ok := serviceErrorMap[err]
	if !ok {
		meta = serviceErrorMap[service.ErrInternalError]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.Status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Message: meta.Message,
		},
	})
}

func WriteError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Message: message,
		},
	})
}
