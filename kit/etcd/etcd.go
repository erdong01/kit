package etcd

import (
	"context"
	"fmt"
	"github.com/erDong01/micro-kit/config"
	config2 "github.com/erDong01/micro-kit/config"
	"github.com/go-kit/kit/sd/etcdv3"
	"time"
)

func New() etcdv3.Client {
	ctx := context.Background()
	configEtcd := config.GetEtcd()
	config2.New().Get(configEtcd, "etcd")
	//etcd连接参数
	option := etcdv3.ClientOptions{DialTimeout: time.Second * 3, DialKeepAlive: time.Second * 3}
	client, err := etcdv3.NewClient(ctx, configEtcd.Addr, option)
	if err != nil {
		fmt.Println("etcd链接失败")
		panic(err)
	}
	return client
}
