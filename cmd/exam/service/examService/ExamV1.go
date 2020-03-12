package examService

import (
	"errors"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"github.com/shopspring/decimal"
	"net/http"
	"reflect"
	"rxt/cmd/exam/dao/knowledgeDao"
	"rxt/cmd/exam/dao/questionDao"
	"rxt/cmd/exam/logic/aiLogic"
	"rxt/cmd/exam/logic/knowledgeLogic"
	"rxt/cmd/exam/logic/questionLogic"
	"rxt/cmd/exam/model"
	exam "rxt/cmd/exam/proto/sc"
	"rxt/cmd/exam/service/baseService"
	"rxt/cmd/exam/uitl"
	"rxt/internal/wrong"
	"strconv"
	"strings"
)

type V1 struct {
	baseService.Service
}

func (c V1) Submit(exam *exam.ExamRequest) (param Param, err error) {
	var examNo int64 = uitl.ExamNo()
	var questionNo []int64
	for _, examV := range exam.Exam {
		questionNo = append(questionNo, examV.QuestionNo)
	}
	demandMap := initDemand()                                            //报告知识点要求
	reportQuestionDifficultyMap := initQuestionDifficulty()              //题目难度
	questionDifficulty := questionDifficulty()                           //题目难度等级
	questionData := questionDao.New(c.Core).ExamQuestionByNo(questionNo) //考试题目查询 通过编号
	questionMap := map[int64]model.Question{}
	questionSmallMap := map[int64]model.QuestionSmall{}
	examQuestionTypeMap := map[int64]model.ExamQuestionType{}
	examQuestionMap := map[int64]ExamQuestionData{}
	knowledgeMap := map[int64]Knowledge{}
	examQuestionSmallMap := map[int64]model.ExamQuestionSmall{}
	examQuestionSmallKnowledgeMap := map[int64]model.ExamQuestionSmallKnowledge{}
	var examScore float64
	var examActualScore float64
	var questionNoArr []int64
	var examQuestionCount, examQuestionCorrectCount, examQuestionErrorCount int32
	var examPaperQuestionTypeOrder, examPaperQuestionOrder int16
	var qkDataArr []aiLogic.QkData   //传给AI 题目和知识点关联
	var qDdDataArr []aiLogic.QDdData //传给AI 题目难度系数
	var ansDataArr []aiLogic.AnsData //传给AI 学生答题结果
	for _, questionValue := range questionData {
		questionMap[questionValue.QuestionNo] = questionValue
		for _, questionSmallValue := range questionValue.QuestionSmall {
			questionSmallMap[questionSmallValue.QuestionSmallId] = questionSmallValue
		}
	}
	scStudentQuestionMap := map[int64]model.ScStudentQuestion{}
	if exam.ScStudentUserNo > 0 {
		scStudentQuestionMap = c.getScStudentQuestion(exam.ScStudentUserNo, questionNoArr)
	}

	for _, examValue := range exam.Exam {
		var examQuestionSmallOrder int32                       //考试小题编号
		var examQuestionScore, examQuestionActualScore float64 //考试题目分数、考试题目实际得分
		questionNoArr = append(questionNoArr, examValue.QuestionNo)
		questionIsRight := 1 //题目对错 1:对 2:错
		smallKnowledgeErrorCount := 0
		smallKnowledgeCorrectCount := 0

		examQuestionCount++
		examPaperQuestionTypeOrder++
		examPaperQuestionOrder++
		question := questionMap[examValue.QuestionNo] //题目数据
		for _, examQuestionSmallValue := range examValue.ExamQuestionSmall {
			var examQuestionSmallActualScore float64 //小题实际分数
			var examQuestionSmallIsRight int8
			examScore += examQuestionSmallValue.ExamQuestionSmallScore
			examQuestionSmallOrder++
			examQuestionScore += examQuestionSmallValue.ExamQuestionSmallScore
			questionSmall := questionSmallMap[examQuestionSmallValue.QuestionSmallId]
			var option []string
			for _, questionSmallAnswerOptionValue := range questionSmall.QuestionSmallAnswerOption {
				if questionSmallAnswerOptionValue.OptionName != "" {
					option = append(option, questionSmallAnswerOptionValue.OptionName)
				} else {
					option = append(option, questionSmallAnswerOptionValue.OptionContent)
				}
			}
			optionDiff := uitl.Difference(option, examQuestionSmallValue.ExamQuestionSmallAnswer)
			optionDiff1 := uitl.Difference(examQuestionSmallValue.ExamQuestionSmallAnswer, option)
			if reflect.DeepEqual(optionDiff, optionDiff1) || len(examQuestionSmallValue.ExamQuestionSmallAnswer) != len(option) {
				//错误处理
				questionIsRight = 0
				examQuestionSmallIsRight = 2
				smallKnowledgeErrorCount += len(questionSmall.QuestionSmallKnowledge)
			} else {
				examQuestionSmallIsRight = 1
				smallKnowledgeCorrectCount += len(questionSmall.QuestionSmallKnowledge)
				examQuestionActualScore += examQuestionSmallValue.ExamQuestionSmallScore
				examActualScore += examQuestionSmallValue.ExamQuestionSmallScore
				examQuestionSmallActualScore = examQuestionSmallValue.ExamQuestionSmallScore
			}
			examQuestionSmallAnswer := strings.Replace(strings.Trim(fmt.Sprint(examQuestionSmallValue.ExamQuestionSmallAnswer), "[]"), " ", ",", -1)
			examQuestionSmallMap[examQuestionSmallValue.QuestionSmallId] = model.ExamQuestionSmall{
				QuestionSmallId:              examQuestionSmallValue.QuestionSmallId,
				QuestionNo:                   examValue.QuestionNo,
				ExamQuestionSmallAnswer:      examQuestionSmallAnswer,
				ExamQuestionSmallScore:       examQuestionSmallValue.ExamQuestionSmallScore,
				ExamQuestionSmallActualScore: examQuestionSmallActualScore,
				ExamQuestionSmallIsRight:     examQuestionSmallIsRight,
				ExamQuestionSmallOrder:       examQuestionSmallOrder,
			}
			for _, QuestionSmallKnowledgeValue := range questionSmallMap[examQuestionSmallValue.QuestionSmallId].QuestionSmallKnowledge {
				examQuestionSmallKnowledgeMap[examQuestionSmallValue.QuestionSmallId] = model.ExamQuestionSmallKnowledge{
					QuestionSmallId:              examQuestionSmallValue.QuestionSmallId,
					KnowledgeNo:                  QuestionSmallKnowledgeValue.KnowledgeNo,
					ExamQuestionKnowledgeIsRight: examQuestionSmallIsRight,
				}
				qkDataArr = append(qkDataArr, aiLogic.QkData{
					KnowledgeId: QuestionSmallKnowledgeValue.Knowledge.KnowledgeId,
					QuestionId:  question.QuestionId,
					Value:       1,
				})
				var reportKnowledgeCorrectCount, reportKnowledgeErrorCount int16 = 1, 1
				if examQuestionSmallIsRight == 1 {
					reportKnowledgeErrorCount = 0
				} else {
					reportKnowledgeCorrectCount = 0
				}
				if _, ok := knowledgeMap[QuestionSmallKnowledgeValue.Knowledge.KnowledgeId]; ok {
					knowledge := knowledgeMap[QuestionSmallKnowledgeValue.Knowledge.KnowledgeId]
					knowledge.ReportKnowledgeExamCount += 1
					knowledge.ReportKnowledgeErrorCount += reportKnowledgeErrorCount
					knowledge.ReportKnowledgeCorrectCount += reportKnowledgeCorrectCount
					knowledge.ReportKnowledgeActualScore += examQuestionSmallActualScore
					knowledge.ReportKnowledgeScore += examQuestionSmallValue.ExamQuestionSmallScore
					knowledge.Question = append(knowledge.Question, Question{QuestionNo: question.QuestionNo, QuestionId: question.QuestionId})
				} else {
					knowledgeMap[QuestionSmallKnowledgeValue.Knowledge.KnowledgeId] = Knowledge{
						Knowledge: QuestionSmallKnowledgeValue.Knowledge,
						ReportKnowledge: model.ReportKnowledge{
							ExamNo:                      examNo,
							ReportKnowledgeCorrectCount: reportKnowledgeCorrectCount,
							ReportKnowledgeErrorCount:   reportKnowledgeErrorCount,
							ReportKnowledgeExamCount:    1,
							ReportKnowledgeActualScore:  examQuestionSmallActualScore,
							ReportKnowledgeScore:        examQuestionSmallValue.ExamQuestionSmallScore,
						},
						Question: []Question{{QuestionNo: question.QuestionNo, QuestionId: question.QuestionId}},
					}
				}
				if _, ok := demandMap[QuestionSmallKnowledgeValue.KnowledgeAttributeOne.KnowledgeDemandId].Question[examValue.QuestionNo]; !ok {
					demand := demandMap[QuestionSmallKnowledgeValue.KnowledgeAttributeOne.KnowledgeDemandId]
					demand.Question[examValue.QuestionNo] = examValue.QuestionNo
					demand.KnowledgeDemandExamCount++
					demand.KnowledgeDemandScore += examQuestionSmallValue.ExamQuestionSmallScore
					demand.KnowledgeDemandActualScore += examQuestionSmallValue.ExamQuestionSmallScore
					if examQuestionSmallIsRight == 2 {
						demand.KnowledgeDemandErrorCount++
					} else {
						demand.KnowledgeDemandCorrectCount++
					}
					demandMap[QuestionSmallKnowledgeValue.KnowledgeAttributeOne.KnowledgeDemandId] = demand
				}
			}
		}
		if questionIsRight == 1 {
			examQuestionCorrectCount++
		} else {
			examQuestionErrorCount++
		}
		if exam.ScStudentUserNo > 0 {
			questionLogic.New(c.Core).EditScStudentQuestion(scStudentQuestionMap, exam.StudentUserNo, examValue.QuestionNo, questionIsRight)
		}
		qDdDataArr = append(qDdDataArr, aiLogic.QDdData{
			QuestionId: question.QuestionId,
			Value:      questionDifficulty[questionMap[examValue.QuestionNo].QuestionDifficultyId],
		})
		ansData := aiLogic.AnsData{
			QuestionId: question.QuestionId,
			Value:      0,
		}
		if examQuestionActualScore > 0 {
			ansDataValueDecimal := decimal.NewFromFloat(examQuestionActualScore).
				Div(decimal.NewFromFloat(examQuestionScore)).
				Round(1)
			ansDataValue, _ := ansDataValueDecimal.Float64()
			ansData.Value = ansDataValue
		}
		ansDataArr = append(ansDataArr, ansData)
		if _, ok := examQuestionTypeMap[examValue.QuestionTypeId]; ok {
			examQuestionTypeUpdate := examQuestionTypeMap[examValue.QuestionTypeId]
			examQuestionTypeUpdate.ExamQuestionTypeCount++
			if questionIsRight == 1 {
				examQuestionTypeUpdate.ExamQuestionTypeCorrectCount++
				examQuestionTypeUpdate.ExamQuestionTypeActualScore += examQuestionScore
			} else {
				examQuestionTypeUpdate.ExamQuestionTypeErrorCount++
			}
			examQuestionTypeUpdate.ExamQuestionTypeScore += examQuestionScore
			examQuestionTypeMap[examValue.QuestionTypeId] = examQuestionTypeUpdate
		} else {
			examPaperQuestionOrder = 1
			examQuestionType := model.ExamQuestionType{
				ExamNo:                     examNo,
				QuestionTypeId:             examValue.QuestionTypeId,
				ExamPaperQuestionTypeOrder: examPaperQuestionTypeOrder,
				ExamQuestionTypeCount:      1,
				ExamQuestionTypeScore:      examQuestionScore,
				QuestionTypeCategoryId:     question.QuestionCategoryId,
			}
			if questionIsRight == 1 {
				examQuestionType.ExamQuestionTypeCorrectCount = 1
				examQuestionType.ExamQuestionTypeErrorCount = 0
				examQuestionType.ExamQuestionTypeActualScore = examQuestionScore
			} else {
				examQuestionType.ExamQuestionTypeErrorCount = 1
				examQuestionType.ExamQuestionTypeCorrectCount = 0

			}
			examQuestionTypeMap[examValue.QuestionTypeId] = examQuestionType
		}
		examQuestionMap[examValue.QuestionNo] = ExamQuestionData{
			ExamQuestion: model.ExamQuestion{QuestionNo: examValue.QuestionNo,
				ExamPaperQuestionOrder:  examPaperQuestionOrder,
				ExamQuestionScore:       examQuestionScore,
				ExamQuestionActualScore: examQuestionActualScore,
			},
			QuestionTypeId: examValue.QuestionTypeId,
		}
		reportQuestionDifficulty := reportQuestionDifficultyMap[question.QuestionDifficultyId]
		reportQuestionDifficulty.DifficultyQuestionCount++
		if questionIsRight == 1 {
			reportQuestionDifficulty.DifficultyQuestionCorrectCount++
		} else {
			reportQuestionDifficulty.DifficultyQuestionErrorCount++
		}
		reportQuestionDifficulty.DifficultyQuestionScore += examQuestionScore
		reportQuestionDifficulty.DifficultyQuestionActualScore += examQuestionActualScore
		reportQuestionDifficultyMap[question.QuestionDifficultyId] = reportQuestionDifficulty
	}

	examQuestionTypeId := map[int64]int64{}
	//写入考试题型
	for _, examQuestionTypeValue := range examQuestionTypeMap {
		c.Transaction.Create(&examQuestionTypeValue)
		examQuestionTypeId[examQuestionTypeValue.QuestionTypeId] = examQuestionTypeValue.ExamQuestionTypeId
	}
	//写入考试题目
	examQuestionId := map[int64]int64{}
	for _, examQuestionV := range examQuestionMap {
		examQuestionV.ExamQuestion.ExamQuestionTypeId = examQuestionTypeId[examQuestionV.QuestionTypeId]
		c.Transaction.Create(&examQuestionV.ExamQuestion)
		examQuestionId[examQuestionV.ExamQuestion.QuestionNo] = examQuestionV.ExamQuestion.ExamQuestionId
	}
	//写入考试小题
	examQuestionSmallKnowledge := map[int64]int64{}
	for _, examQuestionSmallV := range examQuestionSmallMap {
		examQuestionSmallV.ExamQuestionId = examQuestionId[examQuestionSmallV.QuestionNo]
		c.Transaction.Create(&examQuestionSmallV)
		examQuestionSmallKnowledge[examQuestionSmallV.QuestionSmallId] = examQuestionSmallV.ExamQuestionSmallId
	}
	for _, examQuestionSmallKnowledgeV := range examQuestionSmallKnowledgeMap {
		examQuestionSmallKnowledgeV.ExamQuestionSmallId = examQuestionSmallKnowledge[examQuestionSmallKnowledgeV.QuestionSmallId]
		c.Transaction.Create(&examQuestionSmallKnowledgeV)
	}
	examActualTime, _ := strconv.ParseFloat(exam.ExamActualTime, 32)
	examCreate := model.Exam{
		ExamName:                 exam.ExamName,
		ExamNo:                   examNo,
		ExamStatus:               3,
		ExamQuestionCount:        examQuestionCount,
		ExamQuestionCorrectCount: examQuestionCorrectCount,
		ExamQuestionErrorCount:   examQuestionErrorCount,
		SubjectId:                exam.SubjectId,
		GradeId:                  exam.GradeId,
		ExamActualTime:           float32(examActualTime),
		ExamAnswerType:           int8(exam.ExamAnswerType),
		ExamLevelType:            int8(exam.ExamLevelType),
		GradeChildrenId:          exam.GradeChildrenId,
		IsReview:                 int8(exam.IsReview),
		EditionNo:                exam.EditionNo,
		ScoringMethod:            int8(exam.ScoringMethod),
		ExamWordStatus:           int8(exam.ExamWordStatus),
		ExamTypeCode:             int(exam.ExamTypeCode),
		ExamScore:                examScore,
		ExamActualScore:          examActualScore,
	}
	c.Transaction.Create(&examCreate)
	param = Param{ExamNo: examNo, StudentUserNo: exam.StudentUserNo,
		ScStudentUserNo:   exam.ScStudentUserNo,
		FrontUserCampusId: 0, StudentUserCampusId: 0,
		Demand:             demandMap,
		QuestionDifficulty: reportQuestionDifficultyMap,
		Knowledge:          knowledgeMap,
		QkData:             qkDataArr,
		QDdData:            qDdDataArr,
		AnsData:            ansDataArr,
		TeahCheck:          nil,
		Exam:               examCreate,
	}
	return param, nil
}

