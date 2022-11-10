package aop

import (
	"context"
)

type Aop struct {
	I           IAop
	beforeSlice []IAop
	afterSlice  []IAop
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

func (a *Aop) Before(I ...IAop) *Aop {
	a.beforeSlice = I
	return a
}

func (a *Aop) After(I ...IAop) *Aop {
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

type Base struct {
	Ctx context.Context
	End bool
}

func (b *Base) Before()  {}
func (b *Base) Handler() {}
func (b *Base) After()   {}
func (b *Base) Context(ctx context.Context) {
	b.Ctx = ctx
}
func (b *Base) GetContext() context.Context {
	return b.Ctx
}

func (b *Base) Set(key, vale any) {
	b.Ctx = context.WithValue(b.Ctx, key, vale)
}

func (b *Base) Get(key any) any {
	return b.Ctx.Value(key)
}
