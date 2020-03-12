package examService

import (
	"rxt/cmd/exam/logic/aiLogic"
	"rxt/cmd/exam/model"
	exam "rxt/cmd/exam/proto/sc"
	"rxt/internal/core"
)

type Exam struct {
	IExam
}

func New(param ...*core.Core) Exam {
	exam := &V1{}
	exam.Init(param...)
	return Exam{exam}
}

type IExam interface {
	Submit(exam *exam.ExamRequest) (param Param, err error)
	ReportCreate(param Param) (err error)
	CreateScExamStudent(examNo int64, studentUserNo int64, bookNo int64) error
}

//----- submit() start -----
type Question struct {
	QuestionNo int64 `json:"question_no"`
	QuestionId int64 `json:"question_id"`
}
type Knowledge struct {
	model.Knowledge
	model.ReportKnowledge
	Question []Question
}
type Demand struct {
	KnowledgeDemand             int16
	KnowledgeDemandExamCount    int16
	KnowledgeDemandErrorCount   int16
	KnowledgeDemandCorrectCount int16
	KnowledgeDemandScore        float64
	KnowledgeDemandActualScore  float64
	Question                    map[int64]int64
	KnowledgeNo                 int64
}

type QuestionDifficulty struct {
	DifficultyQuestionCount        int16
	DifficultyQuestionErrorCount   int16
	DifficultyQuestionCorrectCount int16
	DifficultyQuestionScore        float64
	DifficultyQuestionActualScore  float64
}

type TeahCheck struct {
	Value       int8  `json:"value"`
	QuestionId  int64 `json:"question_id"`
	KnowledgeId int64 `json:"knowledge_id"`
}

type ExamQuestionData struct {
	ExamQuestion   model.ExamQuestion
	QuestionTypeId int64
}

//----- submit() end -----

//----- Create() start -----

type Param struct {
	ExamNo              int64
	StudentUserNo       int64
	ScStudentUserNo     int64
	FrontUserCampusId   int64
	StudentUserCampusId int64
	Demand              map[int64]Demand
	QuestionDifficulty  map[int8]QuestionDifficulty
	Knowledge           map[int64]Knowledge
	QkData              []aiLogic.QkData
	QDdData             []aiLogic.QDdData
	AnsData             []aiLogic.AnsData
	TeahCheck           []TeahCheck
	Exam                model.Exam
}

//----- Create() end -----
