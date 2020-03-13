package exam

import (
	"fmt"
	"net/http"
	pb "rxt/cmd/exam/proto/student"
	"rxt/internal/http/restful"
	"rxt/internal/util"
	"rxt/internal/wrong"
	"rxt/rpc/client/student/examclient"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
)

// Classwork 课堂作业列表接口
func Classwork(context *gin.Context) {
	t, err := ptypes.TimestampProto(util.WeekStartDay().AddDate(0, 0, -7))
	if err != nil {
		fmt.Println(err)
		return
	}

	req := &pb.ClassworkRequest{
		StudentUserNo: 75828069,
		Time:          t,
	}

	res, err := examclient.New().Classwork(req)
	if err != nil {
		restful.Exception(context, wrong.New(http.StatusInternalServerError, err, "操作失败"))
		return
	}

	restful.Success(context, res)

	// restful.Success(context, gin.H{})
}
