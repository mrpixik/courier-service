package server

import (
	"encoding/json"
	"net/http"
)

const ErrRequestCanceled = "request canceled"

// Решил создать именно структуру для ответа, а не декодировать JSON из map,
// чтобы в будущем можно было удобнее добавить какую-то доп информацию для возврата по ручке /ping
type pingHandler struct {
	Message string `json:"message"`
}

// Так же решил сразу добавить проверку отмены контекста, чтобы было
func pingGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		http.Error(w, ErrRequestCanceled, http.StatusRequestTimeout)
	default:

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := pingHandler{Message: "pong"}

		enc := json.NewEncoder(w)
		if err := enc.Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func healthcheckHeadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		http.Error(w, ErrRequestCanceled, http.StatusRequestTimeout)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
