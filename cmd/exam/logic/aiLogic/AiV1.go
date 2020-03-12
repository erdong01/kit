package aiLogic

import (
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"rxt/cmd/exam/logic/base"
)

type V1 struct {
	base.Logic
}

func (c V1) Dina(QkData []QkData, AnsData []AnsData, QDdData []QDdData) ([]MleData, error) {
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"trace_id":     "0",
	})
	dinaPost := map[string]interface{}{
		"-qk_json_string":   QkData,                   //为本次测试题目和知识点关联矩阵，用于计算Q矩阵
		"-ans_json_string":  AnsData,                  //为学生答题的情况
		"-q_dd_json_string": QDdData,                  //为题目的难度系数
		"-sk_json_string":   map[string]interface{}{}, //为学生的知识点掌握历史情况，（为空，后期删除）
		"-teah_json_string": []int64{},                //为教师勾选薄弱知识点，每道题目中，被勾选的知识点value为-1 未被勾选知识点为 1，如果只勾选书写错误，则认为知识点都掌握，value为1
	}
	res, err := req.JSON().Post("http://101.132.145.239:19511/DINA", dinaPost)
	var weakKnowledgeArr []MleData
	if err == nil {
		res.Json(&weakKnowledgeArr)
	}
	return weakKnowledgeArr, err
}

func (c V1) LkA(n []map[string]int64, preKnowledgeMap []map[string]int64, dArr []D) ([]Lk, error) {
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"trace_id":     "0",
	})
	lkaPost := map[string]interface{}{
		"-n": n,               //知识图谱中的知识点id列表
		"-p": preKnowledgeMap, //前置知识点图谱
		"-l": []int{},         //潜在知识点掌握情况，初始为空
		"-d": dArr,            //模型给出的知识点掌握度
	}
	var lk []Lk
	lkRes, err := req.Post("http://101.132.145.239:19511/LKA",
		lkaPost)
	if err == nil {
		lkRes.Json(&lk)
		fmt.Println(lk)
	}
	return lk, err
}

func (c V1) LP(lpKnowledgeMap []LpKnowledge, editionKnowledgeNodeArr []EditionKnowledgeNode, time float64) ([]int64, error) {
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"trace_id":     "0",
	})
	var lp []int64
	lpMap := map[string]interface{}{
		"--sk_data":   lpKnowledgeMap,
		"--node_data": editionKnowledgeNodeArr,
		"--pre_know":  []int{},
		"--time":      time,
		"--rate_flag": 1,
	}
	lpRes, err := req.Post("http://101.132.145.239:19511/LP_time", lpMap)
	if err == nil {
		lpRes.Json(&lp)
		fmt.Print(lp)
	}
	return lp, err
}
