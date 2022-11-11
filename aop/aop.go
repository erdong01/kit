package aop

import (
	"context"
)

type Aop struct {
	beforeSlice []IAop
	afterSlice  []IAop
	Ctx         context.Context
	I           IAop
}

type IAop interface {
	Before()
	Handler()
	After()
	Context(ctx context.Context)
	GetContext() context.Context
}

func New(Ctx context.Context, I IAop) *Aop {
	aop := &Aop{}
	aop.I = I
	aop.I.Context(Ctx)
	return aop
}

func (a *Aop) SetBefore(I ...IAop) *Aop {
	a.beforeSlice = I
	return a
}

func (a *Aop) SetAfter(I ...IAop) *Aop {
	a.afterSlice = I
	return a
}

func (a *Aop) Run() {
	a.I.Before()
	for _, f := range a.beforeSlice {
		f.Context(a.I.GetContext())
		f.Before()
		f.Handler()
		f.After()
	}
	a.I.Context(a.I.GetContext())
	a.I.Handler()
	a.I.After()
	for _, f := range a.afterSlice {
		f.Context(a.I.GetContext())
		f.Before()
		f.Handler()
		f.After()
	}
}

func (b *Aop) Before()  {}
func (b *Aop) Handler() {}
func (b *Aop) After()   {}

func (b *Aop) Context(ctx context.Context) {
	b.Ctx = ctx
}

func (b *Aop) GetContext() context.Context {
	return b.Ctx
}

func (b *Aop) Set(key, vale any) {
	b.Ctx = context.WithValue(b.Ctx, key, vale)
}

func (b *Aop) Get(key any) any {
	return b.Ctx.Value(key)
}
