package report

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	pb "rxt/cmd/report/proto/report"
	"rxt/internal/api"
	"rxt/internal/http/restful"
	"rxt/internal/log"
)

// Show 获取报告列表
func Show(context *gin.Context) {
	report, err := api.New("report", "Show", false).Call(&pb.ReportRequest{ExamId: 19600605})
	if err != nil {

		logrus.Debug("test")
		//var buf [4096]byte
		//n := runtime.Stack(buf[:], false)
		//fmt.Println("file,line:",string(buf[:n]))
		//wrong.Stack(0)
		restful.Exception(context, err)
		return
	}
	result, ok := report.(*pb.ReportResponse)

	log.Error("test")
	if !ok {
		fmt.Println(result)
	}

	//cache.Set("test", "test", 0)
	//var value string
	//err = cache.Get("test", &value)
	//if err != nil {
	//	fmt.Printf("try get key[foo] error[%s]\n", err.Error())
	//	// err_handler(err)
	//}
	restful.Success(context, gin.H{
		"data": result.ExamId,
	})
}
