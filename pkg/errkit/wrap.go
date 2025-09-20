package errkit

import (
	"fmt"
	"strings"
)

func Wrap(s string, err error) error {
	return fmt.Errorf("%s:: %w", s, err)
}

func WrapE(err1 error, err2 error) error {
	return fmt.Errorf("%w:: %w", err1, err2)
}

func Split(err error) []string {
	if err == nil {
		return nil
	}
	return strings.Split(err.Error(), ":: ")
}
