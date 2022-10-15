package util

func Chkfatal(err error) {
	if err != nil {
		panic(err)
	}
}
