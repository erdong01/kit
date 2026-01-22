package goWorker

import (
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
	if len(count) > 0 {
		workerCount = count[0]
	} else {
		workerCount = 64
	}
	return &Pool{
		workerCount: workerCount,
		sem:         make(chan struct{}, workerCount),
	}
}

func (that *Pool) WaitGroup() *Worker {
	return &Worker{pool: that}
}

func (that *Pool) Tune(n int) {
	that.workerCount = n
}

func (that *Worker) Go(f func()) {
	that.pool.sem <- struct{}{}
	that.wg.Add(1)
	go func() {
		defer func() {
			<-that.pool.sem
			that.wg.Add(-1)
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
