package knowledgeLogic

import (
	"rxt/cmd/exam/logic/base"
	"rxt/cmd/exam/model"
)

type V1 struct {
	base.Logic
}

func (c V1) ScStudentKnowledge(scStudentKnowledgeMap map[int64]model.ScStudentKnowledge,
	diff float64, value float64, knowledgeNo int64, studentUserNo int64) {
	var scStudentKnowledge model.ScStudentKnowledge

	if _, ok := scStudentKnowledgeMap[knowledgeNo]; ok {
		scStudentKnowledge = scStudentKnowledgeMap[knowledgeNo]
	} else {
		scStudentKnowledge = model.ScStudentKnowledge{
			StudentUserNo:        studentUserNo,
			StudentKnowledgeType: 1,
			KnowledgeNo:          knowledgeNo,

			StudentKnowledgeFirstProficiency: value,
			StudentKnowledgeWeakStatus:       2,
			IsHistoryWeak:                    2,
		}
	}
	scStudentKnowledge.StudentKnowledgeProficiency = value
	scStudentKnowledge.StudentKnowledgeDiff = diff
	if value <= 0.5 {
		scStudentKnowledge.StudentKnowledgeWeakStatus = 1
	}
	if scStudentKnowledge.StudentKnowledgeId > 0 {
		c.Transaction.Save(scStudentKnowledge)
	} else {
		c.Transaction.Create(scStudentKnowledge)
	}
}
