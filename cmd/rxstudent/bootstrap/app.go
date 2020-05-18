package bootstrap

import (
	"github.com/erDong01/micro-kit/internal/core/register"
)

func App() *register.Register {
	return register.GlobalInit()
}
