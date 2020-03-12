package model

type ExamQuestionSmall struct {
	BaseModel
	ExamQuestionSmallId          int64   `gorm:"primary_key" json:"exam_question_small_id,omitempty"`
	ExamQuestionId               int64   `json:"exam_question_id,omitempty"`
	QuestionNo                   int64   `json:"question_no,omitempty"`
	QuestionSmallId              int64   `json:"question_small_id,omitempty"`
	ExamQuestionSmallAnswer      string  `json:"exam_question_small_answer,omitempty"`
	ExamQuestionSmallScore       float64 `json:"exam_question_small_score,omitempty"`
	ExamQuestionSmallActualScore float64 `json:"exam_question_small_actual_score,omitempty"`
	ExamQuestionSmallIsRight     int8    `json:"exam_question_small_is_right,omitempty"`
	ExamQuestionSmallOrder       int32 `json:"exam_question_small_order,omitempty"`
	ExamOtherError               float32 `json:"exam_other_error,omitempty"`
}

func (ExamQuestionSmall) TableName() string {
	return "rxt_exam_question_small"
}
