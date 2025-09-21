package caller

import (
	"runtime"
	"strings"
)

type OptFuncName struct {
	Skip int
}

type OptionFuncName func(*OptFuncName)

const defaultSkip = 1

// FuncName will return caller function name. Use WithSkip to skip frame.
func FuncName(options ...OptionFuncName) string {
	option := &OptFuncName{Skip: defaultSkip}
	for _, opt := range options {
		opt(option)
	}

	pc, _, _, ok := runtime.Caller(option.Skip)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?"
	}

	funcNameWithModule := fn.Name()
	funcNameWithModuleSplit := strings.Split(funcNameWithModule, "/")
	funcName := funcNameWithModuleSplit[len(funcNameWithModuleSplit)-1]

	return funcName
}

func WithSkip(skip int) OptionFuncName {
	return func(o *OptFuncName) {
		o.Skip = skip + defaultSkip
	}
}
