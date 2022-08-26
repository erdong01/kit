package network

import (
	"net"
	"sync/atomic"

	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools/vector"
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
	MAX_SEND_CHAN  = 512
	HEART_TIME_OUT = 30
)

type (
	ClientClose func(id uint32) error

	PacketFunc   func(packet rpc.Packet) bool //回调函数
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

		totalNum     int
		acceptedNum  int
		connectedNum int

		sendTimes      int
		receiveTimes   int
		PacketFuncList *vector.Vector

		half         bool
		halfSize     int
		packetParser PacketParser
		heartTime    int
		bKcp         bool

		clientClose ClientClose
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
		SendMsg(rpc.RpcHead, string, ...interface{})
		Send(rpc.RpcHead, []byte) int
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
	this.ReceiveBufferSize = 1024
	this.SetState(SSF_NULL)
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
	return false
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

func (this *Socket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
}

func (this *Socket) Send(rpc.RpcHead, []byte) int {
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
	buff := rpc.Marshal(rpc.RpcHead{}, funcName, params...)
	this.HandlePacket(buff)
}

func (this *Socket) HandlePacket(buff []byte) {
	packet := rpc.Packet{Id: this.clientId, Buff: buff}
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