// 报告知识点要求
func initDemand() map[int64]Demand {
	demand := map[int64]Demand{
		1: {
			Question:                    make(map[int64]int64),
			KnowledgeDemand:             1,
			KnowledgeDemandExamCount:    0,
			KnowledgeDemandErrorCount:   0,
			KnowledgeDemandCorrectCount: 0,
			KnowledgeDemandScore:        0,
		},
		2: {
			Question:                    make(map[int64]int64),
			KnowledgeDemand:             2,
			KnowledgeDemandExamCount:    0,
			KnowledgeDemandErrorCount:   0,
			KnowledgeDemandCorrectCount: 0,
			KnowledgeDemandScore:        0,
		},
		3: {
			Question:                    make(map[int64]int64),
			KnowledgeDemand:             3,
			KnowledgeDemandExamCount:    0,
			KnowledgeDemandErrorCount:   0,
			KnowledgeDemandCorrectCount: 0,
			KnowledgeDemandScore:        0,
		},
		4: {
			Question:                    make(map[int64]int64),
			KnowledgeDemand:             4,
			KnowledgeDemandExamCount:    0,
			KnowledgeDemandErrorCount:   0,
			KnowledgeDemandCorrectCount: 0,
			KnowledgeDemandScore:        0,
		},
		5: {
			Question:                    make(map[int64]int64),
			KnowledgeDemand:             5,
			KnowledgeDemandExamCount:    0,
			KnowledgeDemandErrorCount:   0,
			KnowledgeDemandCorrectCount: 0,
			KnowledgeDemandScore:        0,
		},
		6: {
			Question:                    make(map[int64]int64),
			KnowledgeDemand:             6,
			KnowledgeDemandExamCount:    0,
			KnowledgeDemandErrorCount:   0,
			KnowledgeDemandCorrectCount: 0,
			KnowledgeDemandScore:        0,
		},
	}
	return demand
}

