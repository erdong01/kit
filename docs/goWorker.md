## go 协程池
C++/Java基本都有专门的库去实现线程池，而go语言由于有goroutine的存在，并没有提供协程池这样的组件。

**自己写了一个控制goroutine数量组件 使用方法很简单**

```go
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

```

**go worker实现逻辑**

```go

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

```