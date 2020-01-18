package model

type QuestionCategory struct {
	BaseModel
	QuestionCategoryId   int64  `gorm:"primary_key" json:"question_category_id,omitempty"`
	QuestionCategoryName string `json:"question_category_name,omitempty"`
}

func (QuestionCategory) TableName() string {
	return "rxt_question_category"
}
