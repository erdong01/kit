package model

// ScBarrierGame 定义字段
type ScBarrierGame struct {
	BarrierGameID       int64   `gorm:"primary_key" json:"barrier_game_id,omitempty"`
	BarrierGameName     string  `json:"barrier_game_name,omitempty"`
	SubjectID           int64   `json:"subject_id,omitempty"`
	GradeChildrenID     int8    `json:"grade_children_id,omitempty"`
	BookNo              int64   `json:"book_no,omitempty"`
	BarrierGameType     int64   `json:"barrier_game_type,omitempty"`
	BarrierGameStatus   int64   `json:"barrier_game_status,omitempty"`
	EachBarrierQuantity int64   `json:"each_barrier_quantity,omitempty"`
	BarrierSum          int64   `json:"barrier_sum,omitempty"`
	BarrierFinishSum    int64   `json:"barrier_finish_sum,omitempty"`
	GameScore           float64 `json:"game_score,omitempty"`
	GameActualScore     float64 `json:"game_actual_score,omitempty"`
	GameForDays         int64   `json:"game_for_days,omitempty"`
	BarrierGameOrder    int64   `json:"barrier_game_order,omitempty"`
	BaseModel
	KnowledgeList []ScBarrierGameKnowledge
}

// TableName 设置表名
func (t ScBarrierGame) TableName() string {
	return "rxt_sc_barrier_game"
}
