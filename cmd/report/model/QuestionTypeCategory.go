package model

type QuestionTypeCategory struct {
	BaseModel
	QuestionTypeCategoryId   int64  `gorm:"primary_key" json:"question_type_category_id,omitempty"`
	QuestionTypeCategoryName string `json:"question_type_category_name,omitempty"`
}

func (QuestionTypeCategory) TableName() string {
	return "rxt_question_type_category"
}
