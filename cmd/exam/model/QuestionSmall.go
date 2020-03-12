package model

type QuestionSmall struct {
	BaseModel
	QuestionSmallId      int64   `gorm:"primary_key" json:"question_small_id,omitempty"`
	QuestionNo           int64   `json:"question_no,omitempty"`
	QuestionSmallContent string `json:"question_small_content,omitempty"`
	QuestionSmallKnowledge []QuestionSmallKnowledge `json:"knowledge,omitempty" gorm:"foreignkey:QuestionSmallId;association_foreignkey:QuestionSmallId"`
	QuestionSmallAnswerOption []QuestionSmallAnswerOption `json:"option,omitempty" gorm:"foreignkey:QuestionSmallId;association_foreignkey:QuestionSmallId"`
}

func (QuestionSmall) TableName() string {
	return "rxt_question_small"
}
