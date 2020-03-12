package aiLogic

import "rxt/internal/core"

type Ai struct {
	IAi
}

func New(param ...*core.Core) Ai {
	res := &V1{}
	res.Init(param...)
	return Ai{res}
}

type IAi interface {
	Dina(QkData []QkData, AnsData []AnsData, QDdData []QDdData) ([]MleData, error)
	LkA(n []map[string]int64, preKnowledgeMap []map[string]int64, dArr []D) ([]Lk, error)
	LP(lpKnowledgeMap []LpKnowledge, editionKnowledgeNodeArr []EditionKnowledgeNode,time float64) ([]int64, error)
}

//----- Dina() start -----
// ai 返回薄弱知识点结果
type MleData struct {
	KnowledgeId int64   `json:"knowledge_id"`
	Value       float64 `json:"value"`
	Diff        float64 `json:"diff"`
}
type QkData struct {
	KnowledgeId int64 `json:"knowledge_id"`
	QuestionId  int64 `json:"question_id"`
	Value       int8  `json:"value"`
}
type QDdData struct {
	QuestionId int64   `json:"question_id"`
	Value      float32 `json:"value"`
}

type AnsData struct {
	QuestionId int64   `json:"question_id"`
	Value      float64 `json:"value"`
}

//----- Dina() end -----
type D struct {
	KnowledgeId int64   `json:"knowledge_id"`
	Value       float64 `json:"value"`
}
type Lk struct {
	KnowledgeId int64   `json:"knowledge_id"`
	Value       float64 `json:"value"`
	RelateProb  float64 `json:"relate_prob"`
}

type LpKnowledge struct {
	KnowledgeId int64   `json:"knowledge_id"`
	Value       float64 `json:"value"`
	Label       int8    `json:"label"`
}
type EditionKnowledgeNode struct {
	KnowledgeId   int64   `json:"knowledge_id"`
	KnowledgeSort int     `json:"knowledge_sort"`
	Time          float64 `json:"time"`
	PreNumber     int     `json:"pre_number"`
	AfterNumber   int     `json:"after_number"`
	Rate1         float64 `json:"rate1"`
	Rate2         float64 `json:"rate2"`
}
