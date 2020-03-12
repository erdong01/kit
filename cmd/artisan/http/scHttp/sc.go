package scHttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rxt/cmd/artisan/http/errGroup"
	scApp "rxt/cmd/rxsc/bootstrap"
	"rxt/cmd/rxsc/route"
	"rxt/internal/core"
	"time"
)

func Start() *http.Server {
	httpServer := &http.Server{
		Addr:         ":5001",
		Handler:      scRoute(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	errGroup.ErrGroup.Go(func() error {
		scApp.App("", "", "")
		fmt.Println(core.New().GetPort())
		defer core.Close()
		return httpServer.ListenAndServe()
	})
	return httpServer
}
func scRoute() http.Handler {
	e := gin.Default()
	route.Init(e)
	return e
}
