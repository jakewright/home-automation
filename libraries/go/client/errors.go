package client

import (
	"net/http"
)

// ErrorType encapsulates a message and status code for error responses
type ErrorType struct {
	Message string `json:"title"`  // Message of the error
	Status  int    `json:"status"` // The HTTP status code to send in the response
}

var (
	// ErrInternalService is a generic error
	ErrInternalService = ErrorType{Status: http.StatusInternalServerError, Message: "General error"}

	// ErrDecodingJSON is used if the request contained invalid JSON
	ErrDecodingJSON = ErrorType{Status: http.StatusBadRequest, Message: "Could not decode JSON"}

	// ErrBadParam is used if the request contained a malformed or missing parameter
	ErrBadParam = ErrorType{Status: http.StatusBadRequest, Message: "Invalid parameter"}

	// ErrResourceNotFound is used as a 404 error
	ErrResourceNotFound = ErrorType{Status: http.StatusNotFound, Message: "Resource not found"}
)
