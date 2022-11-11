package aop

import (
	"context"
)

type Aop struct {
	Err error
	Ctx context.Context
	end bool
}

type BaseAop struct {
	beforeSlice []IAop
	afterSlice  []IAop
	I           IAop
}

type IAop interface {
	Before()
	Handler()
	After()

	setContext(ctx context.Context)
	getContext() context.Context
	SetErr(err error)
	GetErr() (err error)
	setAop(aop *Aop)
	getAop() (aop *Aop)
	Break(err error)

	getEnd() (b bool)
}

func New(Ctx context.Context, I IAop) *BaseAop {
	aop := &BaseAop{}
	aop.I = I
	aop.I.setContext(Ctx)
	return aop
}

func (a *BaseAop) SetBefore(I ...IAop) *BaseAop {
	a.beforeSlice = I
	return a
}

func (a *BaseAop) SetAfter(I ...IAop) *BaseAop {
	a.afterSlice = I
	return a
}

func (a *BaseAop) Run() (err error) {
	a.I.Before()
	if a.I.getEnd() {
		return a.I.GetErr()
	}
	for _, f := range a.beforeSlice {
		f.setAop(a.I.getAop())
		f.Before()
		f.Handler()
		f.After()
		a.I.setAop(f.getAop())

		if f.getEnd() {
			return f.GetErr()
		}
	}
	a.I.Handler()
	a.I.After()
	if a.I.getEnd() {
		return a.I.GetErr()
	}
	for _, f := range a.afterSlice {
		f.setAop(a.I.getAop())
		f.Before()
		f.Handler()
		f.After()
		a.I.setAop(f.getAop())
		if f.getEnd() {
			return f.GetErr()
		}
	}
	return a.I.GetErr()
}

func (b *Aop) Before()  {}
func (b *Aop) Handler() {}
func (b *Aop) After()   {}

func (a *Aop) setContext(ctx context.Context) {
	a.Ctx = ctx
}

func (a *Aop) getContext() context.Context {
	return a.Ctx
}

func (a *Aop) Set(key, vale any) {
	a.Ctx = context.WithValue(a.Ctx, key, vale)
}

func (a *Aop) Get(key any) any {
	return a.Ctx.Value(key)
}

func (a *Aop) SetErr(err error) {
	a.Err = err
}
func (a *Aop) GetErr() (err error) {
	return a.Err
}

func (a *Aop) setAop(aop *Aop) {
	a.Err = aop.Err
	a.Ctx = aop.Ctx
	a.end = aop.end
}

func (a *Aop) getAop() (aop *Aop) {
	aop = a
	return
}

func (a *Aop) Break(err error) {
	a.end = true
	a.Err = err
}
func (a *Aop) getEnd() (b bool) {
	return a.end
}
