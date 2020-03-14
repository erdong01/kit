package route

import (
	"rxt/cmd/rxstudent/route/middleware"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(g *gin.Engine) *gin.Engine {
	g.Use(middleware.Cors())

	authGroup(g.Group("/v1"))

	// examGroup(g.Group("/v1", middleware.Auth))
	uploadGroup(g.Group("/v1", middleware.Auth))

	return g
}
