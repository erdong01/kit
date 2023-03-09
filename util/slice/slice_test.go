package slice

import (
	"fmt"
	"testing"
)

// 通过查找切片元素删除
func TestDel(t *testing.T) {
	var data = []int{1, 2, 3}
	Del(&data, 1)
	fmt.Println("data", data)
}

// 通过切片下表删除
func TestDelByIndex(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}
	DelByIndex(&data, 10)
	fmt.Println("data", data)
}
