package model

type QuestionSmallKnowledge struct {
	BaseModel
	QuestionSmallId int64 `gorm:"primary_key" json:"question_small_id,omitempty"`
	KnowledgeNo     int64 `json:"knowledge_no,omitempty"`
	KnowledgeAttributeOne KnowledgeAttribute`json:"knowledge_attribute,omitempty" gorm:"foreignkey:KnowledgeNo;association_foreignkey:KnowledgeNo"`
	Knowledge	Knowledge `gorm:"foreignkey:KnowledgeNo;association_foreignkey:KnowledgeNo"`
}

func (QuestionSmallKnowledge) TableName() string {
	return "rxt_question_small_knowledge"
}
