package questionDao

import (
	"rxt/cmd/exam/model"
	"rxt/internal/core"
)

type Question struct {
	IQuestionDao
}

func New(param ...*core.Core) Question {
	res := &V1{}
	res.Init(param...)
	return Question{res}
}

type IQuestionDao interface {
	ExamQuestionByNo(questionNo []int64) []model.Question
}
