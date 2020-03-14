package examService

import (
	"rxt/cmd/exam/logic/aiLogic"
	exam "rxt/cmd/exam/proto/sc"
	"rxt/cmd/exam/service/baseService"
)

type V1 struct {
	baseService.Service
}

func (c V1) Submit(exam *exam.ExamRequest) (str string, err error) {
	aiLogic.New().Dina("test")
	return "test", nil
}
