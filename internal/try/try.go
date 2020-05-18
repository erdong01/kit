package try

import (
	"errors"
	"fmt"
)

type RxtError struct {
	err        error
	error_code int
}

func main() {

	var rxtError RxtError
	rxtError.err = errors.New("搞错了")
	rxtError.error_code = 11

	Try(func() {
		panic(rxtError)
	}, func(err interface{}) {
		errs := (err).(RxtError)
		fmt.Println(errs.err)
	})

	Try(func() {
		test()
	}, func(err interface{}) {
		errs := (err).(RxtError)
		fmt.Println(errs.err)
	})
}

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func test() {
	var rxtError RxtError
	rxtError.err = errors.New("test")
	rxtError.error_code = 11
	panic(rxtError)
}
