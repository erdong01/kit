package base

import (
	"rxt/internal/core"
)

type Dao struct {
	*core.Core
	InitCore bool
	A        interface{}
}

func (s *Dao) Init(param ...*core.Core) *Dao {
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
