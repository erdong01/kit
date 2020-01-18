package auth

import (
	"github.com/gin-gonic/gin"
	auth "rxt/cmd/auth/proto/sc"
	"rxt/internal/api"
	"rxt/internal/http/restful"
)

func Login(context *gin.Context) {
	var AuthRequest = auth.AuthRequest{}
	context.BindJSON(&AuthRequest)
	res, err := api.New("auth", "Login", false).Call(&AuthRequest)
	if err != nil {
		restful.Exception(context, err)
		return
	}
	restful.Success(context, res)
}
