package http

import (
	"encoding/json"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respondErr(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, ErrResponse{Error: msg})
}
