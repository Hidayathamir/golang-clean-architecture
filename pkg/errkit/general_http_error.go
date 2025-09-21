package errkit

import (
	"net/http"
)

func Custom(err error, status int, msg string) error {
	return WrapE(
		&HTTPError{
			HTTPCode: status,
			Message:  msg,
		},
		err,
	)
}

func InternalServerError(err error) error {
	return WrapE(
		&HTTPError{
			HTTPCode: http.StatusInternalServerError,
			Message:  "internal server error",
		},
		err,
	)
}

func BadRequest(err error) error {
	return WrapE(
		&HTTPError{
			HTTPCode: http.StatusBadRequest,
			Message:  "bad request",
		},
		err,
	)
}

func Conflict(err error) error {
	return WrapE(
		&HTTPError{
			HTTPCode: http.StatusConflict,
			Message:  "conflict",
		},
		err,
	)
}

func NotFound(err error) error {
	return WrapE(
		&HTTPError{
			HTTPCode: http.StatusNotFound,
			Message:  "not found",
		},
		err,
	)
}

func Unauthorized(err error) error {
	return WrapE(
		&HTTPError{
			HTTPCode: http.StatusUnauthorized,
			Message:  "unauthorized",
		},
		err,
	)
}
