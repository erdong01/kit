package route

import (
	"github.com/gin-gonic/gin"
	"rxt/cmd/rxstudent/app/handler/auth"
)

// Route 学生端路由
func authGroup(group *gin.RouterGroup) {
	// 用户登录管理
	authRoute := group.Group("/auth")
	authRoute.POST("/login", auth.Login)
}
