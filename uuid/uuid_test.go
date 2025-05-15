package uuid

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGetUUID(t *testing.T) {
	uid := New()
	fmt.Println("uid:", uid)
}

func TestGenerateNumber(t *testing.T) {
	var noMap sync.Map
	var i int
	for i <= 10000 {
		go func() {
			uid := GenerateNumber(8)
			go func(id string) {
				noMap.Store(id, struct{}{})
			}(uid)
		}()
		i++
	}
	var ii int
	time.Sleep(time.Second * 30)
	noMap.Range(func(key, value any) bool {
		ii++
		return true
	})
	fmt.Println(ii)

}
