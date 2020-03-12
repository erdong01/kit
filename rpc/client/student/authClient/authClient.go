package authClient

import (
	"rxt/cmd/auth/handler/student"
	auth "rxt/cmd/auth/proto/student"
	"rxt/internal/api"
)

type Client struct {
	IServer        auth.AuthServiceServer
	IServiceClient auth.AuthServiceClient
	api.Api
}

func New() *Client {
	client := Client{Api: api.Api{ServiceName: "auth", ServiceRemote: false}}
	client.IServer = &student.Server{}
	//client.IServiceClient = auth.NewLoginServiceClient(client.GetConn())
	return &client
}

func (client *Client) Logic(request *auth.LogicRequest) (*auth.LogicResponse, error) {
	return client.IServer.Login(client.GetCtx(), request)
}
