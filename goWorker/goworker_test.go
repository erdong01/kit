package goWorker

import (
	"fmt"
	"testing"
	"time"
)

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