//题目难度
func initQuestionDifficulty() map[int8]QuestionDifficulty {
	return map[int8]QuestionDifficulty{
		1: {DifficultyQuestionCount: 0,
			DifficultyQuestionErrorCount:   0,
			DifficultyQuestionCorrectCount: 0,
			DifficultyQuestionScore:        0,
			DifficultyQuestionActualScore:  0,
		},
		2: {DifficultyQuestionCount: 0,
			DifficultyQuestionErrorCount:   0,
			DifficultyQuestionCorrectCount: 0,
			DifficultyQuestionScore:        0,
			DifficultyQuestionActualScore:  0,
		},
		3: {DifficultyQuestionCount: 0,
			DifficultyQuestionErrorCount:   0,
			DifficultyQuestionCorrectCount: 0,
			DifficultyQuestionScore:        0,
			DifficultyQuestionActualScore:  0,
		},
		4: {DifficultyQuestionCount: 0,
			DifficultyQuestionErrorCount:   0,
			DifficultyQuestionCorrectCount: 0,
			DifficultyQuestionScore:        0,
			DifficultyQuestionActualScore:  0,
		},
		5: {DifficultyQuestionCount: 0,
			DifficultyQuestionErrorCount:   0,
			DifficultyQuestionCorrectCount: 0,
			DifficultyQuestionScore:        0,
			DifficultyQuestionActualScore:  0,
		},
	}
}

