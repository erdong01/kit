package goWorker

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	workerPool := New()
	workerPool.Go(func() {
		fmt.Println("000")
	})
	workerPool.Go(func() {
		fmt.Println("111")
	})
	workerPool.Go(func() {
		fmt.Println("222")
	})
	workerPool.Go(func() {
		fmt.Println("333")
	})

	workerPool.Go(func() {
		fmt.Println("444")
	})
	workerPool.Go(func() {
		fmt.Println("555")
	})
	workerPool.Go(func() {
		fmt.Println("666")
	})
	workerPool.Go(func() {
		fmt.Println("777")
	})
	workerPool.Go(func() {
		fmt.Println("888")
	})
	workerPool.Go(func() {
		fmt.Println("999")
	})
	time.Sleep(1 * time.Second)
}

func TestWork(t *testing.T) {
	workerPool := New(1)
	go func() {
		goWorker := workerPool.WaitGroup()
		goWorker.Go(func() {
			fmt.Println("111")
		})
		goWorker.Go(func() {
			fmt.Println("222")
		})

		goWorker.Go(func() {
			fmt.Println("333")
		})

		goWorker.Go(func() {
			fmt.Println("444")
		})

		goWorker.Go(func() {
			fmt.Println("555")
		})
		goWorker.Go(func() {
			fmt.Println("555")
		})
		goWorker.Go(func() {
			fmt.Println("777")
		})
		goWorker.Go(func() {
			fmt.Println("888")
		})
		goWorker.Go(func() {
			fmt.Println("999")
		})

		goWorker.Go(func() {
			fmt.Println("10001")
		})
		goWorker.Go(func() {
			fmt.Println("10002")
		})
		goWorker.Wait()
		fmt.Println("Wait0")
	}()
	go func() {
		goWorker := workerPool.WaitGroup()
		goWorker.Go(func() {
			fmt.Println("aaa111")
		})
		goWorker.Go(func() {
			fmt.Println("aaa222")
		})
		goWorker.Wait()
		fmt.Println("Wait-A")
	}()
	go func() {
		goWorker := workerPool.WaitGroup()
		goWorker.Go(func() {
			fmt.Println("bbb111")
		})
		goWorker.Go(func() {
			fmt.Println("bbb222")
		})
		goWorker.Wait()
		fmt.Println("Wait-B")
	}()
	time.Sleep(2 * time.Second)
}

func TestErrGroup1(t *testing.T) {
	workerPool := New(3)
	ctx := context.Background()
	errGroup, ctx := workerPool.WithContext(ctx)
	errGroup.Go(func() error {
		fmt.Println("222")
		return nil
	})
	errGroup.Go(func() error {
		fmt.Println("333")
		return errors.New("err333")
	})
	errGroup.Go(func() error {
		fmt.Println("444")
		return errors.New("err444")
	})
	if err := errGroup.Wait(); err != nil {
		fmt.Println("err:", err)
	}
	time.Sleep(1 * time.Second)
}

func TestErrGroup2(t *testing.T) {
	workerPool := New(3)
	ctx := context.Background()
	errGroup, ctx := workerPool.WithContext(ctx)
	errGroup.TryGo(func() error {
		fmt.Println("222")
		return nil
	})
	errGroup.TryGo(func() error {
		fmt.Println("333")
		return errors.New("err333")
	})
	errGroup.TryGo(func() error {
		fmt.Println("444")
		return errors.New("err444")
	})
	if err := errGroup.Wait(); err != nil {
		fmt.Println("err:", err)
	}
	time.Sleep(1 * time.Second)
}
