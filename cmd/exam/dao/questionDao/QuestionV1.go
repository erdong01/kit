package questionDao

import (
	"github.com/jinzhu/gorm"
	"rxt/cmd/exam/dao/base"
	"rxt/cmd/exam/model"
)

type V1 struct {
	base.Dao
}

// 考试题目查询 通过编号
func (c V1) ExamQuestionByNo(questionNo []int64) []model.Question {
	var questionData []model.Question
	c.Db.Select("rxt_question.question_id,rxt_question.question_no,rxt_question.question_type_id,"+
		"rxt_question.question_difficulty_id,rxt_question.question_category_id,rxt_question_type.question_type_category_id").
		Where("rxt_question.question_no in (?)", questionNo).
		Preload("QuestionSmall").
		Preload("QuestionType").
		Preload("QuestionSmall.QuestionSmallKnowledge").
		Preload("QuestionSmall.QuestionSmallKnowledge.Knowledge").
		Preload("QuestionSmall.QuestionSmallKnowledge.KnowledgeAttributeOne", func(db *gorm.DB) *gorm.DB {
			return db.Joins("INNER JOIN rxt_knowledge_demand ON rxt_knowledge_demand.knowledge_demand_id = rxt_knowledge_attribute.knowledge_demand_id")
		}).
		Preload("QuestionSmall.QuestionSmallAnswerOption").
		Joins("INNER JOIN rxt_question_type ON rxt_question_type.question_type_id = rxt_question.question_type_id").
		Find(&questionData)
	return questionData
}