// 题目难度等级
func questionDifficulty() map[int8]float32 {
	return map[int8]float32{
		1: 1,
		2: 0.8,
		3: 0.6,
		4: 0.4,
		5: 0.2,
	}
}

// 学习路径
type course struct {
	minute        float64
	sort          int
	knowledgeList []editionKnowledgeNode
}

type editionKnowledgeTime struct {
	KnowledgeNo int64
	Time        float64
}
type editionKnowledgeNode struct {
	KnowledgeId   int64   `json:"knowledge_id"`
	KnowledgeNo   int64   `json:"knowledge_no"`
	KnowledgeSort int     `json:"knowledge_sort"`
	Time          float64 `json:"time"`
	PreNumber     int     `json:"pre_number"`
	AfterNumber   int     `json:"after_number"`
	Rate1         float64 `json:"rate1"`
	Rate2         float64 `json:"rate2"`
}

//报告创建
func (c V1) ReportCreate(param Param) (err error) {
	// 报告题目难度写入
	for QuestionDifficultyK, QuestionDifficultyV := range param.QuestionDifficulty {
		c.Transaction.Create(&model.ReportQuestionDifficulty{
			ExamNo:                         param.ExamNo,
			QuestionDifficultyId:           int64(QuestionDifficultyK),
			DifficultyQuestionCount:        QuestionDifficultyV.DifficultyQuestionCount,
			DifficultyQuestionErrorCount:   QuestionDifficultyV.DifficultyQuestionErrorCount,
			DifficultyQuestionCorrectCount: QuestionDifficultyV.DifficultyQuestionCorrectCount,
			DifficultyQuestionScore:        QuestionDifficultyV.DifficultyQuestionScore,
			DifficultyQuestionActualScore:  QuestionDifficultyV.DifficultyQuestionActualScore,
		})
	}
	//报告知识点要求写入
	for _, DemandV := range param.Demand {
		c.Transaction.Create(&model.ReportKnowledgeDemand{
			ExamNo:                      param.ExamNo,
			KnowledgeDemand:             DemandV.KnowledgeDemand,
			KnowledgeDemandExamCount:    DemandV.KnowledgeDemandExamCount,
			KnowledgeDemandErrorCount:   DemandV.KnowledgeDemandErrorCount,
			KnowledgeDemandCorrectCount: DemandV.KnowledgeDemandCorrectCount,
			KnowledgeDemandScore:        DemandV.KnowledgeDemandScore,
			KnowledgeDemandActualScore:  DemandV.KnowledgeDemandActualScore,
		})
	}
	reviewKnowledge, editionKnowledge, chapterKnowledgeSort := knowledgeDao.New(c.Core).
		EditionKnowledge(knowledgeDao.Param{
			EditionNo: param.Exam.EditionNo,
			SubjectId: param.Exam.SubjectId,
			IsReview:  param.Exam.IsReview,})
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"trace_id":     "0",
	})

	weakKnowledgeArr, DinaErr := aiLogic.New().Dina(param.QkData, param.AnsData, param.QDdData)
	if DinaErr != nil {
		return wrong.New(http.StatusBadRequest, DinaErr, "dina算法错误！")
	}
	var lpKnowledgeArr []aiLogic.LpKnowledge
	lpKnowledgeMap := map[int64]int64{}
	var knowledgeNoList []int64
	var knowledgeIdArr []int64
	var dArr []aiLogic.D
	mleDataMap := map[int64]aiLogic.MleData{}
	for _, weakKnowledgeV := range weakKnowledgeArr {
		mleDataMap[weakKnowledgeV.KnowledgeId] = weakKnowledgeV
		knowledgeIdArr = append(knowledgeIdArr, weakKnowledgeV.KnowledgeId)
	}
	scStudentKnowledgeMap := c.getScStudentKnowledge(param.StudentUserNo, knowledgeIdArr)
	for _, knowledgeV := range param.Knowledge {
		mleKnowledge := mleDataMap[knowledgeV.KnowledgeId]
		reportKnowledge := model.ReportKnowledge{
			ExamNo:                      param.ExamNo,
			KnowledgeNo:                 knowledgeV.Knowledge.KnowledgeNo,
			ReportKnowledgeIsWeak:       2,
			ReportKnowledgeProficiency:  mleKnowledge.Value,
			ReportKnowledgeDiff:         mleKnowledge.Diff,
			ReportKnowledgeExamCount:    knowledgeV.ReportKnowledgeExamCount,
			ReportKnowledgeErrorCount:   knowledgeV.ReportKnowledgeErrorCount,
			ReportKnowledgeCorrectCount: knowledgeV.ReportKnowledgeCorrectCount,
			ReportKnowledgeScore:        knowledgeV.ReportKnowledgeScore,
			ReportKnowledgeActualScore:  knowledgeV.ReportKnowledgeActualScore,
		}
		if mleKnowledge.Value < 0.5 {
			lpKnowledgeArr = append(lpKnowledgeArr, aiLogic.LpKnowledge{
				KnowledgeId: knowledgeV.Knowledge.KnowledgeId,
				Value:       mleKnowledge.Value,
				Label:       1,
			})
			lpKnowledgeMap[knowledgeV.Knowledge.KnowledgeNo] = knowledgeV.Knowledge.KnowledgeNo
			reportKnowledge.ReportKnowledgeIsWeak = 1
			dArr = append(dArr, aiLogic.D{
				KnowledgeId: knowledgeV.Knowledge.KnowledgeId,
				Value:       mleKnowledge.Value,
			})
		}

		if knowledgeV.Knowledge.HasChildren == 2 {
			knowledgeNoList = append(knowledgeNoList, knowledgeV.Knowledge.KnowledgeNo)
			c.Transaction.Create(&reportKnowledge)

			for _, QuestionV := range knowledgeV.Question {
				reportKnowledgeQuestion := model.ReportKnowledgeQuestion{
					ReportKnowledgeId: reportKnowledge.ReportKnowledgeId,
					QuestionNo:        QuestionV.QuestionNo,
					ExamQuestionId:    QuestionV.QuestionId,
				}
				c.Transaction.Create(&reportKnowledgeQuestion)
			}
		}
		knowledgeLogic.New(c.Core).ScStudentKnowledge(scStudentKnowledgeMap, mleKnowledge.Diff,
			mleKnowledge.Value, knowledgeV.Knowledge.KnowledgeNo, param.StudentUserNo)
	}

	var nodeData map[int64]int64
	nodeData = make(map[int64]int64)
	var n []map[string]int64
	if param.Exam.IsReview == 2 {
		for _, editionKnowledgesV := range editionKnowledge {
			nodeData[editionKnowledgesV.KnowledgeId] = editionKnowledgesV.KnowledgeId

			n = append(n, map[string]int64{"knowledge_id": editionKnowledgesV.KnowledgeId})
		}
	} else {
		for _, reviewKnowledgeV := range reviewKnowledge {
			nodeData[reviewKnowledgeV.Knowledge.KnowledgeId] = reviewKnowledgeV.Knowledge.KnowledgeId
			n = append(n, map[string]int64{"knowledge_id": reviewKnowledgeV.Knowledge.KnowledgeId})
		}
	}
	for weakKnowledgeK := 0; weakKnowledgeK < len(weakKnowledgeArr); weakKnowledgeK++ {
		if _, ok := nodeData[weakKnowledgeArr[weakKnowledgeK].KnowledgeId]; !ok {
			weakKnowledgeArr = append(weakKnowledgeArr[:weakKnowledgeK], weakKnowledgeArr[weakKnowledgeK+1:]...)
			weakKnowledgeK--
		}
	}

	//潜在薄弱知识点
	preKnowledgeData := knowledgeDao.New(c.Core).ListBySubjectId(param.Exam.SubjectId, param.Exam.IsReview, param.Exam.EditionNo)
	var preKnowledgeMap []map[string]int64
	for _, preKnowledgeDataV := range preKnowledgeData {
		preKnowledgeMap = append(preKnowledgeMap, map[string]int64{
			"knowledge_id":     preKnowledgeDataV.Knowledge.KnowledgeId,
			"pre_knowledge_id": preKnowledgeDataV.Knowledge.KnowledgeParentId,
		})
	}
	lk, LkAErr := aiLogic.New().LkA(n, preKnowledgeMap, dArr)
	if LkAErr != nil {
		return wrong.New(http.StatusBadRequest, LkAErr, "潜在薄弱知识点算法错误！")
	}
	var weekKnowledgeNoLearMax int
	for _, knowledgeNoListV := range knowledgeNoList {
		if _, ok := chapterKnowledgeSort[knowledgeNoListV]; ok && weekKnowledgeNoLearMax < chapterKnowledgeSort[knowledgeNoListV] {
			weekKnowledgeNoLearMax = chapterKnowledgeSort[knowledgeNoListV]
		}
	}
	var lkIdArr []int64
	for _, lkValue := range lk {
		lkIdArr = append(lkIdArr, lkValue.KnowledgeId)
	}
	var knowledge []model.Knowledge
	c.Db.Where("knowledge_id IN (?) ", lkIdArr).Find(&knowledge)
	var chapterKnowledgeData map[int64]model.Knowledge
	chapterKnowledgeData = make(map[int64]model.Knowledge)
	for _, knowledgeV := range knowledge {
		chapterKnowledgeData[knowledgeV.KnowledgeId] = knowledgeV
	}
	for lkKey := 0; lkKey < len(lk); lkKey++ {
		lkV := lk[lkKey]
		chapterKnowledgeV := chapterKnowledgeData[lkV.KnowledgeId]
		mleDataV := mleDataMap[lkV.KnowledgeId]
		if _, mleData := mleDataMap[lkV.KnowledgeId]; mleData {
			if mleDataV.Value >= 0.5 {
				lk = append(lk[:lkKey], lk[lkKey+1:]...)
				lkKey--
			}
		}
		if _, ok := chapterKnowledgeSort[chapterKnowledgeV.KnowledgeNo]; ok || chapterKnowledgeSort[chapterKnowledgeV.KnowledgeNo] > weekKnowledgeNoLearMax &&
			param.Exam.IsReview == 2 {
			lk = append(lk[:lkKey], lk[lkKey+1:]...)
			lkKey--
			continue
		}
		lpKnowledgeArr = append(lpKnowledgeArr, aiLogic.LpKnowledge{
			KnowledgeId: chapterKnowledgeV.KnowledgeId,
			Value:       mleDataV.Value,
			Label:       1,
		})
		lpKnowledgeMap[chapterKnowledgeV.KnowledgeNo] = chapterKnowledgeV.KnowledgeNo
	}
	for _, lkV := range lk {
		chapterKnowledgeV := chapterKnowledgeData[lkV.KnowledgeId]
		c.Transaction.Create(model.ReportLatentKnowledge{
			ExamNo:                           param.ExamNo,
			KnowledgeId:                      lkV.KnowledgeId,
			KnowledgeNo:                      chapterKnowledgeV.KnowledgeNo,
			ReportLatentKnowledgeProficiency: lkV.Value,
			RelateProb:                       lkV.RelateProb,
			ReportLatentKnowledgeType:        1,
			ReportLatentKnowledgeIsExceed:    1,
		})
	}
	paperAnalysis := getPaperAnalysisByDistrictAndType(param.Exam.SubjectId,
		0,
		0,
		0,
		param.Exam.IsReview)
	var paperAnalysisKnowledgeArr []model.PaperAnalysisKnowledge
	paperAnalysisKnowledgeMap := map[int64]map[int]model.PaperAnalysisKnowledge{}

	if paperAnalysis.PaperAnalysisNo > 0 {
		c.Db.Where("paper_analysis_no = ?", paperAnalysis.PaperAnalysisNo).Find(&paperAnalysisKnowledgeArr)
		for _, paperAnalysisKnowledgeV := range paperAnalysisKnowledgeArr {
			paperAnalysisKnowledgeMap[paperAnalysisKnowledgeV.KnowledgeNo][paperAnalysisKnowledgeV.PaperTypeId] = paperAnalysisKnowledgeV
		}
	}

	var knowledgePreAfterCountArr []model.KnowledgePreAfterCount
	knowledgePreAfterCountModel := c.Db
	if param.Exam.IsReview == 1 {
		knowledgePreAfterCountModel.Where("is_revision = 1")
	}
	if param.Exam.IsReview == 2 {
		knowledgePreAfterCountModel.Where("is_revision = 2")
	}
	knowledgePreAfterCountModel.Find(&knowledgePreAfterCountArr)

	knowledgePreAfterCountMap := map[int64]model.KnowledgePreAfterCount{}
	for _, knowledgePreAfterCountV := range knowledgePreAfterCountArr {
		knowledgePreAfterCountMap[knowledgePreAfterCountV.KnowledgeNo] = knowledgePreAfterCountV
	}
	var sumTime float64
	var editionKnowledgeIds []int64
	editionKnowledgeTimeMap := map[int64]editionKnowledgeTime{}
	editionKnowledgeNodeMap := map[int64]editionKnowledgeNode{}
	var editionKnowledgeNodeArr []aiLogic.EditionKnowledgeNode

	for editionKnowledgeK, editionKnowledgeV := range editionKnowledge {
		if editionKnowledgeV.KnowledgeId == 0 {
			continue
		}
		editionKnowledgeIds = append(editionKnowledgeIds, editionKnowledgeV.KnowledgeId)
		var time float64 = editionKnowledgeV.KnowledgeTeachingMinute
		if len(editionKnowledgeV.TeachingPlanQuestionMany) > 0 {
			var exampleMinuteSum float64
			for _, TeachingPlanQuestionV := range editionKnowledgeV.TeachingPlanQuestionMany {
				exampleMinuteSum += TeachingPlanQuestionV.ExampleMinute
			}
			time = exampleMinuteSum
		}
		editionKnowledgeTimeMap[editionKnowledgeV.KnowledgeId] = editionKnowledgeTime{
			KnowledgeNo: editionKnowledgeV.KnowledgeNo,
			Time:        time,
		}
		editionKnowledgeNode := editionKnowledgeNode{
			KnowledgeId:   editionKnowledgeV.KnowledgeId,
			KnowledgeNo:   editionKnowledgeV.KnowledgeNo,
			KnowledgeSort: editionKnowledgeK,
			Time:          time,
			PreNumber:     knowledgePreAfterCountMap[editionKnowledgeV.KnowledgeNo].KnowledgePreCount,
			AfterNumber:   knowledgePreAfterCountMap[editionKnowledgeV.KnowledgeNo].KnowledgeAfterCount,
		}
		if param.Exam.IsReview == 1 {
			if _, ok := paperAnalysisKnowledgeMap[editionKnowledgeV.KnowledgeId][1]; ok {
				rate1 := decimal.NewFromInt(int64(paperAnalysisKnowledgeMap[editionKnowledgeV.KnowledgeId][1].Score)).Div(decimal.NewFromInt(1000)).Round(2)
				editionKnowledgeNode.Rate1, _ = rate1.Float64()
			}
			if _, ok := paperAnalysisKnowledgeMap[editionKnowledgeV.KnowledgeId][2]; ok {
				rate2 := decimal.NewFromInt(int64(paperAnalysisKnowledgeMap[editionKnowledgeV.KnowledgeId][2].Score)).Div(decimal.NewFromInt(1000)).Round(2)
				editionKnowledgeNode.Rate2, _ = rate2.Float64()
			}

		} else if param.Exam.IsReview == 2 {
			if _, ok := paperAnalysisKnowledgeMap[editionKnowledgeV.KnowledgeId][4]; ok {
				rate1 := decimal.NewFromInt(int64(paperAnalysisKnowledgeMap[editionKnowledgeV.KnowledgeId][4].Score)).Div(decimal.NewFromInt(1000)).Round(2)
				editionKnowledgeNode.Rate2, _ = rate1.Float64()
			}
			if _, ok := paperAnalysisKnowledgeMap[editionKnowledgeV.KnowledgeId][5]; ok {
				rate2 := decimal.NewFromInt(int64(paperAnalysisKnowledgeMap[editionKnowledgeV.KnowledgeId][5].Score)).Div(decimal.NewFromInt(1000)).Round(2)
				editionKnowledgeNode.Rate2, _ = rate2.Float64()
			}
		}
		editionKnowledgeNodeMap[editionKnowledgeV.KnowledgeId] = editionKnowledgeNode
		editionKnowledgeNodeArr = append(editionKnowledgeNodeArr, aiLogic.EditionKnowledgeNode{
			KnowledgeId:   editionKnowledgeNode.KnowledgeId,
			KnowledgeSort: editionKnowledgeNode.KnowledgeSort,
			Time:          editionKnowledgeNode.Time,
			PreNumber:     editionKnowledgeNode.PreNumber,
			AfterNumber:   editionKnowledgeNode.AfterNumber,
			Rate1:         editionKnowledgeNode.Rate1,
			Rate2:         editionKnowledgeNode.Rate2,
		})
		if _, ok := lpKnowledgeMap[editionKnowledgeV.KnowledgeNo]; ok {
			sumTime += time
		}
	}
	lp, LpErr := aiLogic.New().LP(lpKnowledgeArr, editionKnowledgeNodeArr, sumTime)
	if LpErr != nil {
		return wrong.New(http.StatusBadRequest, LpErr, "学习路径算法错误！")
	}
	var i int = 1
	courseMap := map[int]course{}
	var main float64
	for _, lpV := range lp {
		var knowledgeMinute float64
		if _, ok := editionKnowledgeNodeMap[lpV]; ok && editionKnowledgeNodeMap[lpV].Time > 0 {
			knowledgeMinute += editionKnowledgeNodeMap[lpV].Time
		}
		if knowledgeMinute > 0 {
			main += knowledgeMinute
			course := course{
				minute: 120,
				sort:   i,
			}
			course.knowledgeList = append(course.knowledgeList, editionKnowledgeNodeMap[lpV])
			courseMap[i] = course
		}
		if main >= 100 {
			main = 0
			i++
		}
	}
	for _, courseV := range courseMap {
		reportCourse := model.ReportCourse{
			ExamNo:           param.ExamNo,
			ReportCourseSort: courseV.sort,
		}
		c.Transaction.Create(&reportCourse)
		for _, courseKnowledgeV := range courseV.knowledgeList {
			reportCourseKnowledge := model.ReportCourseKnowledge{
				ReportCourseId: reportCourse.ReportCourseId,
				KnowledgeNo:    courseKnowledgeV.KnowledgeNo,
			}
			c.Transaction.Create(&reportCourseKnowledge)
		}
	}
	return nil
}

