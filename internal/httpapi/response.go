package httpapi

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, errorResponse{Ok: false, Error: err.Error()})
}
