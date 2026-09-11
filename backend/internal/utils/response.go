package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Issue string `json:"issue"`
}

type APIError struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Details   []ErrorDetail `json:"details,omitempty"`
	Timestamp string        `json:"timestamp,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, code, message string) {
	ErrorWithDetails(w, status, code, message, nil)
}

func ErrorWithDetails(w http.ResponseWriter, status int, code, message string, details []ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Code:      code,
			Message:   message,
			Details:   details,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

func ValidationError(w http.ResponseWriter, message string, details []ErrorDetail) {
	if message == "" {
		message = "Data validation failed"
	}
	ErrorWithDetails(w, http.StatusUnprocessableEntity, "INVALID_ARGUMENT", message, details)
}

func Unauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Authentication required"
	}
	Error(w, http.StatusUnauthorized, "UNAUTHENTICATED", message)
}

func Forbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "You do not have permission to access this resource."
	}
	Error(w, http.StatusForbidden, "PERMISSION_DENIED", message)
}

func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(w, http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource state conflict detected"
	}
	Error(w, http.StatusConflict, "RESOURCE_CONFLICT", message)
}

func RateLimitExceeded(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Too many requests. Please try again later."
	}
	Error(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", message)
}

func InternalServerError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "An unexpected error occurred"
	}
	Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}
