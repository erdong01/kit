package core

import (
	"errors"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

type sig struct{}

type f func() error

type Pool struct {
	capacity int32

	running int32

	expiryDuration time.Duration

	workers []*Worker

	release chan sig

	lock sync.Mutex

	once sync.Once
}

const (
	// DefaultAntsPoolSize is the default capacity for a default goroutine pool.
	DefaultAntsPoolSize = math.MaxInt32

	// DefaultCleanIntervalTime is the interval time to clean up goroutines.
	DefaultCleanIntervalTime = time.Second
)

var (
	// Error types for the Ants API.
	//---------------------------------------------------------------------------

	// ErrInvalidPoolSize will be returned when setting a negative number as pool capacity.
	ErrInvalidPoolSize = errors.New("invalid size for pool")

	// ErrLackPoolFunc will be returned when invokers don't provide function for pool.
	ErrLackPoolFunc = errors.New("must provide function for pool")

	// ErrInvalidPoolExpiry will be returned when setting a negative number as the periodic duration to purge goroutines.
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")

	// ErrPoolClosed will be returned when submitting task to a closed pool.
	ErrPoolClosed = errors.New("this pool has been closed")

	// ErrPoolOverload will be returned when the pool is full and no workers available.
	ErrPoolOverload = errors.New("too many goroutines blocked on submit or Nonblocking is set")
	//---------------------------------------------------------------------------

	// workerChanCap determines whether the channel of a worker should be a buffered channel
	// to get the best performance. Inspired by fasthttp at https://github.com/valyala/fasthttp/blob/master/workerpool.go#L139
	workerChanCap = func() int {
		// Use blocking workerChan if GOMAXPROCS=1.
		// This immediately switches Serve to WorkerFunc, which results
		// in higher performance (under go1.5 at least).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// Use non-blocking workerChan if GOMAXPROCS>1,
		// since otherwise the Serve caller (Acceptor) may lag accepting
		// new connections if WorkerFunc is CPU-bound.
		return 1
	}()

	defaultLogger = Logger(log.New(os.Stderr, "", log.LstdFlags))

	// Init a instance pool when importing ants.
	defaultAntsPool, _ = NewPool(DefaultAntsPoolSize)
)

// Logger is used for logging formatted messages.
type Logger interface {
	// Printf must have the same semantics as log.Printf.
	Printf(format string, args ...interface{})
}
type Worker struct {
	pool        *Pool
	task        chan f
	recycleTime time.Time
}

func NewPool(size int)(*Pool,error) {
	return NewTimingPool(size, DefaultCleanIntervalTime)
}

func NewTimingPool(size int, expiry  time.Duration) (*Pool, error) {
	if size == 0 {
		return nil, ErrInvalidPoolSize
	}

	if expiry <= 0 {
		return nil, ErrInvalidPoolExpiry
	}
	p := &Pool{
		capacity:       int32(size),
		release:        make(chan sig, 1),
		expiryDuration: time.Duration(expiry) * time.Second,
	}
	return p, nil
}

func (p *Pool)Submit(task f)error{

}
