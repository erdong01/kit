package model

type ScStudentKnowledge struct {
	BaseModel
	StudentKnowledgeId               int64   `gorm:"primary_key" json:"student_knowledge_id,omitempty"`
	StudentUserNo                    int64   `json:"student_user_no,omitempty"`
	StudentKnowledgeType             int8    `json:"student_knowledge_type,omitempty"`
	KnowledgeNo                      int64   `json:"knowledge_no,omitempty"`
	StudentKnowledgeProficiency      float64 `json:"student_knowledge_proficiency,omitempty"`
	StudentKnowledgeDiff             float64 `json:"student_knowledge_diff,omitempty"`
	StudentKnowledgeFirstProficiency float64 `json:"student_knowledge_first_proficiency,omitempty"`
	StudentKnowledgeExamCount        int     `json:"student_knowledge_exam_count,omitempty"`
	StudentKnowledgeWeakStatus       int8    `json:"student_knowledge_weak_status,omitempty"`
	IsHistoryWeak                    int8    `json:"is_history_weak,omitempty"`
	IsPre                            int8    `json:"is_pre,omitempty"`
	StudentKnowledgeRelateProb       int8    `json:"student_knowledge_relate_prob,omitempty"`
}

func (ScStudentKnowledge) TableName() string {
	return "rxt_sc_student_knowledge"
}
