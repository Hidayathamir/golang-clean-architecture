package errkit

import (
	"fmt"
	"strings"
)

func Wrap(originalErr error, s string) error {
	return fmt.Errorf("%s:: %w", s, originalErr)
}

func WrapE(originalErr error, wrapperErr error) error {
	return fmt.Errorf("%w:: %w", wrapperErr, originalErr)
}

func Split(err error) []string {
	if err == nil {
		return nil
	}
	return strings.Split(err.Error(), ":: ")
}
