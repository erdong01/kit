package bootstrap

import (
	"rxt/internal/core/register"
)

func App() *register.Register {
	return register.GlobalInit()
}
