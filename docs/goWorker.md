# goWorker 协程池

Go 语言提供了轻量的 goroutine，但标准库没有线程/协程池组件。goWorker 通过一个带缓冲的通道实现并发上限控制，并配合 `WaitGroup` 管理任务生命周期。

## 安装

```bash
go get github.com/erdong01/kit/goWorker
```

## 快速开始

```go
package main

import (
	"fmt"

	"github.com/erdong01/kit/goWorker"
)

func main() {
	// 限制同时最多 2 个 goroutine 执行
	pool := goWorker.New(2)

	worker := pool.WaitGroup()
	for i := 0; i < 5; i++ {
		i := i
		worker.Go(func() {
			fmt.Println("job", i)
		})
	}
	worker.Wait()
}
```

## 设计与行为

- 并发上限：通过 `chan struct{}` 作为信号量，任务启动前先占位，执行结束后释放。
- 任务等待：每个 `Worker` 内部自带一个 `sync.WaitGroup`，`Wait()` 只等待该 `Worker` 提交的任务。
- panic 处理：任务内 panic 会被捕获并 `log.Fatal` 终止进程。

## API 说明

```go
// 创建协程池
func New(count ...int) *Pool

// 创建一个 Worker，用于提交任务并等待完成
func (p *Pool) WaitGroup() *Worker

// 调整 worker 数量（注意：仅修改计数，不会调整信号量容量）
func (p *Pool) Tune(n int)

// 提交任务
func (w *Worker) Go(f func())

// 等待当前 Worker 的所有任务完成
func (w *Worker) Wait()
```

## 多批次任务示例

`Worker` 可多次创建，彼此等待独立，但共享同一个池的并发上限。

```go
pool := goWorker.New(2)

batchA := pool.WaitGroup()
batchA.Go(func() { /* ... */ })
batchA.Go(func() { /* ... */ })

batchB := pool.WaitGroup()
batchB.Go(func() { /* ... */ })

batchA.Wait()
batchB.Wait()
```

## 注意事项

- `Tune` 只更新 `workerCount` 字段，不会重建 `sem` 容量，当前实现下对并发上限无实际影响。如需调整并发数，请创建新的 `Pool`。
- 任务内 panic 会导致进程退出，若需要更温和的错误处理，请在任务函数内部自行 `recover`。
