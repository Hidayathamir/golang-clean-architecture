package response_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"

	"github.com/stretchr/testify/assert"
)

func TestLoadErrAsHTTPError_1(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = errkit.Wrap("wrap:: %w", err)
	err = errkit.Unauthorized(err)
	err = errkit.Wrap("wrap2:: %w", err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, "unauthorized", httpErr.Message)
}

func TestLoadErrAsHTTPError_2(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = errkit.Wrap("wrap:: %w", err)
	err = errkit.Wrap("wrap2:: %w", err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_3(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = errkit.Wrap("wrap:: %w", err)
	err = errkit.Unauthorized(err)
	err = errkit.Wrap("wrap2:: %w", err)
	err = errkit.InternalServerError(err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_4(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = errkit.Wrap("wrap:: %w", err)
	err = errkit.Unauthorized(err)
	err = errkit.Wrap("wrap2:: %w", err)
	err = errkit.InternalServerError(err)
	err = errkit.Unauthorized(err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, "unauthorized", httpErr.Message)
}

func TestLoadErrAsHTTPError_5(t *testing.T) {
	var err error = &errkit.HTTPError{
		HTTPCode: http.StatusUnauthorized,
		Message:  "unauthorized",
	}

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, "unauthorized", httpErr.Message)
}

func TestLoadErrAsHTTPError_6(t *testing.T) {
	var err error

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, "internal server error", httpErr.Message)
}

func TestLoadErrAsHTTPError_7(t *testing.T) {
	var err error
	err = &errkit.HTTPError{
		HTTPCode: http.StatusUnauthorized,
		Message:  "unauthorized",
	}
	err = errkit.Wrap("wrap1:: %w", err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, "unauthorized", httpErr.Message)
}
