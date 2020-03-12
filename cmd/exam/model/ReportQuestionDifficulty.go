package model

type ReportQuestionDifficulty struct {
	BaseModel
	ReportQuestionDifficultyId     int64   `gorm:"primary_key" json:"report_question_difficulty_id,omitempty"`
	ExamNo                         int64   `json:"exam_no,omitempty"`
	QuestionDifficultyId           int64   `json:"question_difficulty_id,omitempty"`
	DifficultyQuestionCount        int16   `json:"difficulty_question_count,omitempty"`
	DifficultyQuestionErrorCount   int16   `json:"difficulty_question_error_count,omitempty"`
	DifficultyQuestionCorrectCount int16   `json:"difficulty_question_correct_count,omitempty"`
	DifficultyQuestionScore        float64 `json:"difficulty_question_score,omitempty"`
	DifficultyQuestionActualScore  float64 `json:"difficulty_question_actual_score,omitempty"`
}

func (ReportQuestionDifficulty) TableName() string {
	return "rxt_report_question_difficulty"
}
