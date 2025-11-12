package errkit

import "github.com/Hidayathamir/golang-clean-architecture/pkg/caller"

// AddFuncName wraps the error using the caller's function name automatically.
func AddFuncName(err error) error {
	return Wrap(err, caller.FuncName(caller.WithSkip(2)))
}
