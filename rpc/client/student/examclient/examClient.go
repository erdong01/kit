package examclient

import (
	"rxt/cmd/exam/handler/student"
	pb "rxt/cmd/exam/proto/student"
	"rxt/internal/api"
)

// Client 客户端
type Client struct {
	IServer        pb.ExamServiceServer
	IServiceClient pb.ExamServiceClient
	api.Api
}

// New 构造函数
func New() *Client {
	client := Client{Api: api.Api{ServiceName: "exam", ServiceRemote: false}}
	client.IServer = &student.Server{}
	//client.IServiceClient = auth.NewLoginServiceClient(client.GetConn())
	return &client
}

// Classwork 获取课堂作业
func (client *Client) Classwork(request *pb.ClassworkRequest) (*pb.ClassworkResponse, error) {
	return client.IServer.Classwork(client.GetCtx(), request)
}
