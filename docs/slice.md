切片
===
### 连个切片是否相等/差集
```go
import (
	"fmt"
	"testing"

	"github.com/erdong01/kit/slice"
)

func TestEqual(t *testing.T) {
	var data1 = []int{1, 2, 3, 4, 5, 6}
	var data2 = []int{1, 2, 3, 4, 5}
	b := slice.Equal(data1, data2)
	fmt.Println(b)
}
//输出： false


func TestDiff(t *testing.T) {
	var data1 = []int{ 1, 2, 3, 4, 5, 6}
	var data2 = []int{1, 2, 3, 4, 5}
	data := slice.Diff(data1, data2)
	fmt.Println("data",data)
}
//输出：  data [6]

```
### 在切片头部/任意位置插入元素
```go
import (
	"fmt"
	"testing"

	"github.com/erdong01/kit/slice"
)

func TestPrepend(t *testing.T) {
	var data = []int{ 3, 4, 5}
	slice.Prepend(&data,2)
	slice.Prepend(&data,1)
	fmt.Println("data", data)
}
//输出：  data [1 2 3 4 5]

// 任意位置插入元素
func TestInsert(t *testing.T) {
	var data = []int{1, 3, 4, 5}
	slice.Insert(&data, 1, 2)
	fmt.Println("data", data)
}
//输出：  data [1 2 3 4 5]

```

### 删除切片中元素



```go
import (
	"fmt"
	"testing"

	"github.com/erdong01/kit/slice"
)

// 通过查找切片元素删除
func TestDel(t *testing.T) {
	var data = []int{1, 2, 3,3, 4, 5}
	slice.Del(&data,3) //删除为3元素
	fmt.Println("data", data)
}

//输出：  data [1 2 4 5]

// 通过切片下标删除
func TestDelByIndex(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}
	slice.DelByIndex(&data, 2) //删除下标2的元素
	fmt.Println("data", data)
}
//输出：data [1 2 4 5]


// 删除切片多个元素
func TestDelFunc(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}           //数据
	var term = map[int]struct{}{2: {}, 3: {}} //条件
	slice.DelFunc(&data, func(i int) bool {
		_, ok := term[data[i]]
		return ok
	})
	fmt.Println("data", data)
}
//输出: data [1 4 5]
```



### 切片去重

```go

import (
	"fmt"
	"testing"

	"github.com/erdong01/kit/slice"
)

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
//输出 data [1 2 3 4 5]

//输出 users [{红 1} {绿 2} {黄 3}]

```