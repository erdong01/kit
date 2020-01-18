package route

import (
	"rxt/cmd/rxsc/app/handler/report"

	"github.com/gin-gonic/gin"
)

// Route 学生端路由
func reportGroup(group *gin.RouterGroup) {
	g := group.Group("/report")
	{
		g.GET("/show", report.Show)
	}
}
