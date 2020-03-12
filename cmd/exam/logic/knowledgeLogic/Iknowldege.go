package knowledgeLogic

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
	ScStudentKnowledge(scStudentKnowledgeMap map[int64]model.ScStudentKnowledge,
		diff float64, value float64, knowledgeNo int64, studentUserNo int64)
}
