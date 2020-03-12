package reportClient

import (
	"rxt/cmd/report/handler"
	"rxt/cmd/report/proto/report"
	"rxt/internal/api"
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
