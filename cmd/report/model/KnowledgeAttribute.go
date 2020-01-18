package model

type KnowledgeAttribute struct {
	BaseModel
	KnowledgeAttributeId     int64   `gorm:"primary_key" json:"knowledge_attribute_id,omitempty"`
	KnowledgeNo              string  `json:"knowledge_no,omitempty"`
	KnowledgeImportanceType  int8    `json:"knowledge_importance_type,omitempty"`
	KnowledgeDifficultyType  int8    `json:"knowledge_difficulty_type,omitempty"`
	KnowledgeDemandId        int64   `json:"knowledge_demand_id,omitempty"`
	KnowledgeStudyMinute     float32 `json:"knowledge_study_minute,omitempty"`
	CountryId                int64   `json:"country_id,omitempty"`
	ProvinceId               int64   `json:"province_id,omitempty"`
	CityId                   int64   `json:"city_id,omitempty"`
	KnowledgeAttributeRemark string  `json:"knowledge_attribute_remark,omitempty"`
}

func (KnowledgeAttribute) TableName() string {
	return "rxt_knowledge_attribute"
}
