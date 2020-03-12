package model

type ReportLatentKnowledge struct {
	BaseModel
	ReportLatentKnowledged           int64   `gorm:"primary_key" json:"report_latent_knowledge_id,omitempty"`
	ExamNo                           int64   `json:"exam_no,omitempty"`
	KnowledgeId                      int64   `json:"knowledge_id,omitempty"`
	KnowledgeNo                      int64   `json:"knowledge_no,omitempty"`
	ReportLatentKnowledgeProficiency float64 `json:"report_latent_knowledge_proficiency,omitempty"`
	RelateProb                       float64 `json:"relate_prob,omitempty"`
	ReportLatentKnowledgeType        int8    `json:"report_latent_knowledge_type,omitempty"`
	ReportLatentKnowledgeIsExceed    int8    `json:"report_latent_knowledge_is_exceed,omitempty"`
}

func (ReportLatentKnowledge) TableName() string {
	return "rxt_report_latent_knowledge"
}
