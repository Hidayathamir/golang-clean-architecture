package errkit

import (
	"fmt"
)

func Wrap(originalErr error, s string) error {
	return fmt.Errorf("%s:: %w", s, originalErr)
}

func WrapE(originalErr error, wrapperErr error) error {
	return fmt.Errorf("%w:: %w", wrapperErr, originalErr)
}
