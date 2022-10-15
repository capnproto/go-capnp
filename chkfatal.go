package util

// Chkfatal panics if err is not nil
func Chkfatal(err error) {
	if err != nil {
		panic(err)
	}
}
