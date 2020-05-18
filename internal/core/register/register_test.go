package register

import (
	"fmt"
	"github.com/erDong01/gin-kit/internal/core"
	"testing"
)

func TestGlobalInit(t *testing.T) {
	GlobalInit()
	fmt.Println(core.New().GetName())
}
