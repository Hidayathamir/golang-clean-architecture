package errkit

import (
	"fmt"
	"net/http"
)

func InternalServerError(err error) error {
	return fmt.Errorf(
		"%w: %w",
		&HTTPError{
			HTTPCode: http.StatusInternalServerError,
			Message:  "internal server error",
		},
		err,
	)
}

func BadRequest(err error) error {
	return fmt.Errorf(
		"%w: %w",
		&HTTPError{
			HTTPCode: http.StatusBadRequest,
			Message:  "bad request",
		},
		err,
	)
}

func Conflict(err error) error {
	return fmt.Errorf(
		"%w: %w",
		&HTTPError{
			HTTPCode: http.StatusConflict,
			Message:  "conflict",
		},
		err,
	)
}

func NotFound(err error) error {
	return fmt.Errorf(
		"%w: %w",
		&HTTPError{
			HTTPCode: http.StatusNotFound,
			Message:  "not found",
		},
		err,
	)
}

func Unauthorized(err error) error {
	return fmt.Errorf(
		"%w: %w",
		&HTTPError{
			HTTPCode: http.StatusUnauthorized,
			Message:  "unauthorized",
		},
		err,
	)
}
