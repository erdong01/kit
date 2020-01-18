package model

type ExamQuestionType struct {
	BaseModel
	ExamQuestionTypeId                    int64   `gorm:"primary_key" json:"exam_question_type_id,omitempty"`
	ExamNo                                int64   `json:"exam_no,omitempty"`
	QuestionTypeId                        int64   `json:"question_type_id,omitempty"`
	ExamPaperQuestionTypeOrder            int16   `json:"exam_paper_question_type_order,omitempty"`
	QuestionTypeCategoryId                int8    `json:"question_type_category_id,omitempty"`
	ExamQuestionTypeScore                 float32 `json:"exam_question_type_score,omitempty"`
	ExamQuestionTypeActualScore           float32 `json:"exam_question_type_actual_score,omitempty"`
	ExamQuestionTypeCount                 int16   `json:"exam_question_type_count,omitempty"`
	ExamQuestionTypeErrorCount            int16   `json:"exam_question_type_error_count,omitempty"`
	ExamQuestionTypeCorrectCount          int16   `json:"exam_question_type_correct_count,omitempty"`
	ExamQuestionTypeKnowledgeCount        int16   `json:"exam_question_type_knowledge_count,omitempty"`
	ExamQuestionTypeKnowledgeErrorCount   int16   `json:"exam_question_type_knowledge_error_count,omitempty"`
	ExamQuestionTypeKnowledgeCorrectCount int16   `json:"exam_question_type_knowledge_correct_count,omitempty"`
	QuestionType	QuestionType `gorm:"ForeignKey:question_type_id;AssociationForeignKey:question_type_id"`
}

func (ExamQuestionType) TableName() string {
	return "rxt_exam_question_type"
}
