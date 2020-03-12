package model

type ExamQuestionSmallKnowledge struct {
	BaseModel
	ExamQuestionSmallKnowledgeId int64  `gorm:"primary_key" json:"exam_question_small_knowledge_id,omitempty"`
	ExamQuestionSmallId          int64  `json:"exam_question_small_id,omitempty"`
	QuestionSmallId              int64  `json:"question_small_id,omitempty"`
	KnowledgeNo                  int64  `json:"knowledge_no,omitempty"`
	ExamQuestionKnowledgeIsRight int8 `json:"exam_question_knowledge_is_right,omitempty"`
}

func (ExamQuestionSmallKnowledge) TableName() string {
	return "rxt_exam_question_small_knowledge"
}
