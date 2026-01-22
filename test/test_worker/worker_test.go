package testworker

import (
	"fmt"
	"testing"
	"time"

	"github.com/erdong01/kit/goWorker"
)

func TestWork(t *testing.T) {
	workerPool := goWorker.New(2)
	go func() {
		worker := workerPool.WaitGroup()
		worker.Go(func() {
			fmt.Println("111")
		})
		worker.Go(func() {
			fmt.Println("222")
		})

		worker.Go(func() {
			fmt.Println("333")
		})

		worker.Go(func() {
			fmt.Println("444")
		})

		worker.Go(func() {
			fmt.Println("555")
		})
		worker.Go(func() {
			fmt.Println("555")
		})
		worker.Go(func() {
			fmt.Println("777")
		})
		worker.Go(func() {
			fmt.Println("888")
		})
		worker.Go(func() {
			fmt.Println("999")
		})
		worker.Wait()
		fmt.Println("Wait0")
	}()
	go func() {
		worker := workerPool.WaitGroup()
		worker.Go(func() {
			fmt.Println("aaa111")
		})
		worker.Go(func() {
			fmt.Println("aaa222")
		})
		worker.Wait()
		fmt.Println("Wait-A")
	}()
	go func() {
		worker := workerPool.WaitGroup()
		worker.Go(func() {
			fmt.Println("bbb111")
		})
		worker.Go(func() {
			fmt.Println("bbb222")
		})
		worker.Wait()
		fmt.Println("Wait-B")
	}()
	time.Sleep(1 * time.Second)
}
