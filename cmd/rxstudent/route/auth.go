package route

import (
	"rxt/cmd/rxstudent/app/handler/auth"

	"github.com/gin-gonic/gin"
)

// Route 学生端路由
func authGroup(group *gin.RouterGroup) {
	// 用户登录管理
	authRoute := group.Group("/auth")
	authRoute.POST("/login", auth.Login)
}
