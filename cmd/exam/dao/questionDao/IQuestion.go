package questionDao

import (
	"github.com/erDong01/micro-kit/cmd/exam/model"
	"github.com/erDong01/micro-kit/internal/core"
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
