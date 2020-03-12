package examService

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kirinlabs/HttpRequest"
	"rxt/cmd/exam/dao/knowledgeDao"
	"rxt/cmd/exam/logic/aiLogic"
	"rxt/internal/core"
	"rxt/internal/core/register"
	"testing"
)

func TestGet(t *testing.T) {
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{"Content-Type": "application/json"})
	res, err := req.Post("https://oa-api.517rxt.com/v1/test")
	var m map[string]interface{}
	if err == nil {
		res.Json(&m)
		fmt.Println(m)
	} else {
		fmt.Print(err)
	}
}

func TestMle(t *testing.T) {
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZGVudGl0eSI6InJ4dCIsImlhdCI6MTU4MjYxNTg0MywiZnJlc2giOmZhbHNlLCJqdGkiOiI2OGRiYTI0ZS05NDVlLTRlNDEtOGNhOS01N2VmMmFmNzg1NzciLCJ0eXBlIjoiYWNjZXNzIiwibmJmIjoxNTgyNjE1ODQzLCJleHAiOjE1ODI2MTU5MDN9.bdVbI28TBVTjLlSZ3Nc_Nak11UdY517xc4Ozam7uH58",
	})

	qkData := []aiLogic.QkData{}
	qkData = append(qkData, aiLogic.QkData{
		Value:       1,
		QuestionId:  1068191,
		KnowledgeId: 190,
	})
	qkData = append(qkData, aiLogic.QkData{
		Value:       1,
		QuestionId:  1068191,
		KnowledgeId: 191,
	})
	qkData = append(qkData, aiLogic.QkData{
		Value:       1,
		QuestionId:  1038189,
		KnowledgeId: 182,
	})
	qkData = append(qkData, aiLogic.QkData{
		Value:       1,
		QuestionId:  1069074,
		KnowledgeId: 182,
	})
	qkData = append(qkData, aiLogic.QkData{
		Value:       1,
		QuestionId:  1068180,
		KnowledgeId: 188,
	})
	qkData = append(qkData, aiLogic.QkData{
		Value:       1,
		QuestionId:  1068180,
		KnowledgeId: 190,
	})
	qkData = append(qkData, aiLogic.QkData{
		Value:       1,
		QuestionId:  1068180,
		KnowledgeId: 191,
	})
	qkData = append(qkData, aiLogic.QkData{
		Value:       1,
		QuestionId:  1038192,
		KnowledgeId: 182,
	})

	test := map[string]interface{}{
		"-qk_json_string": qkData,
		"-ans_json_string": []map[string]interface{}{
			{
				"value":       "1.0",
				"question_id": 1068191,
			},
			{
				"value":       "0.4",
				"question_id": 1038189,
			},
			{
				"value":       "0.7",
				"question_id": 1069074,
			},
			{
				"value":       "0.7",
				"question_id": 1068180,
			},
			{
				"value":       "0.4",
				"question_id": 1038192,
			},
		},
		"-q_dd_json_string": []map[string]interface{}{
			{
				"value":       "0.4",
				"question_id": 1068191,
			},
			{
				"value":       "0.6",
				"question_id": 1038189,
			},
			{
				"value":       "0.6",
				"question_id": 1069074,
			},
			{
				"value":       "0.4",
				"question_id": 1068180,
			},
			{
				"value":       "0.4",
				"question_id": 1038192,
			},
		},
		"-sk_json_string": map[string]interface{}{},
		"-teah_json_string": []map[string]interface{}{
			{"value": 1, "question_id": 1068191, "knowledge_id": 190},
			{"value": 1, "question_id": 1068191, "knowledge_id": 191},
			{"value": -1, "question_id": 1038189, "knowledge_id": 182},
			{"value": -1, "question_id": 1069074, "knowledge_id": 182},
			{"value": -1, "question_id": 1068180, "knowledge_id": 188},
			{"value": -1, "question_id": 1068180, "knowledge_id": 190},
			{"value": 1, "question_id": 1068180, "knowledge_id": 191},
			{"value": -1, "question_id": 1038192, "knowledge_id": 182},
		},
	}
	qkJson, _ := json.Marshal(test)
	fmt.Println(string(qkJson))

	res, err := req.JSON().Post("http://101.132.145.239:19511/DINA",
		test)
	fmt.Print(err)
	var m []map[string]interface{}
	if err == nil {
		b, _ := res.Body()
		fmt.Print(string(b))
		res.Json(&m)
		fmt.Print(11111111111111, m)
	} else {
		fmt.Print(err)
	}
}

