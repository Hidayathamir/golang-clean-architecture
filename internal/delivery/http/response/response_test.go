package response_test

import (
	"errors"
	"fmt"
	"golang-clean-architecture/internal/delivery/http/response"
	"golang-clean-architecture/pkg/httperror"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadErrAsHTTPError_1(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = fmt.Errorf("wrap: %w", err)
	err = errors.Join(httperror.Unauthorized(), err)
	err = fmt.Errorf("wrap2: %w", err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, httperror.Unauthorized().ID, httpErr.ID)
	assert.Equal(t, httperror.Unauthorized().Message, httpErr.Message)
}

func TestLoadErrAsHTTPError_2(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = fmt.Errorf("wrap: %w", err)
	err = fmt.Errorf("wrap2: %w", err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, httperror.InternalServerError().ID, httpErr.ID)
	assert.Equal(t, httperror.InternalServerError().Message, httpErr.Message)
}

func TestLoadErrAsHTTPError_3(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = fmt.Errorf("wrap: %w", err)
	err = errors.Join(httperror.Unauthorized(), err)
	err = fmt.Errorf("wrap2: %w", err)
	err = errors.Join(httperror.InternalServerError(), err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, httperror.InternalServerError().ID, httpErr.ID)
	assert.Equal(t, httperror.InternalServerError().Message, httpErr.Message)
}

func TestLoadErrAsHTTPError_4(t *testing.T) {
	var err error
	err = errors.New("dummy err 1")
	err = fmt.Errorf("wrap: %w", err)
	err = errors.Join(httperror.Unauthorized(), err)
	err = fmt.Errorf("wrap2: %w", err)
	err = errors.Join(httperror.InternalServerError(), err)
	err = errors.Join(httperror.Unauthorized(), err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, httperror.Unauthorized().ID, httpErr.ID)
	assert.Equal(t, httperror.Unauthorized().Message, httpErr.Message)
}

func TestLoadErrAsHTTPError_5(t *testing.T) {
	var err error = httperror.Unauthorized()

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, httperror.Unauthorized().ID, httpErr.ID)
	assert.Equal(t, httperror.Unauthorized().Message, httpErr.Message)
}

func TestLoadErrAsHTTPError_6(t *testing.T) {
	var err error

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, httperror.InternalServerError().ID, httpErr.ID)
	assert.Equal(t, httperror.InternalServerError().Message, httpErr.Message)
}

func TestLoadErrAsHTTPError_7(t *testing.T) {
	var err error
	err = httperror.Unauthorized()
	err = fmt.Errorf("wrap1: %w", err)

	httpErr := response.LoadErrAsHTTPError(err)

	assert.Equal(t, httperror.Unauthorized().ID, httpErr.ID)
	assert.Equal(t, httperror.Unauthorized().Message, httpErr.Message)
}
