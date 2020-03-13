package student

import (
	"context"
	pb "rxt/cmd/exam/proto/student"
	workservice "rxt/cmd/exam/service/workService"

	"github.com/golang/protobuf/ptypes"
)

// Server 服务类
type Server struct {
}

// Classwork 获取课堂作业
func (s Server) Classwork(ctx context.Context, request *pb.ClassworkRequest) (*pb.ClassworkResponse, error) {
	time, err := ptypes.Timestamp(request.Time)
	if err != nil {
		return nil, err
	}

	resp, err := workservice.New().Classwork(int(request.StudentUserNo), time)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
