package route

import (
	"rxt/cmd/rxsc/app/handler/barrier"

	"github.com/gin-gonic/gin"
)

// Route 知识点闯关路由
func barrierGroup(group *gin.RouterGroup) {
	g := group.Group("/barrier-game")
	{
		g.POST("/barrier-game-skip/:BarrierGameId/:KnowledgeNo", barrier.Skip)
	}
}
