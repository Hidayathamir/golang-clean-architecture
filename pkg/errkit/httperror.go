package errkit

import (
	"fmt"
)

type HTTPError struct {
	HTTPCode int
	Message  string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", e.HTTPCode, e.Message)
}
