package model

type KnowledgeInfo struct {
	BaseModel
	KnowledgeInfoId         int64   `gorm:"primary_key" json:"knowledge_info_id,omitempty"`
	KnowledgeNo             string  `json:"knowledge_no,omitempty"`
	KnowledgeTeachingMinute float64 `json:"knowledge_teaching_minute,omitempty"`
	KnowledgeStudyTarget    string  `json:"knowledge_study_target,omitempty"`
	KnowledgeBasicKnowledge string  `json:"knowledge_basic_knowledge,omitempty"`
	KnowledgeInfoExtend     string  `json:"knowledge_info_extend,omitempty"`
}

func (KnowledgeInfo) TableName() string {
	return "rxt_knowledge_info"
}