func getPaperAnalysisByDistrictAndType(subjectId int64, provinceId int64, cityId int64, districtId int64, paperAnalysisType int8) model.PaperAnalysis {

	param := map[string]interface{}{
		"district_id":         districtId,
		"paper_analysis_type": paperAnalysisType,
		"subject_id":          subjectId,
	}
	PaperAnalysis, recordNotFound := knowledgeDao.New().GetPaperAnalysis(param)
	if !recordNotFound {
		return PaperAnalysis
	}

	PaperAnalysis2, recordNotFound2 := knowledgeDao.New().GetPaperAnalysis(map[string]interface{}{
		"city_id":             cityId,
		"district_id":         0,
		"paper_analysis_type": paperAnalysisType,
		"subject_id":          subjectId,
	})
	if !recordNotFound2 {
		return PaperAnalysis2
	}
	PaperAnalysis3, recordNotFound3 := knowledgeDao.New().GetPaperAnalysis(map[string]interface{}{
		"city_id":             0,
		"province_id":         provinceId,
		"district_id":         0,
		"paper_analysis_type": paperAnalysisType,
		"subject_id":          subjectId,
	})
	if !recordNotFound3 {
		return PaperAnalysis3
	}
	PaperAnalysis4, _ := knowledgeDao.New().GetPaperAnalysis(map[string]interface{}{
		"district_id":         0,
		"province_id":         0,
		"city_id":             0,
		"paper_analysis_type": paperAnalysisType,
		"subject_id":          subjectId,
	})
	return PaperAnalysis4

}

