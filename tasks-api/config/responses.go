package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    any               `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// WriteResponse writes a JSON response with the given status code and response object
func WriteResponse(w http.ResponseWriter, status int, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// WriteSuccessResponse writes a success response
func WriteSuccessResponse(w http.ResponseWriter, message string, data any) {
	resp := NewSuccessResponse(message, data)
	WriteResponse(w, http.StatusOK, resp)
}

// WriteErrorResponse writes an error response
func WriteErrorResponse(w http.ResponseWriter, status int, message string, err error) {
	resp := NewErrorResponse(message, err)
	WriteResponse(w, status, resp)
}

// WriteCreatedResponse writes a response for successful creation
func WriteCreatedResponse(w http.ResponseWriter, message string, data any) {
	resp := NewSuccessResponse(message, data)
	WriteResponse(w, http.StatusCreated, resp)
}

// WriteNoContentResponse writes a 204 No Content response
func WriteNoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, validationErrors any) *Response {
	resp := &Response{Status: "error", Message: message}

	if validationErrors != nil {
		switch v := validationErrors.(type) {
		case validator.ValidationErrors:
			errMap := make(map[string]string)
			for _, err := range v {
				field := strings.ToLower(err.Field())
				var msg string
				switch err.Tag() {
				case "email":
					msg = fmt.Sprintf("%s must be a valid email address", field)
				case "required":
					msg = fmt.Sprintf("%s is required", field)
				case "min":
					msg = fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
				case "max":
					msg = fmt.Sprintf("%s must be at most %s characters long", field, err.Param())
				default:
					// For other validation tags, use validator's default message
					msg = fmt.Sprintf("%s %s", field, err.Tag())
					if err.Param() != "" {
						msg = fmt.Sprintf("%s %s %s", field, err.Tag(), err.Param())
					}
				}
				errMap[field] = msg
			}
			resp.Errors = errMap
		case error:
			resp.Message = v.Error()
		}
	}

	return resp
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data any) *Response {
	return &Response{Status: "success", Message: message, Data: data}
}
