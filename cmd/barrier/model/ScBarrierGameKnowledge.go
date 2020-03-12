package model

// ScBarrierGameKnowledge 定义字段
type ScBarrierGameKnowledge struct {
	BarrierGameKnowledgeID          int64   `gorm:"primary_key" json:"barrier_game_knowledge_id,omitempty"`
	BarrierGameID                   int64   `json:"barrier_game_id,omitempty"`
	KnowledgeNo                     int64   `json:"knowledge_no,omitempty"`
	BarrierKnowledgeOrder           int64   `json:"barrier_knowledge_order,omitempty"`
	BarrierGameKnowledgeStatus      int64   `json:"barrier_game_knowledge_status,omitempty"`
	BarrierGameKnowledgeScore       int64   `json:"barrier_game_knowledge_score,omitempty"`
	BarrierGameKnowledgeActualScore float64 `json:"barrier_game_knowledge_actual_score,omitempty"`
	BarrierGameKnowledgeLearnMinute float64 `json:"barrier_game_knowledge_learn_minute,omitempty"`
	BaseModel
}

// TableName 设置表名
func (t ScBarrierGameKnowledge) TableName() string {
	return "rxt_sc_barrier_game_knowledge"
}
