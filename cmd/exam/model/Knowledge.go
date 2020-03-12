package model

type Knowledge struct {
	BaseModel
	KnowledgeId       int64   `gorm:"primary_key" json:"knowledge_id,omitempty"`
	KnowledgeNo       int64   `json:"knowledge_no,omitempty"`
	KnowledgeName     string  `json:"knowledge_name,omitempty"`
	KnowledgeSort     int     `json:"knowledge_sort,omitempty"`
	KnowledgeParentId int64   `json:"knowledge_parent_id,omitempty"`
	KnowledgeTopId    int64   `json:"knowledge_top_id,omitempty"`
	HasChildren       int8    `json:"has_children,omitempty"`
	SubjectId         int8    `json:"subject_id,omitempty"`
	PhaseId           int8    `json:"phase_id,omitempty"`
	MappingValue      float32 `json:"mapping_value,omitempty"`
	KnowledgeTypeId   int64   `json:"knowledge_type_id,omitempty"`
	DefaultLearned    int8    `json:"default_learned,omitempty"`
	KnowledgeAttributeOne	KnowledgeAttribute `gorm:"foreignkey:KnowledgeNo;association_foreignkey:KnowledgeNo"`
}

func (Knowledge) TableName() string {
	return "rxt_knowledge"
}
