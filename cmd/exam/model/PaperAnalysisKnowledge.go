package model

import "github.com/jinzhu/gorm"

type PaperAnalysisKnowledge struct {
	gorm.Model
	PaperAnalysisKnowledgeId int64 `gorm:"primary_key"json:"paper_analysis_knowledge_id,omitempty"`
	KnowledgeNo              int64 `json:"knowledge_no,omitempty"`
	PaperAnalysisNo          int64 `json:"paper_analysis_no,omitempty"`
	Score                    int   `json:"score,omitempty"`
	Count                    int   `json:"count,omitempty"`
	PaperTypeId              int   `json:"paper_type_id,omitempty"`
	HighestDifficulty        int8  `json:"highest_difficulty,omitempty"`
	LowestDifficulty         int8  `json:"lowest_difficulty,omitempty"`
	AvgDifficulty            int   `json:"avg_difficulty,omitempty"`
}

func (PaperAnalysisKnowledge) TableName() string {
	return "rxt_paper_analysis_knowledge"
}
