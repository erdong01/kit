package main

import (
	"fmt"
	"github.com/erDong01/micro-kit/cluster"
	"github.com/erDong01/micro-kit/cluster/common"
	"github.com/erDong01/micro-kit/examples/account"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//账号管理类
	AccountMgr := new(account.AccountMgr)
	AccountMgr.Init(1000)

	var s network.ServerSocket
	s.Init("192.168.2.177", 8001)
	s.SetConnectType(network.CLIENT_CONNECT)
	s.BindPacketFunc(AccountMgr.PacketFunc)
	s.Start()
	clustert := new(cluster.Cluster)
	clustert.Init(1000, &common.ClusterInfo{Type: rpc3.SERVICE_ACCOUNTSERVER, Ip: "192.168.2.177", Port: 8001}, []string{"192.168.2.129:2379"}, "192.168.2.129:4222")
	clustert.BindPacketFunc(AccountMgr.PacketFunc)
	clustert.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	t := <-c
	fmt.Printf("server【%s】 exit ------- signal:[%v]", t)
}
