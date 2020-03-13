package ads

import (
	"rxt/cmd/ads/model"
	"rxt/cmd/ads/proto"
	"rxt/cmd/ads/service/base"
	"rxt/internal/log"
)

type Service struct {
	base.Service
}

func (service *Service) Find(request *ads.FindRequest) (list []model.Advertisement, err error) {
	db := service.Db
	// 开启sql日志
	log.SetSqlLogger(db)
	selected := []string{
		"advertisement_id",
		"advertisement_no",
		"advertisement_application_position_id",
		"advertisement_index",
		"link_type",
		"advertisement_img_url",
		"advertisement_link_url",
		"application_url",
		"status",
		"created_at",
	}
	db.Select(selected).Where("advertisement_application_position_id = ? and status = ?", request.PositionId, 2).Find(&list)

	return
}
