package register

import (
	"fmt"
	"github.com/erDong01/micro-kit/core"
	"testing"
)

func TestGlobalInit(t *testing.T) {
	GlobalInit()
	fmt.Println(core.New().GetName())
}
