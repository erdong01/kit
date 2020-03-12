package api
//
//import (
//	"context"
//	"fmt"
//	"reflect"
//	authSc "rxt/cmd/auth/handler/sc"
//	auth "rxt/cmd/auth/proto/sc"
//	barrierHandle "rxt/cmd/barrier/handler/sc"
//	examSc "rxt/cmd/exam/handler/sc"
//	exam "rxt/cmd/exam/proto/sc"
//	"time"
//
//	"google.golang.org/grpc"
//
//	reportHandler "rxt/cmd/report/handler"
//	"rxt/cmd/report/proto/report"
//)
//
//var err error
//
//var serviceAddress = make(map[string]string)
//var rpcClientMap = make(map[string]func(cc *grpc.ClientConn) interface{})
//var innerMethodMap = make(map[string]func() interface{})
//
//func init() {
//	serviceAddress["report"] = "localhost:50001"
//	serviceAddress["auth"] = "localhost:50001"
//	serviceAddress["exam"] = "localhost:50001"
//	serviceAddress["barrier"] = "localhost:50001"
//
//	rpcClientMap["auth"] = func(cc *grpc.ClientConn) interface{} {
//		return auth.NewLoginServiceClient(cc)
//	}
//
//	rpcClientMap["report"] = func(cc *grpc.ClientConn) interface{} {
//		return report.NewReportRpcClient(cc)
//	}
//
//	rpcClientMap["exam"] = func(cc *grpc.ClientConn) interface{} {
//		return exam.NewExamServiceClient(cc)
//	}
//
//	innerMethodMap["auth"] = func() interface{} {
//		return &authSc.Server{}
//	}
//
//	innerMethodMap["report"] = func() interface{} {
//		return &reportHandler.Server{}
//	}
//
//	innerMethodMap["exam"] = func() interface{} {
//		return &examSc.Server{}
//	}
//
//	innerMethodMap["barrier"] = func() interface{} {
//		return &barrierHandle.Server{}
//	}
//}
//
//// IAPI api接口 支持grpc和内部调用
//type IAPI interface {
//	Call(params interface{}) (interface{}, error)
//}
//
//type api struct {
//	name   string
//	method string
//	remote bool
//	conn   *grpc.ClientConn
//}
//
//func getAddress(name string) string {
//	return serviceAddress[name]
//}
//
//// New 构造函数
//func New(name string, method string, remote bool) IAPI {
//	client := &api{
//		name:   name,
//		method: method,
//		remote: remote,
//	}
//
//	if client.remote {
//		address := getAddress(client.name)
//		client.conn, err = grpc.Dial(address, grpc.WithInsecure())
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	return client
//}
//
//func (client *api) Call(params interface{}) (interface{}, error) {
//	var c interface{}
//	if client.remote {
//		defer client.conn.Close()
//		c = rpcClientMap[client.name](client.conn)
//	} else {
//		c = innerMethodMap[client.name]()
//	}
//
//	// 10秒超时 仅remote call有效
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//	defer cancel()
//
//	value := make([]reflect.Value, 2)
//	value[0] = reflect.ValueOf(ctx)
//	value[1] = reflect.ValueOf(params)
//
//	result := reflect.ValueOf(c).MethodByName(client.method).Call(value)
//
//	if !result[1].IsNil() {
//		return nil, result[1].Interface().(error)
//	}
//
//	return result[0].Interface(), nil
//}
