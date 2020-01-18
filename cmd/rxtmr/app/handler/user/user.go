package user

import (
	pb "rxt/cmd/report/proto/report"
	"rxt/internal/api"
	"rxt/internal/http/restful"

	"github.com/gin-gonic/gin"
)

// Login 登录接口
func Login(context *gin.Context) {
	report, err := api.New("report", "Show", false).Call(&pb.ReportRequest{ExamId: 19600605})
	if err != nil {
		restful.Exception(context, err)
		return
	}

	restful.Success(context, gin.H{
		"data": report.(*pb.ReportResponse).ExamId,
	})
}
