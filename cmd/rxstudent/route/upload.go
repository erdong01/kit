package route

import (
	"rxt/cmd/rxstudent/app/handler/upload"

	"github.com/gin-gonic/gin"
)

// Route 学生端路由
func uploadGroup(group *gin.RouterGroup) {
	// 用户登录管理
	uploadRoute := group.Group("/file-upload")
	{
		uploadRoute.GET("/token", upload.Uptoken)
		uploadRoute.POST("/attachment-info", upload.StoreAttachmentInfo)
	}
}
