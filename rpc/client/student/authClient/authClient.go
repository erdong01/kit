package authClient

import (
	"github.com/erDong01/micro-kit/cmd/auth/handler/student"
	auth "github.com/erDong01/micro-kit/cmd/auth/proto/student"
	"github.com/erDong01/micro-kit/internal/api"
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
