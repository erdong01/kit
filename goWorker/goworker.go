package goWorker

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

type worker struct {
	// Record the number of running workers
	//协程数量
	workerCount int
	sem         chan struct{}
	wg          sync.WaitGroup
}

func NewWorker(count ...int) (p *worker) {
	p = &worker{}
	p.workerCount = 32
	if len(count) > 0 {
		p.workerCount = count[0]
	}
	p.sem = make(chan struct{}, p.workerCount)
	return p
}

func (that *worker) Go(f func()) {
	that.sem <- struct{}{}
	that.wg.Add(1)
	go func() {
		defer that.wg.Add(-1)
		defer func() { <-that.sem }()
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("GOWorker: panic %s", debug.Stack())
				log.Fatal(msg)
			}
		}()
		f()
	}()
}

func (that *worker) Wait() {
	that.wg.Wait()
}
