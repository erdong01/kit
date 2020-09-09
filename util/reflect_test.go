package util

import (
	"fmt"
	"testing"
)

func TestReflect(t *testing.T) {
	test := Reflect("1.02").Float64()
	fmt.Println(test)
}
