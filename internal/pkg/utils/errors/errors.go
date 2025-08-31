package errors

import (
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Type:       "BAD_REQUEST",
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Type:       "UNAUTHORIZED",
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		Message:    message,
		Type:       "FORBIDDEN",
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    message,
		Type:       "NOT_FOUND",
	}
}

func Internal(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Type:       "INTERNAL_SERVER_ERROR",
	}
}