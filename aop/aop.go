package aop

import "context"

type Aop struct {
	Context context.Context
	I       IAop
}

type IAop interface {
	Before()
	Handler()
	After()
}

func New(Ctx context.Context, I IAop) *Aop {
	aop := &Aop{}
	aop.Context = Ctx
	aop.I = I
	return aop
}

func (a *Aop) Run() {
	a.I.Before()
	a.I.Handler()
	a.I.After()
}

type Base struct{}

func (b *Base) Before()  {}
func (b *Base) Handler() {}
func (b *Base) After()   {}
