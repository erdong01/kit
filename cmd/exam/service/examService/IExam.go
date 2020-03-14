package examService

import (
	exam "rxt/cmd/exam/proto/sc"
	"rxt/internal/core"
)

type Exam struct {
	IExam
}

func New(param ...*core.Core) Exam {
	exam := &V1{}
	exam.Init(param...)
	return Exam{exam}
}

type IExam interface {
	Submit(exam *exam.ExamRequest) (str string, err error)
}
