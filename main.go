package main

import (
	"github.com/gin-gonic/gin"
	"rxt/cmd/rxsc/bootstrap"
	"rxt/cmd/rxsc/route"
	"rxt/internal/core"
)

func main() {
	engine := gin.Default()
	bootstrap.App("se", "sdf", "v1")
	defer core.Close()
	route.Init(engine)
	engine.Run()
}
