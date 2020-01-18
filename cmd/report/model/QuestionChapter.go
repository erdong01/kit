package model

type QuestionChapter struct {
	BaseModel
	QuestionChapterId int64 `gorm:"primary_key" json:"question_chapter_id,omitempty"`
	QuestionNo        int64 `json:"question_no,omitempty"`
	ChapterNo         int64 `json:"chapter_no,omitempty"`
}

func (QuestionChapter) TableName() string {
	return "rxt_question_category"
}
