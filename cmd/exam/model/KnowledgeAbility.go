package model

type KnowledgeAbility struct {
	BaseModel
	KnowledgeAbilityId     int64  `gorm:"primary_key" json:"knowledge_ability_id,omitempty"`
	KnowledgeAbilityName   string `json:"knowledge_ability_name,omitempty"`
	KnowledgeAbilityRemark string `json:"knowledge_ability_remark,omitempty"`
	PhaseId                int64  `json:"phase_id,omitempty"`
	SubjectId              int64  `json:"subject_id,omitempty"`
}

func (KnowledgeAbility) TableName() string {
	return "rxt_knowledge_ability"
}
