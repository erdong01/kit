package model

type QuestionSmallAnswerOption struct {
	BaseModel
	QuestionSmallAnswerOptionId int64  `gorm:"primary_key" json:"question_small_answer_option_id,omitempty"`
	QuestionSmallId             int64  `json:"question_small_id,omitempty"`
	OptionName                  string `json:"option_name,omitempty"`
	OptionContent               string `json:"option_content,omitempty"`
	IsRight                     int8   `json:"is_right,omitempty"`
}

func (QuestionSmallAnswerOption) TableName() string {
	return "rxt_question_small_answer_option"
}
