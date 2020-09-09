package util

import (
	"fmt"
	"testing"
)

func TestStr_Float64(t *testing.T) {
	f := Str("0.1").Float64()
	fmt.Println(f)
}

func TestStr_Bytes(t *testing.T) {
	b := Str("世界你好").Bytes()
	fmt.Println(b)
}
func TestStr_ToJSON(t *testing.T) {
	var a interface{}
	json := Str("世界你好").ToJSON(&a)
	fmt.Println(json, a)
}
