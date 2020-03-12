package model

type KnowledgeReview struct {
	BaseModel
	KnowledgeReviewId int64 `gorm:"primary_key" json:"knowledge_review_id,omitempty"`
	KnowledgeNo       int64 `json:"knowledge_no,omitempty"`
	EditionNo         int64 `json:"edition_no,omitempty"`
	IsReviewKnowledge int64 `json:"is_review_knowledge,omitempty"`
}

func (KnowledgeReview) TableName() string {
	return "rxt_exam_question_type"
}
