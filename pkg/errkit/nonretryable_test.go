package errkit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNonRetryableError(t *testing.T) {
	baseErr := errors.New("base error")
	nrErr := WrapNonRetryable(baseErr)

	require.True(t, IsNonRetryable(nrErr))
	require.Equal(t, "[NonRetryable]:: base error", nrErr.Error())

	var nre *NonRetryableError
	require.True(t, errors.As(nrErr, &nre))
	require.ErrorIs(t, nrErr, baseErr)
}

func TestNonRetryableErrorIsNotRetryable(t *testing.T) {
	baseErr := errors.New("base error")
	require.False(t, IsNonRetryable(baseErr))
}

func TestNonRetryableErrorWrapped(t *testing.T) {
	baseErr := errors.New("base error")
	nrErr := WrapNonRetryable(baseErr)
	wrappedErr := Wrap(nrErr, "wrapped")

	require.True(t, IsNonRetryable(wrappedErr))
	require.ErrorIs(t, wrappedErr, baseErr)
}

func TestNonRetryableNilErr(t *testing.T) {
	require.False(t, IsNonRetryable(nil))
}
