package report

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	pb "rxt/cmd/report/proto/report"
	"rxt/internal/http/restful"
	"rxt/rpc/client/sc/reportClient"
)

// Show 获取报告列表
func Show(context *gin.Context) {
	report, err := reportClient.New().Show(&pb.ReportRequest{ExamId: 19600605})
	if err != nil {

		logrus.Debug("test")
		//var buf [4096]byte
		//n := runtime.Stack(buf[:], false)
		//fmt.Println("file,line:",string(buf[:n]))
		//wrong.Stack(0)
		restful.Exception(context, err)
		return
	}

	//cache.Set("test", "test", 0)
	//var value string
	//err = cache.Get("test", &value)
	//if err != nil {
	//	fmt.Printf("try get key[foo] error[%s]\n", err.Error())
	//	// err_handler(err)
	//}
	restful.Success(context, gin.H{
		"data": report.ExamId,
	})
}
