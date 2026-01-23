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

type NextFunc func() error

// Arounder provides AOP-style chaining. Implement to control when next runs.
type Arounder interface {
	Around(next NextFunc) error
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
		if f.getEnd() {
			a.I.setAop(f.getAop())
			return f.GetErr()
		}
		f.Handler()
		if f.getEnd() {
			a.I.setAop(f.getAop())
			return f.GetErr()
		}
		f.After()
		a.I.setAop(f.getAop())

		if f.getEnd() {
			return f.GetErr()
		}
	}
	a.I.Handler()
	if a.I.getEnd() {
		return a.I.GetErr()
	}
	a.I.After()
	if a.I.getEnd() {
		return a.I.GetErr()
	}
	for _, f := range a.afterSlice {
		f.setAop(a.I.getAop())
		f.Before()
		if f.getEnd() {
			a.I.setAop(f.getAop())
			return f.GetErr()
		}
		f.Handler()
		if f.getEnd() {
			a.I.setAop(f.getAop())
			return f.GetErr()
		}
		f.After()
		a.I.setAop(f.getAop())
		if f.getEnd() {
			return f.GetErr()
		}
	}
	return a.I.GetErr()
}

// RunAround executes a full AOP chain using Around if implemented.
// Order is: beforeSlice -> I -> afterSlice. Each link may call next() or short-circuit.
func (a *BaseAop) RunAround() (err error) {
	chain := make([]IAop, 0, len(a.beforeSlice)+1+len(a.afterSlice))
	chain = append(chain, a.beforeSlice...)
	chain = append(chain, a.I)
	chain = append(chain, a.afterSlice...)

	state := &Aop{}
	state.setAop(a.I.getAop())

	var runAt func(i int) error
	runAt = func(i int) error {
		if i >= len(chain) {
			return nil
		}
		if state.getEnd() {
			return state.GetErr()
		}

		cur := chain[i]
		cur.setAop(state)

		if ar, ok := cur.(Arounder); ok {
			next := func() error {
				state.setAop(cur.getAop())
				err := runAt(i + 1)
				cur.setAop(state)
				if err != nil {
					cur.SetErr(err)
				}
				return err
			}
			err := ar.Around(next)
			if err != nil && cur.GetErr() == nil {
				cur.SetErr(err)
			}
			state.setAop(cur.getAop())
			if cur.getEnd() {
				return cur.GetErr()
			}
			if err != nil {
				return err
			}
			return cur.GetErr()
		}

		cur.Before()
		if cur.getEnd() {
			state.setAop(cur.getAop())
			return cur.GetErr()
		}
		cur.Handler()
		if cur.getEnd() {
			state.setAop(cur.getAop())
			return cur.GetErr()
		}

		state.setAop(cur.getAop())
		err := runAt(i + 1)
		cur.setAop(state)
		if err != nil {
			cur.SetErr(err)
		}
		cur.After()
		if cur.getEnd() {
			state.setAop(cur.getAop())
			return cur.GetErr()
		}
		state.setAop(cur.getAop())
		if err != nil {
			return err
		}
		return cur.GetErr()
	}

	err = runAt(0)
	a.I.setAop(state)
	return err
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
