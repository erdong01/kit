package knowledgeDao

import (
	"rxt/cmd/exam/model"
	"rxt/internal/core"
)

type Knowledge struct {
	IKnowledgeDao
}

func New(param ...*core.Core) Knowledge {
	res := &V1{}
	res.Init(param...)
	return Knowledge{res}
}

type IKnowledgeDao interface {
	EditionKnowledge(param Param) ([]Review, []ChapterKnowledge, map[int64]int)
	ListBySubjectId(subjectId int64, isReview int8, editionNo int64) []preKnowledge
	GetPaperAnalysis(whereMap map[string]interface{}) (model.PaperAnalysis, bool)
	EditionChapterTree(editionNo int64, gradeId int64, bookNo int64) []model.Chapter
}
type Param struct {
	TeachingPlanQuestionType      int8
	TeachingPlanQuestionLevelType int8
	EditionNo                     int64
	SubjectId                     int64
	IsReview                      int8
}
