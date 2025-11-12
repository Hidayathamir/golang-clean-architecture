package errkit

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	HTTPCode int
	Message  string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", e.HTTPCode, e.Message)
}

func GetHTTPError(err error) *HTTPError {
	httpErr := &HTTPError{
		HTTPCode: http.StatusInternalServerError,
		Message:  "internal server error",
	}

	errors.As(err, &httpErr)

	return httpErr
}

func SetHTTPError(err error, code int) error {
	return WrapE(
		err,
		&HTTPError{
			HTTPCode: code,
			Message:  http.StatusText(code),
		},
	)
}

func SetMessage(err error, message string) error {
	httpErr := GetHTTPError(err)

	return WrapE(
		err,
		&HTTPError{
			HTTPCode: httpErr.HTTPCode,
			Message:  message,
		},
	)
}

func SetCode(err error, code int) error {
	httpErr := GetHTTPError(err)

	return WrapE(
		err,
		&HTTPError{
			HTTPCode: code,
			Message:  httpErr.Message,
		},
	)
}
