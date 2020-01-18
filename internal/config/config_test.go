package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {

	var  rawVal int
	rawVal = 111
	aa:=reflect.TypeOf(rawVal)
	if aa.Name() == "int"{
		fmt.Println("tseete")
	}

}
