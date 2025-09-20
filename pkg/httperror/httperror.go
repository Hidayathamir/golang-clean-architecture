package httperror

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	HTTPCode int
	ID       string
	Message  string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%d] [%s] %s", e.HTTPCode, e.ID, e.Message)
}

func InternalServerError() *HTTPError {
	return &HTTPError{
		HTTPCode: http.StatusInternalServerError,
		ID:       "InternalServerError",
		Message:  "internal server error",
	}
}

func BadRequest() *HTTPError {
	return &HTTPError{
		HTTPCode: http.StatusBadRequest,
		ID:       "BadRequest",
		Message:  "bad request",
	}
}

func Conflict() *HTTPError {
	return &HTTPError{
		HTTPCode: http.StatusConflict,
		ID:       "Conflict",
		Message:  "conflict",
	}
}

func NotFound() *HTTPError {
	return &HTTPError{
		HTTPCode: http.StatusNotFound,
		ID:       "Not Found",
		Message:  "not found",
	}
}

func Unauthorized() *HTTPError {
	return &HTTPError{
		HTTPCode: http.StatusUnauthorized,
		ID:       "Unauthorized",
		Message:  "unauthorized",
	}
}
