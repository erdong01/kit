package model

type QuestionCreator struct {
	BaseModel
	QuestionCreatorId   int64 `gorm:"primary_key" json:"question_creator_id,omitempty"`
	CreatorUserNo       int64 `json:"creator_user_no,omitempty"`
	QuestionNo          int64 `json:"question_no,omitempty"`
	QuestionCreatorType int8 `json:"question_creator_type,omitempty"`
	CompanyNo           int64 `json:"company_no,omitempty"`
	CampusNo            int64 `json:"campus_no,omitempty"`
}

func (QuestionCreator) TableName() string {
	return "rxt_question_creator"
}
