package main

import (
	"fmt"
	"github.com/erDong01/micro-kit/cluster"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/test/account"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//账号管理类
	AccountMgr := new(account.AccountMgr)
	AccountMgr.Init(1000)

	var s network.ServerSocket
	s.Init("127.0.0.1", 8001)
	s.BindPacketFunc(AccountMgr.PacketFunc)
	s.StartTcpServer()


	clustert := new(cluster.Cluster)
	clustert.Init(1000, &cluster.ClusterInfo{Type: rpc3.SERVICE_ACCOUNTSERVER, Ip: "127.0.0.1", Port: 7001}, []string{}, "127.0.0.1")
	clustert.BindPacketFunc(AccountMgr.PacketFunc)
	clustert.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	t := <-c
	fmt.Printf("server【%s】 exit ------- signal:[%v]", t)
}
