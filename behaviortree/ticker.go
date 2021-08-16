package behaviortree

import (
	"context"
	"errors"
	"sync"
	"time"
)

type (
	Ticker interface {
		Done() <-chan struct{}

		Err() error

		Stop()
	}

	tickerCore struct {
		ctx    context.Context
		cancel context.CancelFunc
		node   Node
		ticker *time.Ticker
		done   chan struct{}
		stop   chan struct{}
		once   sync.Once
		mutex  sync.Mutex
		err    error
	}

	tickerStopOnFailure struct {
		Ticker
	}
)

var (
	errExitOnFailure = errors.New("errExitOnFailure")
)

func NewTicker(ctx context.Context, duration time.Duration, node Node) Ticker {
	if ctx == nil {
		panic(errors.New("behaviortree.NewTicker nil context"))
	}
	if duration <= 0 {
		panic(errors.New("behaviortree.NewTicker duration <= 0"))
	}
	if node == nil {
		panic(errors.New("behaviortree.NewTicker nil node"))
	}
	result := &tickerCore{
		node:   node,
		ticker: time.NewTicker(duration),
		done:   make(chan struct{}),
		stop:   make(chan struct{}),
	}
	result.ctx, result.cancel = context.WithCancel(ctx)
	go result.run()
	return result
}

func NewTickerStopOnFailure(ctx context.Context, duration time.Duration, node Node) Ticker {
	if node == nil {
		panic(errors.New("behaviortree.NewTickerStopOnFailure nil node"))
	}
	return tickerStopOnFailure{
		Ticker: NewTicker(
			ctx,
			duration,
			func() (Tick, []Node) {
				tick, children := node()
				if tick == nil {
					return nil, children
				}
				return func(children []Node) (Status, error) {
					status, err := tick(children)
					if err == nil && status == Failure {
						err = errExitOnFailure
					}
					return status, err
				}, children
			},
		),
	}
}

func (t *tickerCore) run() {
	var err error
TickLoop:
	for err == nil {
		select {
		case <-t.ctx.Done():
			err = t.ctx.Err()
		case <-t.stop:
			break TickLoop
		case <-t.ticker.C:
			_, err = t.node.Tick()
		}
	}
	t.mutex.Lock()
	t.err = err
	t.mutex.Unlock()
	t.Stop()
	t.cancel()
	close(t.done)
}

func (t *tickerCore) Done() <-chan struct{} {
	return t.done
}
func (t *tickerCore) Err() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.err
}

func (t *tickerCore) Stop() {
	t.once.Do(func() {
		t.ticker.Stop()
		close(t.stop)
	})
}

func (t tickerStopOnFailure) Err() error {
	err := t.Ticker.Err()
	if err == errExitOnFailure {
		return nil
	}
	return err
}
