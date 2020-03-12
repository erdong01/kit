package examClient

import (
	"rxt/cmd/exam/handler/sc"
	exam "rxt/cmd/exam/proto/sc"
	"rxt/internal/api"
)

type Client struct {
	IServer        exam.ExamServiceServer
	IServiceClient exam.ExamServiceClient
	api.Api
}

func New() *Client {
	client := Client{IServer: &sc.Server{},
		Api: api.Api{ServiceName: "exam", ServiceRemote: false},
	}
	//client.IServiceClient = exam.NewExamServiceClient(client.GetConn())
	return &client
}

//func (client *Client) Submit(request *exam.ExamRequest) (*exam.ExamResponse, error) {
//	client.Method = "Submit"
//	res, err := client.Call(client, request)
//	return res.(*exam.ExamResponse), err
//}

func (client *Client) Submit(request *exam.ExamRequest) (*exam.ExamResponse, error) {
	return client.IServer.Submit(client.GetCtx(), request)
}
