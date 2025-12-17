package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"service-order-avito/internal/adapters"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
)

// mockgen -source="internal/http/server/handler/delivery/delivery.go" -destination="internal/http/server/handler/delivery/mocks/mock_delivery_service.go"
type deliveryService interface {
	Assign(context.Context, *dto.AssignDeliveryRequest) (*dto.AssignDeliveryResponse, error)
	Unassign(context.Context, *dto.UnassignDeliveryRequest) (*dto.UnassignDeliveryResponse, error)
}

type deliveryHandler struct {
	service deliveryService
}

func NewDeliveryHandler(service deliveryService) *deliveryHandler {
	return &deliveryHandler{service: service}
}

func (dh *deliveryHandler) PostAssign(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		adapters.WriteError(w, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	res, err := dh.service.Assign(r.Context(), &req)
	if err != nil {
		adapters.WriteServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

func (dh *deliveryHandler) PostUnassign(w http.ResponseWriter, r *http.Request) {
	var req dto.UnassignDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		adapters.WriteError(w, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	res, err := dh.service.Unassign(r.Context(), &req)
	if err != nil {
		adapters.WriteServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
