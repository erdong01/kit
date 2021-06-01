package network

import (
	"github.com/erDong01/micro-kit/tools/vector"
	"net"
)

type HandlePacket func(buff []byte)
type Socket struct {
	IP                string
	Port              int
	Zone              string
	ReceiveBufferSize int //单次接收缓存
	Conn              net.Conn
	PacketParser      PacketParser
	PacketFuncList    *vector.Vector
}

const (
	MAX_SEND_CHAN = 100
)

func (this *Socket) Init(string, int) {
	this.ReceiveBufferSize = 1024
}
