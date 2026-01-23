package goWorker

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

type Pool struct {
	workerCount int
	sem         chan struct{}
}

type Worker struct {
	pool *Pool
	wg   sync.WaitGroup
}

func New(count ...int) *Pool {
	var workerCount int
	var sem chan struct{}
	if len(count) > 0 {
		workerCount = count[0]
		sem = make(chan struct{}, workerCount)
	} else {
		sem = make(chan struct{}, 1)
	}
	return &Pool{
		workerCount: workerCount,
		sem:         sem,
	}
}

func (that *Pool) Go(f func()) {
	if f == nil {
		return
	}
	that.sem <- struct{}{}
	go func() {
		defer func() {
			<-that.sem
			if r := recover(); r != nil {
				msg := fmt.Sprintf("GOWorker: panic %s", debug.Stack())
				log.Fatal(msg)
			}
		}()
		f()
	}()
}

func (that *Pool) WaitGroup() *Worker {
	return &Worker{pool: that}
}

func (that *Worker) Go(f func()) {
	if f == nil {
		return
	}
	that.pool.sem <- struct{}{}
	that.wg.Add(1)
	go func() {
		defer func() {
			<-that.pool.sem
			that.wg.Done()
			if r := recover(); r != nil {
				msg := fmt.Sprintf("GOWorker: panic %s", debug.Stack())
				log.Fatal(msg)
			}
		}()
		f()
	}()
}

func (that *Worker) Wait() {
	that.wg.Wait()
}

type Group struct {
	pool   *Pool
	cancel func(error)

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

func (g *Group) done() {
	if g.pool.sem != nil {
		<-g.pool.sem
	}
	g.wg.Done()
}

func (that *Pool) WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{
		pool:   that,
		cancel: cancel,
	}, ctx
}

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel(g.err)
	}
	return g.err
}

func (g *Group) Go(f func() error) {
	if g.pool.sem != nil {
		g.pool.sem <- struct{}{}
	}

	g.wg.Add(1)
	go func() {
		defer g.done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel(g.err)
				}
			})
		}
	}()
}

func (g *Group) TryGo(f func() error) bool {
	if g.pool.sem != nil {
		select {
		case g.pool.sem <- struct{}{}:
		default:
			return false
		}
	}

	g.wg.Add(1)
	go func() {
		defer g.done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel(g.err)
				}
			})
		}
	}()
	return true
}
