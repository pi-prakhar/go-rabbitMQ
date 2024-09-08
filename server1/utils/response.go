package utils

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	Write(w http.ResponseWriter) error
}

type SuccessResponse[T any] struct {
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
	Data       T      `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
	Error      string `json:"error"`
}

func (r *SuccessResponse[T]) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		return err
	}

	return nil
}

func (r *ErrorResponse) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		return err
	}

	return nil

}
