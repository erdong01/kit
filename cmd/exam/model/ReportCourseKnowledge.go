package model

type ReportCourseKnowledge struct {
	BaseModel
	ReportCourseKnowledgeId int64 `gorm:"primary_key" json:"report_course_knowledge_id,omitempty"`
	ReportCourseId          int64 `json:"report_course_id,omitempty"`
	KnowledgeNo             int64 `json:"knowledge_no,omitempty"`
}

func (ReportCourseKnowledge) TableName() string {
	return "rxt_report_course_knowledge"
}
