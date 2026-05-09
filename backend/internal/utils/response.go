package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code"`
}

// WriteJSON writes JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteSuccess writes success response
func WriteSuccess(w http.ResponseWriter, data interface{}) error {
	return WriteJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Code:    0,
	})
}

// WriteError writes error response
func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
		Code:    status,
	})
}
