package handlers

import (
	"encoding/json"
	"net/http"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
)

func PingGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		http.Error(w, server.ErrRequestCanceled, http.StatusRequestTimeout)
	default:

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := dto.PingResponse{Message: "pong"}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HealthcheckHeadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		http.Error(w, server.ErrRequestCanceled, http.StatusRequestTimeout)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
