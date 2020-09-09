package wrong

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()

}

func Panic(code int, err error, message ...string) {

	var msg string
	msg = ""
	if len(message) > 0 {
		msg = message[0]
	}
	errStruct := New(code, err, msg)
	panic(errStruct)
}
