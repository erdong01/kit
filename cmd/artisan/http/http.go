package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	scApp "rxt/cmd/rxsc/bootstrap"
	"rxt/cmd/rxsc/route"
	tmrApp "rxt/cmd/rxtmr/bootstrap"
	"rxt/internal/core"
	"rxt/internal/log"
	"time"
)

var (
	g errgroup.Group
)

func Sc() http.Handler {
	e := gin.Default()
	route.Init(e)
	return e
}

func Tmr() http.Handler {
	e := gin.Default()

	return e
}

func New() {
	server01 := &http.Server{
		Addr:         ":5001",
		Handler:      Sc(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	g.Go(func() error {
		scApp.App("", "", "")
		fmt.Println(core.New().GetPort())
		defer core.Close()
		return server01.ListenAndServe()
	})
	server02 := &http.Server{
		Addr:         ":5002",
		Handler:      Tmr(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	g.Go(func() error {
		tmrApp.App("", "", "")
		fmt.Println(core.New().GetPort())
		defer core.Close()
		return server02.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
