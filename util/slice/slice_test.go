package slice

import (
	"fmt"
	"testing"
)

// 通过查找切片元素删除
func TestDel(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}
	Del(&data, 1)
	fmt.Println("data", data)
}

// 通过切片下表删除
func TestDelByIndex(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}
	DelByIndex(&data, 2)
	fmt.Println("data", data)
}

type Ball struct {
	Name string
	Id   int
}

func TestUnique(t *testing.T) {
	var data = []int{1, 2, 2, 3, 3, 4, 5}
	Unique(&data)
	fmt.Println("data", data)

	var users = []Ball{Ball{Name: "红", Id: 1}, Ball{Name: "红", Id: 1}, Ball{Name: "绿", Id: 2}, Ball{Name: "黄", Id: 3}}
	Unique(&users)
	fmt.Println("users", users)
}
