package errkit

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetMessage(t *testing.T) {
	baseErr := errors.New("base")

	err := SetCode(baseErr, http.StatusBadRequest)
	err = SetMessage(err, "custom message")

	httpErr := GetHTTPError(err)

	require.Equal(t, http.StatusBadRequest, httpErr.HTTPCode)
	require.Equal(t, "custom message", httpErr.Message)
	require.ErrorIs(t, err, baseErr)
}

func TestSetCode(t *testing.T) {
	baseErr := errors.New("base")

	err := SetCode(baseErr, http.StatusBadRequest)
	err = SetMessage(err, "custom message")
	err = SetCode(err, http.StatusConflict)

	httpErr := GetHTTPError(err)

	require.Equal(t, http.StatusConflict, httpErr.HTTPCode)
	require.Equal(t, "custom message", httpErr.Message)
	require.ErrorIs(t, err, baseErr)
}

func TestLoadErrAsHTTPError_1(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = Wrap(err, "wrap")
	err = SetCode(err, http.StatusUnauthorized)
	err = Wrap(err, "wrap2")

	httpErr := GetHTTPError(err)

	require.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_2(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = Wrap(err, "wrap")
	err = Wrap(err, "wrap2")

	httpErr := GetHTTPError(err)

	require.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_3(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = Wrap(err, "wrap")
	err = SetCode(err, http.StatusUnauthorized)
	err = Wrap(err, "wrap2")
	err = SetCode(err, http.StatusInternalServerError)

	httpErr := GetHTTPError(err)

	require.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_4(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = Wrap(err, "wrap")
	err = SetCode(err, http.StatusUnauthorized)
	err = Wrap(err, "wrap2")
	err = SetCode(err, http.StatusInternalServerError)
	err = SetCode(err, http.StatusUnauthorized)

	httpErr := GetHTTPError(err)

	require.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_5(t *testing.T) {
	var err error = &HTTPError{
		HTTPCode: http.StatusUnauthorized,
		Message:  "unauthorized",
	}

	httpErr := GetHTTPError(err)

	require.Equal(t, "unauthorized", httpErr.Message)
}

func TestLoadErrAsHTTPError_6(t *testing.T) {
	var err error

	httpErr := GetHTTPError(err)

	require.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_7(t *testing.T) {
	var err error
	err = &HTTPError{
		HTTPCode: http.StatusUnauthorized,
		Message:  "unauthorized",
	}
	err = Wrap(err, "wrap1")

	httpErr := GetHTTPError(err)

	require.Equal(t, "unauthorized", httpErr.Message)
}
