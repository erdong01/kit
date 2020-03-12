package questionLogic

import (
	"rxt/cmd/exam/model"
	"rxt/internal/core"
)

type Container struct {
	I
}

func New(param ...*core.Core) Container {
	res := &V1{}
	res.Init(param...)
	return Container{res}
}

type I interface {
	EditScStudentQuestion(scStudentQuestionMap map[int64]model.ScStudentQuestion, studentUserNo int64, questionNo int64, questionIsRight int)
}
