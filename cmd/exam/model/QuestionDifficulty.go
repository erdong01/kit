package model

type QuestionDifficulty struct {
	BaseModel
	QuestionDifficultyId   int64  `gorm:"primary_key" json:"question_difficulty_id,omitempty"`
	QuestionDifficultyName string `json:"question_difficulty_name,omitempty"`
}

func (QuestionDifficulty) TableName() string {
	return "rxt_question_difficulty"
}
