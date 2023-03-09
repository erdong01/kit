package slice

import (
	"fmt"
	"testing"

	"github.com/erdong01/kit/util/slice"
)

// 通过查找切片元素删除
func TestDel(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}
	slice.Del(&data, 3)
	fmt.Println("data", data)
}

// 通过切片下标删除
func TestDelByIndex(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}
	slice.DelByIndex(&data, 2)
	fmt.Println("data", data)
}
