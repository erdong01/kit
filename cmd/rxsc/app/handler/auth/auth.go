package auth

import (
	"github.com/gin-gonic/gin"
	auth "rxt/cmd/auth/proto/sc"
	"rxt/internal/http/restful"
	"rxt/rpc/client/sc/authClient"
)

func Login(context *gin.Context) {
	var AuthRequest *auth.AuthRequest
	context.BindJSON(&AuthRequest)
	res, err := authClient.New().Logic(AuthRequest)
	if err != nil {
		restful.Exception(context, err)
		return
	}
	restful.Success(context, res)
}
