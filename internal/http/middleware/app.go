package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Auth 测试鉴权中间件
func Auth(context *gin.Context) {
	fmt.Print("before\n")
	context.Next()
	fmt.Print("after\n")
}
