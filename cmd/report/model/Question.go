package model

type Question struct {
	BaseModel
	QuestionId           int64  `gorm:"primary_key" json:"question_id,omitempty"`
	QuestionNo           int64  `json:"question_no,omitempty"`
	SubjectId            int64  `json:"subject_id,omitempty"`
	QuestionContent      string `json:"question_content,omitempty"`
	QuestionAnalysis     string `json:"question_analysis,omitempty"`
	QuestionDifficultyId int8   `json:"question_difficulty_id,omitempty"`
	QuestionTypeId       int64  `json:"question_type_id,omitempty"`
	QuestionCategoryId   int8   `json:"question_category_id,omitempty"`
	AvgScore             int8   `json:"avg_score,omitempty"`
	GroupCount           int    `json:"group_count,omitempty"`
	ResponseCount        int    `json:"response_count,omitempty"`
	QuestionStatus       int8   `json:"question_status,omitempty"`
	IsQuestionBest       int8   `json:"is_question_best,omitempty"`
	IsQuestionReview     int8   `json:"is_question_review,omitempty"`
	IsQuestionNew        int8   `json:"is_question_new,omitempty"`
	QuestionYear         int8   `json:"question_year,omitempty"`
}

func (Question) TableName() string {
	return "rxt_exam_question_type"
}
