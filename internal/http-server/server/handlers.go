package router

import (
	"encoding/json"
	"net/http"
)

func pingGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	reportOk := "{ \"message\": \"pong\" }"
	enc := json.NewEncoder(w)
	if err := enc.Encode(reportOk); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
