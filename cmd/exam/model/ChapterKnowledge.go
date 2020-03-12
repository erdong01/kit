package model

type ChapterKnowledge struct {
	BaseModel
	ChapterKnowledgeId       int64                  `gorm:"primary_key" json:"chapter_knowledge_id,omitempty"`
	EditionNo                int64                  `json:"edition_no,omitempty"`
	KnowledgeNo              int64                  `json:"knowledge_no,omitempty"`
	ChapterNo                int64                  `json:"chapter_no,omitempty"`
	ChapterKnowledgeSort     int                    `json:"chapter_knowledge_sort,omitempty"`
	ChapterKnowledgeType     int8                   `json:"chapter_knowledge_type,omitempty"`
	SubjectId                int64                  `json:"subject_id,omitempty"`
	ChapterKnowledgeReview   int8                   `json:"chapter_knowledge_review,omitempty"`
	ChapterKnowledgeEmphasis int8                   `json:"chapter_knowledge_emphasis,omitempty"`
	TeachingPlanQuestionMany []TeachingPlanQuestion `gorm:"ForeignKey:KnowledgeNo;AssociationForeignKey:KnowledgeNo"`
	Knowledge []Knowledge `gorm:"ForeignKey:KnowledgeNo;AssociationForeignKey:KnowledgeNo"`
}

func (ChapterKnowledge) TableName() string {
	return "rxt_chapter_knowledge"
}
