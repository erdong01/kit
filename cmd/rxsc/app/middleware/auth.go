package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	auth "rxt/cmd/auth/proto/sc"
	"rxt/internal/http/restful"
	"rxt/internal/jwt"
	"rxt/internal/wrong"
	"rxt/rpc/client/sc/authClient"
)

// Auth 测试鉴权中间件
func Auth(context *gin.Context) {
	token, err := jwt.GetToken(context)
	if err != nil || token == "" {
		restful.Exception(context, wrong.New(http.StatusUnauthorized, err, "请输入Token"))
		return
	}

	result, err := authClient.New().Validate(&auth.TokenRequest{Token: token})
	if err != nil || result == nil {
		restful.Exception(context, wrong.New(http.StatusUnauthorized, err, "Token解析失败"))
		return
	}

	if result.StudentUserNo == 0 {
		restful.Exception(context, wrong.New(http.StatusUnauthorized, errors.New("获取用户失败"), "获取用户失败"))
		return
	}

	context.Set("auth", result)
	context.Next()
}
