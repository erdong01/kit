package adsClient

import (
	"rxt/cmd/ads/handle/student"
	ads "rxt/cmd/ads/proto"
	"rxt/internal/api"
)

type Client struct {
	Server       ads.AdsServiceServer
	ServerClient ads.AdsServiceClient
	api.Api
}

func New() *Client {
	client := &Client{Api: api.Api{ServiceName: "ads", ServiceRemote: false}}
	client.Server = &student.Server{}
	return client
}

func (client *Client) Find(request *ads.FindRequest) (*ads.FindResponse, error) {
	return client.Server.Find(client.GetCtx(), request)
}
