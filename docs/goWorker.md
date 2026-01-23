# goWorker go协程池

Go 语言提供了轻量的 goroutine，但标准库没有线程/协程池组件。goWorker 通过一个带缓冲的通道实现并发上限控制，并提供三种使用方式：

- `Pool.Go`：直接提交任务（不带等待）
- `Worker`：分批提交任务并等待完成
- `Group`：带错误与取消语义的并发执行（类似 `errgroup`）

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
	// 限制同时最多 10 个 goroutine 执行
	pool := goWorker.New(10)

	worker := pool.WaitGroup()
	for i := 0; i < 50; i++ {
		i := i
		worker.Go(func() {
			fmt.Println("job", i)
		})
	}
	worker.Wait()
}
```

## 使用方式

- 直接任务：`pool.Go(func(){ ... })` 只负责并发上限，不提供等待。
- 分批等待：`worker := pool.WaitGroup()` 后用 `worker.Go` 提交并 `worker.Wait()` 等待。
- 错误与取消：`group, ctx := pool.WithContext(ctx)` 支持首错返回与 `ctx` 取消。

## 设计与行为

- 并发上限：通过 `chan struct{}` 作为信号量，任务启动前先占位，执行结束后释放。
- 任务等待：每个 `Worker` 内部自带一个 `sync.WaitGroup`，`Wait()` 只等待该 `Worker` 提交的任务。
- 错误与取消：`Group` 记录首个错误并取消上下文，其余任务仍会继续执行直到结束。
- panic 处理：任务内 panic 会被捕获并 `log.Fatal` 终止进程。

## API 说明

```go
// 创建协程池；不传参时默认并发为 1
func New(count ...int) *Pool

// 创建一个 Worker，用于提交任务并等待完成
func (p *Pool) WaitGroup() *Worker

// 提交任务
func (p *Pool) Go(f func())

// 提交任务
func (w *Worker) Go(f func())

// 等待当前 Worker 的所有任务完成
func (w *Worker) Wait()

// 创建一个带取消能力的任务组
func (p *Pool) WithContext(ctx context.Context) (*Group, context.Context)

// 提交一个返回 error 的任务
func (g *Group) Go(f func() error)

// 尝试提交任务，若达到并发上限立即返回 false
func (g *Group) TryGo(f func() error) bool

// 等待所有任务结束并返回首个错误
func (g *Group) Wait() error
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

## 带错误与取消的并发示例

`WithContext` 返回的 `ctx` 会在首个错误发生时被取消。

```go
import (
	"context"
	"errors"
	"fmt"

	"github.com/erdong01/kit/goWorker"
)

pool := goWorker.New(3)
ctx := context.Background()

group, ctx := pool.WithContext(ctx)

group.Go(func() error {
	// do something
	return nil
})

group.Go(func() error {
	return errors.New("boom")
})

if err := group.Wait(); err != nil {
	fmt.Println("err:", err)
}
<-ctx.Done() // ctx 已取消
```

## 注意事项

- `New` 的 `count` 需大于 0。传入 0 会导致任务提交阻塞。
- 任务内 panic 会导致进程退出，若需要更温和的错误处理，请在任务函数内部自行 `recover`。

## 测试案例

推荐执行：

```bash
go test ./goWorker -run Test
```

示例测试（仓库内：`goWorker/goWorker_cases_test.go`）包含：

- 并发上限是否生效
- `Worker.Wait` 是否等待所有任务完成
- `Group.WithContext` 的错误与取消行为
- `Group.TryGo` 的限流行为
