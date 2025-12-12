package uuid

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetUUID(t *testing.T) {
	uid := New()
	fmt.Println("uid:", uid)
	fmt.Println("uid32:", GenerateNumber(32))
}

func TestGenerateNumber(t *testing.T) {
	var noMap sync.Map
	var dupMu sync.Mutex
	dup := 0
	var n int64
	deadline := time.After(60 * time.Second)
LOOP:
	for {
		select {
		case <-deadline:
			break LOOP
		default:
			go func() {
				id := GenerateNumber(18)
				fmt.Println("id:", id)
				atomic.AddInt64(&n, 1)
				if _, loaded := noMap.LoadOrStore(id, struct{}{}); loaded {
					dupMu.Lock()
					dup++
					dupMu.Unlock()
				}
			}()
		}
	}
	unique := 0
	noMap.Range(func(key, value any) bool {
		unique++
		return true
	})
	fmt.Printf("generated=%d unique=%d dup=%d\n", n, unique, dup)
}
