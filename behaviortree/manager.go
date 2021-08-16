package behaviortree

import (
	"sync"

	"github.com/joeycumines/go-bigbuff"
)

type (
	Manager interface {
		Ticker
		Add(ticker Ticker) error
	}
	manager struct {
		mu      sync.RWMutex
		once    sync.Once
		worker  bigbuff.Worker
		done    chan struct{}
		stop    chan struct{}
		tickers chan managerTicker
		errs    []error
	}

	managerTicker struct {
		Ticker Ticker
		Done   func()
	}

	errManagerTicker []error

	errManagerStopped struct{ error }
)
