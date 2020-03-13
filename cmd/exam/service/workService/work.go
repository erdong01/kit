package workservice

import (
	pb "rxt/cmd/exam/proto/student"

	"rxt/cmd/exam/service/baseService"
	"time"
)

type Work struct {
	baseService.Service
}

func (w *Work) Classwork(studentUserNo int, time time.Time) (*pb.ClassworkResponse, error) {
	return &pb.ClassworkResponse{}, nil
}
