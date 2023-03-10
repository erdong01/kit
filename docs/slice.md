切片
===

### 删除切片中元素



```go
import (
	"fmt"
	"testing"

	"github.com/erdong01/kit/util/slice"
)

// 通过查找切片元素删除
func TestDel(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}
	slice.Del(&data,3)
	fmt.Println("data", data)
}

//输出：  data [1 2 4 5]

// 通过切片下标删除
func TestDelByIndex(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5}
	slice.DelByIndex(&data, 2)
	fmt.Println("data", data)
}
//输出：data [1 2 4 5]
```



### 切片去重

```go

import (
	"fmt"
	"testing"

	"github.com/erdong01/kit/util/slice"
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