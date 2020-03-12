package tmrHttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rxt/cmd/artisan/http/errGroup"
	"rxt/cmd/rxtmr/route"
	"rxt/internal/core"
	"time"
)

func Start() *http.Server {
	httpServer := &http.Server{
		Addr:         ":5002",
		Handler:      tmrRoute(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errGroup.ErrGroup.Go(func() error {
		core.Make(
			core.Env(""),
			core.Name(""),
			core.Version(""),
			core.DbRegister(),
			core.RedisRegister(),
			core.ConfigRegister(),
			core.Port(5002),
		).Init()
		fmt.Println(core.New().GetPort())
		defer core.Close()
		return httpServer.ListenAndServe()
	})
	return httpServer
}
func tmrRoute() http.Handler {
	e := gin.Default()
	route.Init(e) //注册路由
	return e
}
