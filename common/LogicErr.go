package common

import "github.com/erDong01/micro-kit/base"

func DBERROR(msg string, err error) {
	base.LOG.Printf("db [%s] error [%s]", msg, err.Error())
}
