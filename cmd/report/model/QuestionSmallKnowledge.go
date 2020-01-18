package model

type QuestionSmallKnowledge struct {
	BaseModel
	QuestionSmallId int64 `gorm:"primary_key" json:"question_small_id,omitempty"`
	KnowledgeNo     int64 `json:"knowledge_no,omitempty"`
}

func (QuestionSmallKnowledge) TableName() string {
	return "rxt_question_small_knowledge"
}
