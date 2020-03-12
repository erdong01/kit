package baseService

import (
	"fmt"
	"rxt/internal/core/register"
	"testing"
)

func TestGlobalBegin(t *testing.T) {
	new(register.Register).ConfigRegister().
		RedisRegister().
		DbRegister().
		FacadeCacheRegister().
		SetPort(5001)
	s := GlobalBegin()
	fmt.Print(s)
}
