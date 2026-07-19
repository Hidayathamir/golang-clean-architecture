package errkit

func AddFuncName(err error, funcName ...string) error {
	if len(funcName) > 0 {
		return Wrap(err, funcName[0])
	}
	return Wrap(err, "no func name provided")
}
