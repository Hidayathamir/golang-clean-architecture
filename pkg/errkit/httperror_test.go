package errkit

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetMessage(t *testing.T) {
	baseErr := errors.New("base")

	err := SetHTTPError(baseErr, http.StatusBadRequest)
	err = SetMessage(err, "custom message")

	httpErr := GetHTTPError(err)

	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPCode)
	assert.Equal(t, "custom message", httpErr.Message)
	assert.ErrorIs(t, err, baseErr)
}

func TestSetCode(t *testing.T) {
	baseErr := errors.New("base")

	err := SetHTTPError(baseErr, http.StatusBadRequest)
	err = SetMessage(err, "custom message")
	err = SetCode(err, http.StatusConflict)

	httpErr := GetHTTPError(err)

	assert.Equal(t, http.StatusConflict, httpErr.HTTPCode)
	assert.Equal(t, "custom message", httpErr.Message)
	assert.ErrorIs(t, err, baseErr)
}

func TestLoadErrAsHTTPError_1(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = Wrap(err, "wrap")
	err = Unauthorized(err)
	err = Wrap(err, "wrap2")

	httpErr := GetHTTPError(err)

	assert.Equal(t, "Unauthorized", httpErr.Message)
}

func TestLoadErrAsHTTPError_2(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = Wrap(err, "wrap")
	err = Wrap(err, "wrap2")

	httpErr := GetHTTPError(err)

	assert.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_3(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = Wrap(err, "wrap")
	err = Unauthorized(err)
	err = Wrap(err, "wrap2")
	err = InternalServerError(err)

	httpErr := GetHTTPError(err)

	assert.Equal(t, "Internal Server Error", httpErr.Message)
}

func TestLoadErrAsHTTPError_4(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = Wrap(err, "wrap")
	err = Unauthorized(err)
	err = Wrap(err, "wrap2")
	err = InternalServerError(err)
	err = Unauthorized(err)

	httpErr := GetHTTPError(err)

	assert.Equal(t, "Unauthorized", httpErr.Message)
}

func TestLoadErrAsHTTPError_5(t *testing.T) {
	var err error = &HTTPError{
		HTTPCode: http.StatusUnauthorized,
		Message:  "unauthorized",
	}

	httpErr := GetHTTPError(err)

	assert.Equal(t, "unauthorized", httpErr.Message)
}

func TestLoadErrAsHTTPError_6(t *testing.T) {
	var err error

	httpErr := GetHTTPError(err)

	assert.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_7(t *testing.T) {
	var err error
	err = &HTTPError{
		HTTPCode: http.StatusUnauthorized,
		Message:  "unauthorized",
	}
	err = Wrap(err, "wrap1")

	httpErr := GetHTTPError(err)

	assert.Equal(t, "unauthorized", httpErr.Message)
}
