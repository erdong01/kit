package network

import (
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools/vector"
	"net"
)

const (
	SSF_ACCEPT    = iota
	SSF_CONNECT   = iota
	SSF_SHUT_DOWN = iota //已经关闭
)
const (
	MAX_SEND_CHAN = 100
)

const (
	CLIENT_CONNECT = iota //对外
	SERVER_CONNECT = iota //对内
)

type (
	PacketFunc   func(packet rpc3.Packet) bool //回调函数
	HandlePacket func(buff []byte)
	Socket       struct {
		IP                string
		state             int
		Port              int
		Zone              string
		ReceiveBufferSize int //单次接收缓存
		connectType       int

		clientId uint32
		seq      int64

		Conn net.Conn

		sendTimes    int
		receiveTimes int
		shuttingDown bool

		PacketFuncList *vector.Vector

		half         bool
		halfSize     int
		packetParser PacketParser
	}

	ISocket interface {
		Init(string, int) bool
		Start() bool
		Stop() bool
		Run() bool
		Restart() bool
		Connect() bool
		Disconnect(bool) bool
		OnNetFail(int)
		Clear()
		Close()
		SendMsg(rpc3.RpcHead, string, ...interface{}) int
		Send(rpc3.RpcHead, []byte) int
		CallMsg(string, ...interface{}) //回调消息处理

		GetId() uint32
		GetState() int
		SetReceiveBufferSize(int)
		GetReceiveBufferSize() int
		SetMaxPacketLen(int)
		GetMaxPacketLen() int
		BindPacketFunc(PacketFunc)
		SetConnectType(int)
		SetTcpConn(net.Conn)
		HandlePacket([]byte)
	}
)

func (this *Socket) Init(string, int) {
	this.PacketFuncList = vector.NewVector()
	this.ReceiveBufferSize = 1024
	this.state = SSF_SHUT_DOWN
	this.connectType = SERVER_CONNECT
	this.half = false
	this.halfSize = 0
	this.packetParser = NewPacketParser(PacketConfig{Func: this.HandlePacket})
}

func (this *Socket) Start() bool {
	return true
}
func (this *Socket) Stop() bool {
	this.shuttingDown = true
	return true
}
func (this *Socket) Run() bool {
	return true
}

func (this *Socket) Restart() bool {
	return true
}

func (this *Socket) Connect() bool {
	return true
}

func (this *Socket) OnNetFail(int) {
	this.Stop()
}

func (this *Socket) GetId() uint32 {
	return this.clientId
}
func (this *Socket) Disconnect(bool) bool {
	return true
}

func (this *Socket) GetState() int {
	return this.state
}
func (this *Socket) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) int {
	return 0
}

func (this *Socket) Send(rpc3.RpcHead, []byte) int {
	return 0
}

func (this *Socket) Clear() {
	this.state = SSF_SHUT_DOWN
	this.Conn = nil
	this.ReceiveBufferSize = 1024
	this.shuttingDown = false
	this.half = false
	this.halfSize = 0
}

func (this *Socket) Close() {
	if this.Conn != nil {
		this.Conn.Close()
	}
	this.Clear()
}

func (this *Socket) GetMaxPacketLen() int {
	return this.packetParser.MaxPacketLen
}

func (this *Socket) SetMaxPacketLen(maxReceiveSize int) {
	this.packetParser.MaxPacketLen = maxReceiveSize
}

func (this *Socket) GetReceiveBufferSize() int {
	return this.ReceiveBufferSize
}

func (this *Socket) SetReceiveBufferSize(maxSendSize int) {
	this.ReceiveBufferSize = maxSendSize
}

func (this *Socket) SetConnectType(nType int) {
	this.connectType = nType
}

func (this *Socket) SetTcpConn(conn net.Conn) {
	this.Conn = conn
}

func (this *Socket) BindPacketFunc(callfunc PacketFunc) {
	this.PacketFuncList.PushBack(callfunc)
}

func (this *Socket) CallMsg(funcName string, params ...interface{}) {
	buff := rpc.Marshal(rpc3.RpcHead{}, funcName, params...)
	this.HandlePacket(buff)
}

func (this *Socket) HandlePacket(buff []byte) {
	packet := rpc3.Packet{Id: this.clientId, Buff: buff}
	for _, v := range this.PacketFuncList.Values() {
		if v.(PacketFunc)(packet) {
			break
		}
	}
}
