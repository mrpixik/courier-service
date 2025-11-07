package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"service-order-avito/internal/http/server/dto"
	"service-order-avito/internal/service"
	"strconv"
)

type CourierService interface {
	CreateCourier(context.Context, *dto.CourierCreateRequest) (int, error)
	GetCourier(context.Context, int) (*dto.Courier, error)
	GetAllCouriers(context.Context) ([]dto.Courier, error)
	UpdateCourier(context.Context, *dto.CourierUpdateRequest) error
}

type courierHandler struct {
	service CourierService
}

func NewCourierHandler(service CourierService) *courierHandler {
	return &courierHandler{service: service}
}

// Я не всоем понял, стоит ли добавлять такую же систему, как и в хендлерах /ping и /healthcheck
// которая бы проверяла завершение контекста через select.
func (ch *courierHandler) Post(w http.ResponseWriter, r *http.Request) {
	var req dto.CourierCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	id, err := ch.service.CreateCourier(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCourierAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, service.ErrInvalidName):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, service.ErrInvalidPhone):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, service.ErrInvalidStatus):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, ErrInternalError, http.StatusInternalServerError)
		}
		return
	}
	res := dto.CourierCreateResponse{
		Id:      id,
		Message: "courier's profile created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (ch *courierHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, ErrInvalidCourierId, http.StatusBadRequest)
		return
	}

	var courier *dto.Courier
	courier, err = ch.service.GetCourier(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrCourierNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, ErrInternalError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courier)
}

func (ch *courierHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	couriers, err := ch.service.GetAllCouriers(r.Context())
	if err != nil {
		http.Error(w, ErrInternalError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(couriers)
}

func (ch *courierHandler) Put(w http.ResponseWriter, r *http.Request) {
	var req dto.CourierUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	err := ch.service.UpdateCourier(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCourierAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, service.ErrInvalidName):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, service.ErrInvalidPhone):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, service.ErrInvalidStatus):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	res := dto.CourierUpdateResponse{
		Message: "courier's profile updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
