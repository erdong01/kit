package model

type ScStudentQuestion struct {
	BaseModel
	StudentQuestionId           int64 `gorm:"primary_key" json:"student_question_id,omitempty"`
	StudentUserNo               int64 `json:"student_user_no,omitempty"`
	QuestionNo                  int64 `json:"question_no,omitempty"`
	ReviseStatus                int8  `json:"revise_status,omitempty"`
	StudentQuestionExamCount    int   `json:"student_question_exam_count,omitempty"`
	StudentQuestionErrorCount   int   `json:"student_question_error_count,omitempty"`
	StudentQuestionCorrectCount int   `json:"student_question_correct_count,omitempty"`
}

func (ScStudentQuestion) TableName() string {
	return "rxt_sc_student_question"
}
