package examService

import (
	"github.com/erDong01/micro-kit/cmd/exam/logic/aiLogic"
	exam "github.com/erDong01/micro-kit/cmd/exam/proto/sc"
	"github.com/erDong01/micro-kit/cmd/exam/service/baseService"
)

type V1 struct {
	baseService.Service
}

func (c V1) Submit(exam *exam.ExamRequest) (str string, err error) {
	aiLogic.New().Dina("test")
	return "test", nil
}
