package model

type KnowledgeMapping struct {
	BaseModel
	KnowledgeMappingId    int64   `gorm:"primary_key" json:"knowledge_mapping_id,omitempty"`
	SourceKnowledgeId     int64   `json:"source_knowledge_id,omitempty"`
	TargetKnowledgeId     int64   `json:"target_knowledge_id,omitempty"`
	SubjectId             int64   `json:"subject_id,omitempty"`
	KnowledgeMappingValue float64 `json:"knowledge_mapping_value,omitempty"`
}

func (KnowledgeMapping) TableName() string {
	return "rxt_knowledge_mapping"
}
