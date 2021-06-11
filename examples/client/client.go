package main

// golang实现带有心跳检测的tcp长连接
// server

import (
	"fmt"
	"github.com/erDong01/micro-kit/examples/message"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
)

var Dch chan bool
var (
	CLIENT *network.ClientSocket
)

func main() {

	message.InitClient()

	Dch = make(chan bool)

	CLIENT = new(network.ClientSocket)

	CLIENT.Init("192.168.2.231", 31700)

	PACKET = new(EventProcess)

	PACKET.Init(1)

	CLIENT.BindPacketFunc(PACKET.PacketFunc)

	PACKET.Client = CLIENT

	CLIENT.Start()

	PACKET.LoginGame()

	//go Handler()
	select {
	case <-Dch:
		fmt.Println("关闭连接")
	}
}

func Handler() {
	// 直到register ok
	head := rpc3.RpcHead{Code: 100, ActorName: "UserPrcoess"}
	byteD := rpc.Marshal(head, "C_G_LogoutRequest", "test", 88, 88, 88)
	CLIENT.Send(head, byteD)
}
