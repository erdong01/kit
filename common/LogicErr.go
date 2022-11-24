package common

import "github.com/erdong01/kit/base"

func DBERROR(msg string, err error) {
	base.LOG.Printf("db [%s] error [%s]", msg, err.Error())
}
