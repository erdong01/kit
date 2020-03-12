package moselService

import "rxt/internal/core"

type container struct {
	I
}

func New(param ...*core.Core) (c *container) {
	exam := &V1{}
	exam.Init(param...)
	return &container{exam}
}

type I interface {
}
