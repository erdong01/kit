package ads

import (
	"rxt/cmd/ads/model"
	ads "rxt/cmd/ads/proto"
)

type Ads struct {
	IAds
}

func New() Ads {
	service := &Service{}
	service.Init()
	return Ads{service}
}

type IAds interface {
	Find(request *ads.FindRequest) ([]model.Advertisement, error)
}
