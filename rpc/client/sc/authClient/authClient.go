package authClient

import (
	"rxt/cmd/auth/handler/sc"
	auth "rxt/cmd/auth/proto/sc"
	"rxt/internal/api"
)

type Client struct {
	IServer        auth.AuthServiceServer
	IServiceClient auth.AuthServiceClient
	api.Api
}

func New() *Client {
	client := Client{Api: api.Api{ServiceName: "auth", ServiceRemote: false}}
	client.IServer = &sc.Server{}
	//client.IServiceClient = auth.NewLoginServiceClient(client.GetConn())
	return &client
}

func (client *Client) Logic(request *auth.AuthRequest) (*auth.AuthResponse, error) {
	return client.IServer.Login(client.GetCtx(), request)
}
func (client *Client) Validate(request *auth.TokenRequest) (*auth.UserResponse, error) {
	return client.IServer.Validate(client.GetCtx(), request)
}
