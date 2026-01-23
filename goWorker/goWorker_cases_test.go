package goWorker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPoolConcurrencyLimit(t *testing.T) {
	pool := New(2)

	var current int32
	var max int32

	start := make(chan struct{})
	var wg sync.WaitGroup

	total := 4
	wg.Add(total)
	for i := 0; i < total; i++ {
		go func() {
			pool.Go(func() {
				<-start
				c := atomic.AddInt32(&current, 1)
				for {
					old := atomic.LoadInt32(&max)
					if c <= old || atomic.CompareAndSwapInt32(&max, old, c) {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				atomic.AddInt32(&current, -1)
				wg.Done()
			})
		}()
	}

	close(start)
	wg.Wait()

	if max > 2 {
		t.Fatalf("max concurrency %d > 2", max)
	}
}

func TestWorkerWaitBlocksUntilDone(t *testing.T) {
	pool := New(3)
	worker := pool.WaitGroup()

	start := make(chan struct{})
	var count int32

	tasks := 5
	for i := 0; i < tasks; i++ {
		worker.Go(func() {
			<-start
			atomic.AddInt32(&count, 1)
		})
	}

	waitDone := make(chan struct{})
	go func() {
		worker.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		t.Fatal("Wait returned before tasks started")
	case <-time.After(50 * time.Millisecond):
	}

	close(start)
	select {
	case <-waitDone:
	case <-time.After(1 * time.Second):
		t.Fatal("Wait did not return in time")
	}

	if got := atomic.LoadInt32(&count); got != int32(tasks) {
		t.Fatalf("expected %d tasks, got %d", tasks, got)
	}
}

func TestGroupWithContextCancel(t *testing.T) {
	pool := New(2)
	ctx := context.Background()
	group, ctx := pool.WithContext(ctx)

	cancelled := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(cancelled)
	}()

	group.Go(func() error {
		return errors.New("boom")
	})

	if err := group.Wait(); err == nil {
		t.Fatal("expected error, got nil")
	}

	select {
	case <-cancelled:
	case <-time.After(1 * time.Second):
		t.Fatal("context not cancelled")
	}
}

func TestGroupTryGoLimit(t *testing.T) {
	pool := New(1)
	group, _ := pool.WithContext(context.Background())

	block := make(chan struct{})
	started := make(chan struct{})

	if ok := group.TryGo(func() error {
		close(started)
		<-block
		return nil
	}); !ok {
		t.Fatal("expected first TryGo to succeed")
	}

	<-started
	if ok := group.TryGo(func() error { return nil }); ok {
		t.Fatal("expected second TryGo to fail when at limit")
	}

	close(block)
	if err := group.Wait(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
