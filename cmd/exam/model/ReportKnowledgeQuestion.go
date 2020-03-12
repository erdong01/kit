package model

type ReportKnowledgeQuestion struct {
	BaseModel
	ReportKnowledgeQuestionId int64 `gorm:"primary_key" json:"report_knowledge_question_id,omitempty"`
	ReportKnowledgeId         int64 `json:"report_knowledge_id,omitempty"`
	QuestionNo                int64 `json:"question_no,omitempty"`
	ExamQuestionId            int64 `json:"exam_question_id,omitempty"`
}

func (ReportKnowledgeQuestion) TableName() string {
	return "rxt_report_knowledge_question"
}
