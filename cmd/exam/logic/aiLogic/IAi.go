package aiLogic

import "rxt/internal/core"

type Ai struct {
	IAi
}

func New(param ...*core.Core) Ai {
	res := &V1{}
	res.Init(param...)
	return Ai{res}
}

type IAi interface {
	Dina(param string) (string, error)
}
