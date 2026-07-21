package errkit

import "errors"

type NonRetryableError struct {
	Err error
}

func (e *NonRetryableError) Error() string {
	return e.Err.Error()
}

func (e *NonRetryableError) Unwrap() error {
	return e.Err
}

func WrapNonRetryable(err error) error {
	return Wrap(&NonRetryableError{Err: err}, "[NonRetryable]")
}

func IsNonRetryable(err error) bool {
	var nre *NonRetryableError
	return errors.As(err, &nre)
}
