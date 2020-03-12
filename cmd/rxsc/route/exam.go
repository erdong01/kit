package route

import (
	"github.com/gin-gonic/gin"
	"rxt/cmd/rxsc/app/handler/exam"
)

// Route 学生端路由
func examGroup(group *gin.RouterGroup) {
	g := group.Group("/exam")
	{
		g.POST("/submit", exam.Submit)
	}
}
