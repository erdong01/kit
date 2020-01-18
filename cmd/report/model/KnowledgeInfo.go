package model

type KnowledgeInfo struct {
	BaseModel
	KnowledgeInfoId         int64   `gorm:"primary_key" json:"knowledge_info_id,omitempty"`
	KnowledgeNo             string  `json:"knowledge_no,omitempty"`
	KnowledgeTeachingMinute int8    `json:"knowledge_teaching_minute,omitempty"`
	KnowledgeStudyTarget    int8    `json:"knowledge_study_target,omitempty"`
	KnowledgeBasicKnowledge int64   `json:"knowledge_basic_knowledge,omitempty"`
	KnowledgeInfoExtend     float32 `json:"knowledge_info_extend,omitempty"`
}

func (KnowledgeInfo) TableName() string {
	return "rxt_knowledge_info"
}
