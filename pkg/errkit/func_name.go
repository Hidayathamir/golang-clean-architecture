package errkit

import (
	"golang-clean-architecture/pkg/caller"
)

type OptAddFuncName struct {
	Skip int
}

type OptionAddFuncName func(*OptAddFuncName)

const defaultSkip = 1

// AddFuncName will wrap error with caller function name. Use WithSkip to skip frame.
func AddFuncName(err error, options ...OptionAddFuncName) error {
	option := &OptAddFuncName{Skip: defaultSkip}
	for _, opt := range options {
		opt(option)
	}

	callerFuncName := caller.FuncName(caller.WithSkip(option.Skip))
	return Wrap(callerFuncName, err)
}

func WithSkip(skip int) OptionAddFuncName {
	return func(wo *OptAddFuncName) {
		wo.Skip = skip + defaultSkip
	}
}
