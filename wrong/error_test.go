package wrong

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	pkgErr := errors.New("testaaaaaa")
	errS := New(http.StatusOK, pkgErr, "")

	s := errS.Format()
	fmt.Println(s)
}

func errs() error {
	return errors.New("whoops")
}
func TestTry(t *testing.T) {
	Try(func() {
		//panic(New(http.StatusExpectationFailed, errors.New(http.StatusText(http.StatusExpectationFailed)), "9999"))
		Panic(http.StatusExpectationFailed, errors.New(http.StatusText(http.StatusExpectationFailed)))
	}, func(i interface{}) {
		err := i.(*Err)
		fmt.Println(err.Errord)
	})
}
