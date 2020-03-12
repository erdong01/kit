package model

type ExamQuestion struct {
	BaseModel
	ExamQuestionId          int64    `gorm:"primary_key" json:"exam_question_id,omitempty"`
	ExamQuestionTypeId      int64    `json:"exam_question_type_id,omitempty"`
	QuestionNo              int64    `json:"question_no,omitempty"`
	ExamPaperQuestionOrder  int16    `json:"exam_paper_question_order,omitempty"`
	ExamQuestionScore       float64  `json:"exam_question_score,omitempty"`
	ExamQuestionActualScore float64  `json:"exam_question_actual_score,omitempty"`
	Question                Question `gorm:"foreignkey:QuestionNo;association_foreignkey:QuestionNo"`
}

func (ExamQuestion) TableName() string {
	return "rxt_exam_question"
}
