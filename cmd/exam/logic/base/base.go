package base

import (
	"rxt/internal/core"
)

type Logic struct {
	*core.Core
	InitCore bool
	A        interface{}
}

func (s *Logic) Init(param ...*core.Core) *Logic {
	if s.InitCore == false {
		if param != nil {
			s.Core = param[0]
		} else {
			s.Core = core.New()
		}
		s.InitCore = true
	}
	return s
}
