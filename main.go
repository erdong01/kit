package main

import (
	"fmt"
	"github.com/erDong01/micro-kit/db/mongoDB/dbo"
	"github.com/erDong01/micro-kit/db/mongoDB/drive"
	"reflect"
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
	drive.Init()
	dbo.Init("mongodb://111.229.20.134:27017")


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
		fmt.Println(flag)
		user.Elem().FieldByName("Name").SetString("testset")
		sysConfig.Field(1).Set(user)
		fmt.Println(sysConfig.FieldByIndex([]int{1})) //12
	}
	fmt.Println(test.User)
}
