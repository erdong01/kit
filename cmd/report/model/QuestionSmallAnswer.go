package model

type QuestionSmallAnswer struct {
	BaseModel
	QuestionSmallId     int64  `gorm:"primary_key" json:"question_small_id,omitempty"`
	QuestionSmallAnswer string `json:"question_small_answer,omitempty"`
	OptionCount         int8   `json:"option_count,omitempty"`
}

func (QuestionSmallAnswer) TableName() string {
	return "rxt_question_review"
}
