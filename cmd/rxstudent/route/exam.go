package route

import (
	"rxt/cmd/rxstudent/app/handler/exam"

	"github.com/gin-gonic/gin"
)

// Route 学生端路由
func examGroup(group *gin.RouterGroup) {
	// 用户登录管理
	examRoute := group.Group("/exam")
	examRoute.GET("/classwork", exam.Classwork)
}
