package config

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {

	var rawVal int
	rawVal = 111
	aa := reflect.TypeOf(rawVal)
	if aa.Name() == "int" {
		fmt.Println("tseete")
	}

}
func TestBig(t *testing.T) {
	test :=decimal.NewFromFloat(1.1).Div(decimal.NewFromFloat(3)).Round(2)
	fmt.Print(test)
}
