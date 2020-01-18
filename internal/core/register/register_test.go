package register

import (
	"fmt"
	"rxt/internal/core"
	"testing"
)

func TestGlobalInit(t *testing.T) {
	GlobalInit()
	fmt.Println(core.New().GetName())
}
