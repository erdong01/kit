package bootstrap

import (
	"rxt/internal/core/register"
)

func App(name, env, version string) {
	register.GlobalInit()
}
