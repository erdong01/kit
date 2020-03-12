package skip

import (
	"errors"
	"net/http"
	"rxt/cmd/barrier/model"
	pb "rxt/cmd/barrier/proto/sc"
	"rxt/cmd/barrier/service/base"
	"rxt/internal/wrong"

	"github.com/jinzhu/gorm"
)

type ScSkipV1 struct {
	base.Service
}

// Skip 跳关
func (sc *ScSkipV1) Skip(request *pb.Request) (bool, error) {
	ScBarrierGame := model.ScBarrierGame{}

	sc.Db.Find(&ScBarrierGame, request.BarrierGameId)
	sc.Db.Model(&ScBarrierGame).Order("barrier_knowledge_order").Related(&ScBarrierGame.KnowledgeList, "barrier_game_id")

	var key int
	list := make(map[int]model.ScBarrierGameKnowledge)
	for _, v := range ScBarrierGame.KnowledgeList {
		list[int(v.BarrierKnowledgeOrder)] = v
		if v.KnowledgeNo == request.KnowledgeNo {
			key = int(v.BarrierKnowledgeOrder)
		}
	}

	if err := unlockBarrierKnowledge(sc.Db, list, key); err != nil {
		return false, wrong.New(http.StatusExpectationFailed, err, "")
	}

	if ScBarrierGame.BarrierGameStatus != 2 && getBarrierGameStatus(&ScBarrierGame) {
		sc.Db.Model(&ScBarrierGame).Update("barrier_game_status", 2)
	}

	return true, nil
}

func unlockBarrierKnowledge(db *gorm.DB, list map[int]model.ScBarrierGameKnowledge, key int) error {
	if _, exist := list[key]; !exist {
		return wrong.New(http.StatusExpectationFailed, errors.New("找不到知识点"), "找不到知识点")
	}

	if list[key].BarrierGameKnowledgeStatus == 3 || list[key].BarrierGameKnowledgeStatus == 4 {
		return wrong.New(http.StatusExpectationFailed, errors.New("当前关卡不允许跳过"), "当前关卡不允许跳过")
	}

	if err := db.Model(list[key]).Update("barrier_game_knowledge_status", 3).Error; err != nil {
		return wrong.New(http.StatusInternalServerError, err, "更新关卡知识点失败")
	}

	// 查找下一关 没有下一关直接返回
	if _, nextExist := list[key+1]; !nextExist {
		return nil
	}

	// 查找上一关 无上一关或者上一关已跳过或结束 则开始解锁下两关
	if _, lastExist := list[key-1]; !lastExist || list[key-1].BarrierGameKnowledgeStatus == 3 || list[key-1].BarrierGameKnowledgeStatus == 4 {
		var target int
		if list[key+1].BarrierGameKnowledgeStatus == 3 || list[key+1].BarrierGameKnowledgeStatus == 4 {
			target = key + 3
		} else {
			target = key + 1
		}
		// 解锁下两关 如果下一关已解锁 则顺延一关(即解锁下下下关)
		if err := db.Model(&model.ScBarrierGameKnowledge{}).Where("barrier_game_id = ? and barrier_game_knowledge_status = ? and (barrier_knowledge_order = ? or barrier_knowledge_order = ?) ", list[key].BarrierGameID, 1, key+2, target).Update("barrier_game_knowledge_status", 2).Error; err != nil {
			return wrong.New(http.StatusInternalServerError, err, "")
		}
	}

	return nil
}

func getBarrierGameStatus(barrier *model.ScBarrierGame) bool {
	for _, v := range barrier.KnowledgeList {
		if v.BarrierGameKnowledgeStatus != 3 && v.BarrierGameKnowledgeStatus != 4 {
			return false
		}
	}

	return true
}
