package x

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfErrForDefer(f func() error) {
	err := f()
	if err != nil {
		panic(err)
	}
}
