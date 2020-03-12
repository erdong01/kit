package barrier

import (
	"errors"
	"net/http"
	pb "rxt/cmd/barrier/proto/sc"
	"rxt/internal/http/restful"
	"rxt/internal/wrong"
	"rxt/rpc/client/sc/barrierClient"

	"github.com/gin-gonic/gin"
)

// Skip 跳关
func Skip(context *gin.Context) {
	// auth, _ := context.MustGet("auth").(*auth.UserResponse)

	var request pb.Request
	if err := context.BindUri(&request); err != nil {
		wrong.New(http.StatusExpectationFailed, err, "参数错误")
		return
	}

	result, err := barrierClient.New().Skip(&request)
	if err != nil {
		restful.Exception(context, err)
		return
	}

	if !result.GetResult() {
		restful.Exception(context, wrong.New(http.StatusExpectationFailed, errors.New("操作失败"), "操作失败"))
		return
	}

	restful.SuccessCreated(context, gin.H{})
}
