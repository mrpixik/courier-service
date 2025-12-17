package courier

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"service-order-avito/internal/adapters"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
	"strconv"
)

// mockgen -source="internal/http/server/handler/courier/courier.go" -destination="internal/http/server/handler/courier/mocks/mock_courier_service.go"
type сourierService interface {
	CreateCourier(context.Context, *dto.CreateCourierRequest) (*dto.CreateCourierResponse, error)
	GetCourier(context.Context, *dto.GetCourierRequest) (*dto.GetCourierResponse, error)
	GetAllCouriers(context.Context) ([]dto.GetCourierResponse, error)
	UpdateCourier(context.Context, *dto.UpdateCourierRequest) error
	DeleteCourier(context.Context, *dto.DeleteCourierRequest) error
}

type courierHandler struct {
	service сourierService
}

func NewCourierHandler(service сourierService) *courierHandler {
	return &courierHandler{service: service}
}

func (ch *courierHandler) Post(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		adapters.WriteError(w, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	res, err := ch.service.CreateCourier(r.Context(), &req)
	if err != nil {
		adapters.WriteServiceError(w, err)
		return
	}
	res.Message = "courier's profile created successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

func (ch *courierHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		adapters.WriteError(w, server.ErrInvalidCourierId, http.StatusBadRequest)
		return
	}

	courierReq := &dto.GetCourierRequest{
		Id: id,
	}
	courier, err := ch.service.GetCourier(r.Context(), courierReq)

	if err != nil {
		adapters.WriteServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(courier)
}

func (ch *courierHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	couriers, err := ch.service.GetAllCouriers(r.Context())
	if err != nil {
		adapters.WriteServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(couriers)
}

func (ch *courierHandler) Put(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		adapters.WriteError(w, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	err := ch.service.UpdateCourier(r.Context(), &req)
	if err != nil {
		adapters.WriteServiceError(w, err)
		return
	}
	res := dto.UpdateCourierResponse{
		Message: "courier's profile updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (ch *courierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		adapters.WriteError(w, server.ErrInvalidCourierId, http.StatusBadRequest)
		return
	}

	courierReq := &dto.DeleteCourierRequest{Id: id}
	err = ch.service.DeleteCourier(r.Context(), courierReq)

	if err != nil {
		adapters.WriteServiceError(w, err)
		return
	}
	res := dto.DeleteCourierResponse{
		Message: "courier's profile deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
