package student

import (
	"context"
	pb "rxt/cmd/ads/proto"
	"rxt/cmd/ads/service/ads"
)

type Server struct{}

func (server *Server) Find(ctx context.Context, request *pb.FindRequest) (*pb.FindResponse, error) {
	response := &pb.FindResponse{}
	res, err := ads.New().Find(request)
	if err != nil {
		return nil, err
	}
	for _, v := range res {
		response.Result = append(response.Result, &pb.Advertisement{
			AdvertisementId:                    v.AdvertisementId,
			AdvertisementNo:                    v.AdvertisementNo,
			AdvertisementApplicationPositionId: v.AdvertisementApplicationId,
			AdvertisementIndex:                 v.AdvertisementIndex,
			LinkType:                           v.LinkType,
			AdvertisementImgUrl:                v.AdvertisementImgUrl,
			AdvertisementLinkUrl:               v.AdvertisementLinkUrl,
			ApplicationUrl:                     v.ApplicationUrl,
			Status:                             v.Status,
			CreatedAt:                          v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return response, nil
}
