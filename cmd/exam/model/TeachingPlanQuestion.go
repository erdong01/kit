package model

type TeachingPlanQuestion struct {
	BaseModel
	TeachingPlanQuestionId        int64   `gorm:"primary_key" json:"teaching_plan_question_id,omitempty"`
	ChapterKnowledgeId            int64   `json:"chapter_knowledge_id,omitempty"`
	KnowledgeNo                   int64   `json:"knowledge_no,omitempty"`
	QuestionNo                    int64   `json:"question_no,omitempty"`
	TeachingPlanQuestionType      int8    `json:"teaching_plan_question_type,omitempty"`
	ExampleMinute                 float64 `json:"example_minute,omitempty"`
	TeachingPlanQuestionSort      int     `json:"teaching_plan_question_sort,omitempty"`
	TeachingPlanQuestionLevelType int8    `json:"teaching_plan_question_level_type,omitempty"`
	TmrUserNo                     int64   `json:"tmr_user_no,omitempty"`
}

func (TeachingPlanQuestion) TableName() string {
	return "rxt_teaching_plan_question"
}
