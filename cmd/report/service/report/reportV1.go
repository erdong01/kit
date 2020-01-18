package report

import (
	"net/http"
	"rxt/cmd/report/model"
	"rxt/cmd/report/service/base"
	"rxt/internal/wrong"

	"github.com/jinzhu/gorm"
)

type ServiceV1 struct {
	base.Service
}

func (s *ServiceV1) Show(examNo int64) (ShowReport, error) {
	var report ShowReport
	var exam model.Exam
	err := s.Db.Model(&exam).Preload("ExamQuestionType", func(db *gorm.DB) *gorm.DB {
		return db.Select("rxt_exam_question_type.exam_no,rxt_exam_question_type.question_type_id")
	}).
		Preload("ExamQuestionType.QuestionType", func(db *gorm.DB) *gorm.DB {
			return db.Select("question_type_id,question_type_name")
		}).
		Select("exam_no,exam_name,exam_id").
		Where(&model.Exam{ExamNo: examNo, ExamStatus: 3}).
		First(&exam).Error
	if err != nil {
		err = wrong.New(http.StatusExpectationFailed, err, "考试不存在!")
	}
	report.Exam = exam
	return report, err
}
