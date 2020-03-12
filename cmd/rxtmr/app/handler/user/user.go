package user

import (
	pb "rxt/cmd/report/proto/report"
	"rxt/internal/http/restful"
	"rxt/rpc/client/sc/reportClient"

	"github.com/gin-gonic/gin"
)

// Login 登录接口
func Login(context *gin.Context) {
	report, err := reportClient.New().Show(&pb.ReportRequest{ExamId: 19600605})
	if err != nil {
		restful.Exception(context, err)
		return
	}

	restful.Success(context, gin.H{
		"data": report.ExamId,
	})
}
