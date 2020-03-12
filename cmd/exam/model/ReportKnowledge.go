package model

type ReportKnowledge struct {
	BaseModel
	ReportKnowledgeId           int64   `gorm:"primary_key" json:"report_knowledge_id,omitempty"`
	ExamNo                      int64   `json:"exam_no,omitempty"`
	KnowledgeNo                 int64   `json:"knowledge_no,omitempty"`
	ReportKnowledgeIsWeak       int8    `json:"report_knowledge_is_weak,omitempty"`
	ReportKnowledgeProficiency  float64 `json:"report_knowledge_proficiency,omitempty"`
	ReportKnowledgeDiff         float64 `json:"report_knowledge_diff,omitempty"`
	ReportKnowledgeScore        float64 `json:"report_knowledge_score,omitempty"`
	ReportKnowledgeActualScore  float64 `json:"report_knowledge_actual_score,omitempty"`
	ReportKnowledgeExamCount    int16   `json:"report_knowledge_exam_count,omitempty"`
	ReportKnowledgeErrorCount   int16   `json:"report_knowledge_error_count,omitempty"`
	ReportKnowledgeCorrectCount int16   `json:"report_knowledge_correct_count,omitempty"`
}

func (ReportKnowledge) TableName() string {
	return "rxt_report_knowledge"
}
