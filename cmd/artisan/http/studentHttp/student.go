package studentHttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rxt/cmd/artisan/http/errGroup"
	student "rxt/cmd/rxstudent/bootstrap"
	"rxt/cmd/rxstudent/route"
	"rxt/internal/core"
	"time"
)

func Start() *http.Server {
	httpServer := &http.Server{
		Addr:         ":5003",
		Handler:      studentRoute(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	errGroup.ErrGroup.Go(func() error {
		student.App().SetPort(5003)
		fmt.Println(core.New().GetPort())
		defer core.Close()
		return httpServer.ListenAndServe()
	})
	return httpServer
}
func studentRoute() http.Handler {
	e := gin.Default()
	route.Init(e) //注册路由
	return e
}
