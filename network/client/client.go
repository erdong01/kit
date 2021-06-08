package main

// golang实现带有心跳检测的tcp长连接
// server

import (
	"fmt"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
)

var Dch chan bool
var (
	CLIENT *network.ClientSocket
)

func main() {
	Dch = make(chan bool)
	//	conn, err := net.Dial("tcp", "127.0.0.1:6666")
	CLIENT = new(network.ClientSocket)
	CLIENT.Init("192.168.2.129", 8001)
	CLIENT.Start()
	go Handler()
	select {
	case <-Dch:
		fmt.Println("关闭连接")
	}
}

func Handler() {
	// 直到register ok
	head := rpc3.RpcHead{Code: 100, ActorName: "Account"}
	byteD := rpc.Marshal(head, "Account_Login", "test", 88, 88, 88)
	CLIENT.Send(head, byteD)
}
