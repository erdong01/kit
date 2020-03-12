package questionLogic

import (
	"rxt/cmd/exam/logic/base"
	"rxt/cmd/exam/model"
)

type V1 struct {
	base.Logic
}

func (c V1) EditScStudentQuestion(scStudentQuestionMap map[int64]model.ScStudentQuestion, studentUserNo int64, questionNo int64, questionIsRight int) {
	var scStudentQuestion model.ScStudentQuestion
	if _, ok := scStudentQuestionMap[questionNo]; ok {
		scStudentQuestion = scStudentQuestionMap[questionNo]
	} else {
		scStudentQuestion = model.ScStudentQuestion{
			StudentUserNo: studentUserNo,
			QuestionNo:    questionNo,
		}
	}

	if questionIsRight == 1 {
		scStudentQuestion.StudentQuestionCorrectCount++
	} else {
		scStudentQuestion.StudentQuestionErrorCount++
	}
	scStudentQuestion.StudentQuestionExamCount++
	scStudentQuestionMap[questionNo] = scStudentQuestion

	if scStudentQuestion.StudentQuestionId > 0 {
		c.Transaction.Save(&scStudentQuestion)
	} else {
		c.Transaction.Create(&scStudentQuestion)
	}

}
