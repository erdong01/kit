package main

import (
	"github.com/erDong01/micro-kit/db/mongoDB/dbo"
	"reflect"
	"time"
)

type Test struct {
	Name string
	User *User
}
type User struct {
	Age  int
	Name string
}

var test *Test

func main() {
	dbo.Init("mongodb://127.0.0.1:27017")
	go func() {
		for i := 0; i < 10; i++ {
			u := User{
				Age:  i,
				Name: "test",
			}
			dbo.MgoClient.TestDB.Insert("user", &u)
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			u := User{
				Age:  i,
				Name: "test",
			}
			dbo.MgoClient.MasterDB.Insert("user", &u)
		}

	}()
	time.Sleep(time.Second * 5)
	test = &Test{}
	sysConfig := reflect.ValueOf(test).Elem()
	//sysType := reflect.TypeOf(test).Elem()
	//fmt.Println(sysConfig.FieldByIndex([]int{1})) //打印User层属性
	//fmt.Println(sysConfig.FieldByName("Name"))
	//fmt.Println(sysConfig.FieldByName("Age")) //直接打印下一层不成功
	//sysConfigNew := reflect.New(sysType)
	//sysConfigNew.Elem().FieldByName("Name").SetString("test")
	//fmt.Println(sysConfigNew.Elem().FieldByIndex([]int{1}))
	//fmt.Println(sysConfigNew.Elem().FieldByName("Name"))
	//fmt.Println(test.Name)
	user := reflect.ValueOf(&User{})
	flag := sysConfig.FieldByName("Age") == reflect.Value{}
	if flag {
		//fmt.Println(flag)
		user.Elem().FieldByName("Name").SetString("testset")
		sysConfig.Field(1).Set(user)
		//fmt.Println(sysConfig.FieldByIndex([]int{1})) //12
	}
	//fmt.Println(test.User)
}