func TestQ(t *testing.T) {
	new(register.Register).ConfigRegister().
		RedisRegister().
		DbRegister().
		FacadeCacheRegister().
		SetPort(5001)
	c := core.New()
	var chapterKnowledge []knowledgeDao.ChapterKnowledge
	c.Db.
		Select("rxt_chapter_knowledge.knowledge_no,rxt_chapter_knowledge.chapter_no,chapter_knowledge_sort,chapter_knowledge_id,"+
			"rxt_knowledge_info.knowledge_teaching_minute,"+
			"rxt_knowledge.knowledge_id,rxt_knowledge.knowledge_parent_id,rxt_knowledge.knowledge_name").
		Joins("INNER JOIN rxt_knowledge ON rxt_knowledge.knowledge_no = rxt_chapter_knowledge.knowledge_no").
		Joins("LEFT JOIN rxt_knowledge_info ON rxt_knowledge.knowledge_no = rxt_knowledge_info.knowledge_no").
		Preload("TeachingPlanQuestionMany", func(db *gorm.DB) *gorm.DB {
			return db.Where("teaching_plan_question_type = ?", 1).
				Where("teaching_plan_question_level_type = ?", 1)
		}).
		Where("rxt_chapter_knowledge.edition_no = ?", 10000061).
		Order("chapter_knowledge_sort asc").
		Find(&chapterKnowledge)

	JsonStr(chapterKnowledge)
}

func TestEditionChapterTree(t *testing.T) {
	new(register.Register).ConfigRegister().
		RedisRegister().
		DbRegister().
		FacadeCacheRegister().
		SetPort(5001)
	bookChapter := knowledgeDao.New().EditionChapterTree(10000061, 0, 0)
	JsonStr(bookChapter)
}
func JsonStr(data interface{}) {
	json, _ := json.Marshal(data)
	fmt.Println(string(json))
	fmt.Println("")
}

func TestD(t *testing.T) {
	allCate := []cate{
		cate{1, "a", 0},
		cate{2, "b", 0},
		cate{3, "c", 0},
		cate{4, "aa", 1},
		cate{5, "bb", 2},
		cate{6, "cc", 3},
		cate{7, "aaa", 4},
		cate{8, "bbb", 5},
		cate{9, "ccc", 6},
	}

	//实现无限级分类

	//arr := superCategory(allCate, 0)
	arr := iteration(allCate, 0)

	for _, value := range arr {
		fmt.Println(value)
	}

}

type cate struct {
	id   int
	name string
	pid  int
	//son  *cate
}
type cateTree struct {
	id   int
	name string
	pid  int
	son  []cateTree
}

//迭代实现
func iteration(allCate []cate, pid int) []cate {
	task := []int{}
	task = append(task, 0)
	res := []cate{}
	hasChild := false
	var parent int = pid
	for len(task) > 0 {
		hasChild = false
		for k, v := range allCate {
			if v.id == -1 {
				continue
			}
			if parent == v.pid {
				res = append(res, v)
				task = append(task, v.id)
				allCate[k].id = -1 //奖该数据删除
				hasChild = true
				parent = v.id
				break
			}
		}

		if !hasChild {
			end := len(task) - 1
			task = task[0:end] //将该数据删除
			//继续找它上级得其他子类
			if len(task) > 0 {
				end = end - 1
				parent = task[end]
			}

		}

	}
	return res
}

func TestNew(t *testing.T) {
	var a int
	a++
	fmt.Println(a)
}
