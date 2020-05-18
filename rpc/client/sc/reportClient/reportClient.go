package reportClient

import (
	"github.com/erDong01/gin-kit/cmd/report/handler"
	"github.com/erDong01/gin-kit/cmd/report/proto/report"
	"github.com/erDong01/gin-kit/internal/api"
)

type Client struct {
	IServer        report.ReportRpcServer
	IServiceClient report.ReportRpcClient
	api.Api
}

func New() *Client {
	client := Client{Api: api.Api{ServiceName: "auth", ServiceRemote: false}}
	client.IServer = &handler.Server{}
	//client.IServiceClient = auth.NewLoginServiceClient(client.GetConn())
	return &client
}

func (client *Client) Show(request *report.ReportRequest) (*report.ReportResponse, error) {
	return client.IServer.Show(client.GetCtx(), request)
}
