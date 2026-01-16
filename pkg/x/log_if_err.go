package x

func LogIfErr(err error) {
	if err != nil {
		Logger.Warn(err)
	}
}
