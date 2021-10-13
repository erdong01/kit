package network

import (
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools/vector"
	"net"
	"sync/atomic"
)

const (
	SSF_NULL = iota
	SSF_RUN  = iota
	SSF_STOP = iota //已经关闭
)

const (
	CLIENT_CONNECT = iota //对外
	SERVER_CONNECT = iota //对内
)
const (
	MAX_SEND_CHAN  = 100
	HEART_TIME_OUT = 30
)

type (
	ClientClose  func(id uint32) error         //客户关闭回调
	PacketFunc   func(packet rpc3.Packet) bool //回调函数
	HandlePacket func(buff []byte)

	Op struct {
		kcp bool
	}

	OpOption func(*Op)

	Socket struct {
		Conn              net.Conn
		Port              int
		IP                string
		state             int32
		connectType       int
		ReceiveBufferSize int //单次接收缓存

		clientId uint32
		seq      int64

		sendTimes      int
		receiveTimes   int
		PacketFuncList *vector.Vector

		half         bool
		halfSize     int
		packetParser PacketParser
		heartTime    int
		bKcp         bool

		clientClose ClientClose //客户关闭回调
	}

	ISocket interface {
		Init(string, int, ...OpOption) bool
		Start() bool
		Stop() bool
		Run() bool
		Restart() bool
		Connect() bool
		Disconnect(bool) bool
		OnNetFail(int)
		Clear()
		Close()
		SendMsg(rpc3.RpcHead, string, ...interface{})
		Send(rpc3.RpcHead, []byte) int
		CallMsg(string, ...interface{}) //回调消息处理

		GetId() uint32
		GetState() int32
		SetState(int32)
		SetReceiveBufferSize(int)
		GetReceiveBufferSize() int
		SetMaxPacketLen(int)
		GetMaxPacketLen() int
		BindPacketFunc(PacketFunc)
		SetConnectType(int)
		SetConn(net.Conn)

		HandlePacket([]byte)

		SetClientClose(ClientClose)
		GetClientClose() ClientClose
	}
)

func (op *Op) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func WithKcp() OpOption {
	return func(op *Op) {
		op.kcp = true
	}
}

func (this *Socket) Init(ip string, port int, params ...OpOption) bool {
	op := &Op{}
	op.applyOpts(params)
	this.PacketFuncList = vector.NewVector()
	this.SetState(SSF_NULL)
	this.ReceiveBufferSize = 1024
	this.connectType = SERVER_CONNECT
	this.half = false
	this.halfSize = 0
	this.heartTime = 0
	this.packetParser = NewPacketParser(PacketConfig{Func: this.HandlePacket})
	if op.kcp {
		this.bKcp = true
	}
	return true
}

func (this *Socket) Start() bool {
	return true
}

func (this *Socket) Stop() bool {
	if this.Conn != nil && atomic.CompareAndSwapInt32(&this.state, SSF_RUN, SSF_STOP) {
		this.Conn.Close()
	}
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

func (this *Socket) Disconnect(bool) bool {
	return true
}

func (this *Socket) OnNetFail(int) {
	this.Stop()
}

func (this *Socket) GetId() uint32 {
	return this.clientId
}

func (this *Socket) GetState() int32 {
	return atomic.LoadInt32(&this.state)
}

func (this *Socket) SetState(state int32) {
	atomic.StoreInt32(&this.state, state)
}

func (this *Socket) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
}

func (this *Socket) Send(rpc3.RpcHead, []byte) int {
	return 0
}

func (this *Socket) Clear() {
	this.SetState(SSF_NULL)
	this.Conn = nil
	this.ReceiveBufferSize = 1024
	this.half = false
	this.halfSize = 0
	this.heartTime = 0
}

func (this *Socket) Close() {
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

func (this *Socket) SetConn(conn net.Conn) {
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

func (this *Socket) SetClientClose(c ClientClose) {
	this.clientClose = c
}
func (this *Socket) GetClientClose() ClientClose {
	return this.clientClose
}
