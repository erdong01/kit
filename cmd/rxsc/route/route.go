package route

import (
	"rxt/cmd/rxsc/app/middleware"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(g *gin.Engine) *gin.Engine {
	g.Use(middleware.Cors())

	authGroup(g.Group("/v1"))
	reportGroup(g.Group("/v1", middleware.Auth))
	examGroup(g.Group("/v1", middleware.Auth))

	return g
}
