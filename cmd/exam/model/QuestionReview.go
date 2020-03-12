package model

type QuestionReview struct {
	BaseModel
	QuestionReviewId           int64   `gorm:"primary_key" json:"question_review_id,omitempty"`
	QuestionNo                   int64   `json:"question_no,omitempty"`
	QuestionReviewTeachingMinute float32 `json:"question_review_teaching_minute,omitempty"`
}

func (QuestionReview) TableName() string {
	return "rxt_question_review"
}
