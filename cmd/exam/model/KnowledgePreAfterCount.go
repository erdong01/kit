package model

import "github.com/jinzhu/gorm"

type KnowledgePreAfterCount struct {
	gorm.Model
	PaperAnalysisKnowledgeId int64 `gorm:"primary_key"json:"paper_analysis_knowledge_id,omitempty"`
	KnowledgeNo              int64 `json:"knowledge_no,omitempty"`
	EditionNo                int64 `json:"edition_no,omitempty"`
	IsRevision               int8  `json:"is_revision,omitempty"`
	KnowledgePreCount        int   `json:"knowledge_pre_count,omitempty"`
	KnowledgeAfterCount      int   `json:"knowledge_after_count,omitempty"`
}

func (KnowledgePreAfterCount) TableName() string {
	return "rxt_knowledge_pre_after_count"
}
