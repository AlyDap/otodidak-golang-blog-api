// utils/error.go
package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendError(w http.ResponseWriter, msg string, statusCode int, err error) {
	if err != nil {
		log.Printf("Internal error: %v", err) // log error asli di server
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
