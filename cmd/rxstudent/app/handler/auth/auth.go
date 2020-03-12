package auth

import (
	"github.com/gin-gonic/gin"
	auth "rxt/cmd/auth/proto/student"
	"rxt/internal/http/restful"
	"rxt/rpc/client/student/authClient"
)

// 用户登录
func Login(context *gin.Context) {
	var logicRequest *auth.LogicRequest
	if err := context.Bind(&logicRequest); err != nil {
		restful.Exception(context, err)
		return
	}
	res, err := authClient.New().Logic(logicRequest)
	if err != nil {
		restful.Exception(context, err)
	}
	restful.Success(context, res)
	return
}
