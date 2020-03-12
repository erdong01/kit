package barrierClient

import (
	"rxt/cmd/barrier/handler/sc"
	barrier "rxt/cmd/barrier/proto/sc"
	"rxt/internal/api"
)

type Client struct {
	IServer        barrier.BarrierServiceServer
	IServiceClient barrier.BarrierServiceClient
	api.Api
}

func New() *Client {
	client := Client{Api: api.Api{ServiceName: "auth", ServiceRemote: false}}
	client.IServer = &sc.Server{}
	//client.IServiceClient = auth.NewLoginServiceClient(client.GetConn())
	return &client
}

func (client *Client) Skip(request *barrier.Request) (*barrier.Response, error) {
	return client.IServer.Skip(client.GetCtx(), request)
}