func (c V1) getScStudentKnowledge(studentUserNo int64, knowledgeIdArr []int64) map[int64]model.ScStudentKnowledge {
	var scStudentKnowledgeArr []model.ScStudentKnowledge
	scStudentKnowledgeMap := map[int64]model.ScStudentKnowledge{}
	c.Db.Select("rxt_sc_student_knowledge.*").
		Where("rxt_sc_student_knowledge.student_user_no = ?", studentUserNo).
		Where("rxt_knowledge.knowledge_id IN(?)", knowledgeIdArr).
		Joins("rxt_knowledge ON rxt_knowledge.knowledge_no = rxt_sc_student_knowledge.knowledge_no").
		Find(&scStudentKnowledgeArr)
	for _, scStudentKnowledgeArrV := range scStudentKnowledgeArr {
		scStudentKnowledgeMap[scStudentKnowledgeArrV.KnowledgeNo] = scStudentKnowledgeArrV
	}
	return scStudentKnowledgeMap
}

func (c V1) getScStudentQuestion(studentUserNo int64, questionNoArr []int64) map[int64]model.ScStudentQuestion {
	var scStudentQuestionArr []model.ScStudentQuestion
	scStudentQuestionMap := map[int64]model.ScStudentQuestion{}
	c.Db.Where("question_no IN(?)", questionNoArr).
		Where("student_user_no = ", studentUserNo).
		Find(&scStudentQuestionArr)
	for _, scStudentQuestionArrV := range scStudentQuestionArr {
		scStudentQuestionMap[scStudentQuestionArrV.QuestionNo] = scStudentQuestionArrV
	}
	return scStudentQuestionMap
}

func (c V1) CreateScExamStudent(examNo int64, studentUserNo int64, bookNo int64) error {
	if studentUserNo == 0 || examNo == 0 {
		wrong.New(http.StatusBadRequest, errors.New("学生或考试编号丢失！"))
	}
	c.Transaction.Create(&model.ScExamStudent{
		StudentUserNo: studentUserNo,
		ExamNo:        examNo,
		BookNo:        bookNo,
	})
	return nil
}
