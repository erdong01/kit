package model

type QuestionGrade struct {
	BaseModel
	QuestionNo int64 `gorm:"primary_key" json:"question_no,omitempty"`
	GradeId    int64 `json:"grade_id,omitempty"`
}

func (QuestionGrade) TableName() string {
	return "rxt_question_grade "
}
