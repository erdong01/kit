package cache

import (
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestHas(t *testing.T) {
	Has("test")
}

func TestGet(t *testing.T) {
	key := "test1"
	Del(key)
	var num1, num2 int
	num1 = 11111
	Set(key, num1, DefaultTTl)
	err := Get(key, &num2)
	if err != nil {
		t.Log(err)
	}
	if num1 != num2 {
		t.Fail()
	}
}

func TestGetJSON(t *testing.T) {

	key := "test2"

	// 简单 map 转义
	Del(key)
	data := JSON{
		"name": "zx",
	}
	err := SetJSON(key, data, DefaultTTl)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	data2 := JSON{}
	GetJSON(key, &data2)

	if data["name"] != data2["name"] {
		t.Log("结果错误")
		t.Fail()
	}
	// 切片 map 转义
	data3 := []JSON{}
	data3 = append(data3, data)
	data3 = append(data3, data)
	data3 = append(data3, data)

	Del(key)
	err = SetJSON(key, data3, ForeverTTl)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	data4 := []JSON{}
	GetJSON(key, &data4)
	if len(data3) != len(data4) {
		t.Fail()
	}

	for _, value := range data4 {
		if value["name"] != data["name"] {
			t.Fail()
		}
	}

	Del(key)
	// 结构体
	data6 := []struct {
		Name string `json:"name"`
	}{}
	err = SetJSON(key, data3, ForeverTTl)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	GetJSON(key, &data6)
	for _, value := range data6 {
		if value.Name != data["name"] {
			t.Fail()
		}
	}
	Del(key)
}

func TestSetJSON(t *testing.T) {
	SetJSON("test2", JSON{
		"name": "hehe",
		"age":  13,
	}, DefaultTTl)
}
