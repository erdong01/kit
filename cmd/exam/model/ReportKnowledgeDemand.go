package model

type ReportKnowledgeDemand struct {
	BaseModel
	ReportKnowledgeDemandId     int64   `gorm:"primary_key" json:"report_knowledge_demand_id,omitempty"`
	ExamNo                      int64   `json:"exam_no,omitempty"`
	KnowledgeDemand             int16   `json:"knowledge_demand,omitempty"`
	KnowledgeDemandExamCount    int16   `json:"knowledge_demand_exam_count,omitempty"`
	KnowledgeDemandErrorCount   int16   `json:"knowledge_demand_error_count,omitempty"`
	KnowledgeDemandCorrectCount int16   `json:"knowledge_demand_correct_count,omitempty"`
	KnowledgeDemandScore        float64 `json:"knowledge_demand_score,omitempty"`
	KnowledgeDemandActualScore  float64 `json:"knowledge_demand_actual_score,omitempty"`
}

func (ReportKnowledgeDemand) TableName() string {
	return "rxt_report_knowledge_demand"
}
