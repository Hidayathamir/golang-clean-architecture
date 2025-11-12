package errkit

import (
	"net/http"
)

func InternalServerError(err error) error {
	return SetHTTPError(err, http.StatusInternalServerError)
}

func BadRequest(err error) error {
	return SetHTTPError(err, http.StatusBadRequest)
}

func Conflict(err error) error {
	return SetHTTPError(err, http.StatusConflict)
}

func NotFound(err error) error {
	return SetHTTPError(err, http.StatusNotFound)
}

func Unauthorized(err error) error {
	return SetHTTPError(err, http.StatusUnauthorized)
}
