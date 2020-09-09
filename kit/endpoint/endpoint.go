package endpoint

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"google.golang.org/grpc"
	"io"
	"reflect"
	"time"
)

type ClientConn func(conn *grpc.ClientConn) interface{}

// Endpoints are exposed
type Endpoints struct {
	Factory    sd.Factory
	ClientFunc interface{}
	Instancer  sd.Instancer
	Logger     log.Logger
	C          ClientConn
}

func New(clientConn ClientConn, instancer sd.Instancer, logger log.Logger, c interface{}) {
	var e Endpoints
	e.ReqFactory(clientConn).MakeEndpoint(instancer, logger, c)
}

func (this *Endpoints) ReqFactory(c ClientConn) *Endpoints {
	this.Factory = func(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			fmt.Println(instanceAddr)
			conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				panic("connect error")
			}
			e := c(conn)
			return e, nil
		}, nil, nil
	}
	return this
}

func (this *Endpoints) MakeEndpoint(instancer sd.Instancer, logger log.Logger, c interface{}) {
	ctx := context.Background()
	//创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
	endpointer := sd.NewEndpointer(instancer, this.Factory, logger) //reqFactory自定义的函数，主要用于端点层（endpoint）接受并显示数据
	endpoints, _ := endpointer.Endpoints()
	fmt.Println("服务有", len(endpoints), "个")
	//创建负载均衡器
	balancer := lb.NewRoundRobin(endpointer)
	//e, err := balancer.Endpoint()
	//fmt.Println(e(ctx, nil), "test")
	/**	也可以通过retry定义尝试次数进行请求	*/
	reqEndPoint := lb.Retry(3, 3*time.Second, balancer)
	s, _ := reqEndPoint(ctx, nil)
	rx := reflect.ValueOf(c).Elem()
	rx.Set(reflect.ValueOf(s))
	return
}
