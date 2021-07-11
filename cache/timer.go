package cache

import (
	"sync/atomic"
	"time"
)

// Timer holds representation of current time.
type Timer interface {
	// Give current time (in seconds)
	Now() uint32
}

// Helper function that returns Unix time in seconds
func getUnixTime() uint32 {
	return uint32(time.Now().Unix())
}

type cacheTimer struct {
	now    uint32
	ticker *time.Ticker
	done   chan bool
}

func (this *cacheTimer) Now() uint32 {
	return atomic.LoadUint32(&this.now)
}

// Stop runtime timer and finish routine that updates time
func (this *cacheTimer) Stop() {
	this.ticker.Stop()
	this.done <- true
	close(this.done)
	this.done = nil
	this.ticker = nil
}

// Periodically check and update  of time
func (this *cacheTimer) update() {
	for {
		select {
		case <-this.done:
			return
		case <-this.ticker.C:
			atomic.StoreUint32(&this.now, getUnixTime())
		}
	}
}
