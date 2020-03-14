package exam

import (
	"github.com/gin-gonic/gin"
	exam "rxt/cmd/exam/proto/sc"
	"rxt/internal/http/restful"
	"rxt/rpc/client/sc/examClient"
)

// Submit 提交阅卷
func Submit(c *gin.Context) {
	var examRequest *exam.ExamRequest
	if err := c.Bind(&examRequest); err != nil {
		restful.Exception(c, err)
		return
	}
	res, err := examClient.New().Submit(examRequest)
	// 10秒超时 仅remote call有效
	if err != nil {
		restful.Exception(c, err)
	}
	//tx.Commit()
	restful.Success(c, res)
	return
}
