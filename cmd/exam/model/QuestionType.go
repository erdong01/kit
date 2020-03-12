package model

type QuestionType struct {
	BaseModel
	QuestionTypeId         int64  `gorm:"primary_key" json:"question_type_id,omitempty"`
	QuestionTypeName       string `json:"question_type_name,omitempty"`
	QuestionTypeCategoryId int8   `json:"question_type_category_id,omitempty"`
	SubjectId              int64  `json:"subject_id,omitempty"`
	PhaseId                int64  `json:"phase_id,omitempty"`
	QuestionTypeSort       int    `json:"question_type_sort,omitempty"`
}

func (QuestionType) TableName() string {
	return "rxt_question_type"
}
