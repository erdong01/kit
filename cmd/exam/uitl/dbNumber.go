package uitl

import (
	"rxt/cmd/exam/model"
	"rxt/internal/core"
)

func ExamNo() int64 {
	Number := GenerateRandomNumber(0, 99999999, 1)
	exam := model.Exam{}
	core.New().Db.Table("rxt_exam").Where("exam_no", Number[0]).First(&exam)
	if exam.ExamNo > 0 {
		ExamNo()
	}
	return int64(Number[0])
}
