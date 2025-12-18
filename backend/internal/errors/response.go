package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents the JSON error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// WriteError writes an error response to the HTTP response writer
func WriteError(w http.ResponseWriter, err error) {
	var appErr *AppError
	var statusCode int
	var message string

	// Check if it's an AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
		statusCode = e.Code
		message = e.Message
	} else {
		// Default to internal server error for unknown errors
		statusCode = http.StatusInternalServerError
		message = "An internal server error occurred"
		log.Printf("Unhandled error: %v", err)
	}

	// Log the error if it's a server error (5xx)
	if statusCode >= 500 {
		if appErr != nil && appErr.Err != nil {
			log.Printf("Server error: %s - %v", message, appErr.Err)
		} else {
			log.Printf("Server error: %s", message)
		}
	}

	WriteJSON(w, statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	})
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// WriteSuccess writes a success response with data
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, SuccessResponse{Data: data})
}

// WriteSuccessWithMessage writes a success response with a message
func WriteSuccessWithMessage(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusOK, SuccessResponse{
		Data:    data,
		Message: message,
	})
}

// WriteCreated writes a 201 Created response
func WriteCreated(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusCreated, SuccessResponse{Data: data})
}

// WriteNoContent writes a 204 No Content response
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
