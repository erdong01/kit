package workservice

import (
	pb "rxt/cmd/exam/proto/student"
	"time"
)

type WorkService struct {
	IWork
}

func New() WorkService {
	work := &Work{}
	work.Init()
	return WorkService{work}
}

type IWork interface {
	Classwork(studentUserNo int, time time.Time) (*pb.ClassworkResponse, error)
}
