package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"rxt/internal/wrong"
	"time"

	"google.golang.org/grpc"
)

type Api struct {
	ServiceName    string
	ServiceMethod  string
	ServiceAddress string
	ServiceRemote  bool
	ServiceConn    *grpc.ClientConn
	Err            error
}

//获取服务地址
func (*Api) getAddress(name string) string {
	serviceAddress := map[string]string{}
	serviceAddress["report"] = "localhost:50001"
	serviceAddress["auth"] = "localhost:50001"
	serviceAddress["exam"] = "localhost:50001"
	serviceAddress["barrier"] = "localhost:50001"
	return serviceAddress[name]
}

func (*Api) GetCtx() context.Context {
	// 10秒超时 仅remote call有效
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return ctx
}
func (client *Api) Call(Server interface{}, params interface{}) (interface{}, error) {
	var c reflect.Value
	server := reflect.ValueOf(Server).Elem()
	if client.ServiceRemote == true {
		c = server.FieldByName("IServiceClient")
	} else {
		c = server.FieldByName("IServer")
	}
	// 10秒超时 仅remote call有效
	ctx := client.GetCtx()
	value := make([]reflect.Value, 2)
	value[0] = reflect.ValueOf(ctx)
	value[1] = reflect.ValueOf(params)
	result := c.MethodByName(client.ServiceMethod).Call(value)
	if !result[1].IsNil() {
		return nil, result[1].Interface().(error)
	}
	return result[0].Interface(), nil
}

//获取连接
func (client *Api) GetConn() *grpc.ClientConn {
	var err error
	address := client.getAddress(client.ServiceName)
	fmt.Println(address)
	client.ServiceConn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		client.Err = wrong.New(http.StatusBadRequest, errors.New("rpc连接失败"))
	}
	defer client.ServiceConn.Close()
	return client.ServiceConn
}
