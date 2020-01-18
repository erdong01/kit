package bootstrap

import (
	"rxt/internal/core"
)

func App(name, env, version string) {
	core.Make(
		core.Env(env),
		core.Name(name),
		core.Version(version),
		core.DbRegister(),
		core.RedisRegister(),
		core.ConfigRegister(),
		core.Port(5002),
		//core.Engine(route.Init),
	).Init()
}
