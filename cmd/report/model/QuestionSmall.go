package model

type QuestionSmall struct {
	BaseModel
	QuestionSmallId      int64   `gorm:"primary_key" json:"question_small_id,omitempty"`
	QuestionNo           int64   `json:"question_no,omitempty"`
	QuestionSmallContent float32 `json:"question_small_content,omitempty"`
}

func (QuestionSmall) TableName() string {
	return "rxt_question_review"
}
