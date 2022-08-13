package core

import (
	"fmt"
	"testing"
)

func TestBind(t *testing.T) {
	var infoa info
	infoa.Name = "11"
	// New().Info = &info{}
	Bind(&New().Info, infoa)
	fmt.Println("111111:", New().Info.Name)

	var TestS TestS
	// New().Info = &info{}
	Bind(&New().TestI, TestS)
	fmt.Println("22222:", New().TestI.Name())
}
