package service

import "github.com/erDong01/micro-kit/core"

type Service struct {
	*core.Core
	InitCore bool
	A        interface{}
}

func (s *Service) Init(param ...*core.Core) *Service {
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
