package actor

import (
	"context"
	"fmt"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var (
	IdSeed int64
)

type (
	Actor struct {
		CallChan  chan CallIO //rpc chan
		ActorChan chan int    //use for states
		id        int64
		CallMap   map[string]*CallFunc
		pTimer    *time.Ticker //定时器
		TimerCall func()       //定时器触发函数
		bStart    bool
		mTrace    traceInfo //trace func
	}
	CallIO struct {
		rpc3.RpcHead
		Buff []byte
	}
	CallFunc struct {
		Func       interface{}
		FuncType   reflect.Type
		FuncVal    reflect.Value
		FuncParams string
	}

	traceInfo struct {
		funcName  string
		fileName  string
		filePath  string
		className string
	}
	IActor interface {
		Init(chanNum int)
		Stop()
		Start()
		FindCall(funcName string) *CallFunc
		RegisterCall(funcName string, call interface{})
		SendMsg(head rpc3.RpcHead, funcName string, params ...interface{})
		Send(head rpc3.RpcHead, buff []byte)
		PacketFunc(packet rpc3.Packet) bool                    //回调函数
		RegisterTimer(duration time.Duration, fun interface{}) //注册定时器,时间为纳秒 1000 * 1000 * 1000
		GetId() int64
		GetRpcHead(ctx context.Context) rpc3.RpcHead //rpc is safe
	}
)

func (this *Actor) Init(chanNum int) {
	this.CallChan = make(chan CallIO, chanNum)
	this.ActorChan = make(chan int, 1)
	this.id = AssignActorId()
	this.CallMap = make(map[string]*CallFunc)
	this.pTimer = time.NewTicker(1<<63 - 1)
	this.mTrace.Init()
}

func (this *Actor) RegisterTimer(duration time.Duration, fun interface{}) {
	this.pTimer.Stop()
	this.pTimer = time.NewTicker(duration)
	this.TimerCall = fun.(func())

}

func (this *Actor) clear() {
	this.id = 0
	this.bStart = false
	if this.pTimer != nil {
		this.pTimer.Stop()
	}
	this.CallMap = make(map[string]*CallFunc)
}

func (this *Actor) Stop() {
	this.ActorChan <- 1
}

func (this *Actor) ClientSocket(ctx context.Context) *network.ServerSocketClient {
	rpcHead := ctx.Value("rpcHead").(rpc3.RpcHead)
	return network.SocketServer.GetClientById(rpcHead.SocketId)
}

func (this *Actor) Start() {
	if this.bStart == false {
		go this.run()
		this.bStart = true
	}
}

func (this *Actor) run() {
	for {
		if !this.loop() {
			break
		}
	}
	this.clear()
}

func (this *Actor) loop() bool {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
		}
	}()
	select {
	case io := <-this.CallChan:
		this.call(io)
	}
	return true
}

func (this *Actor) call(io CallIO) {
	rpcPacket, _ := rpc.Unmarshal(io.Buff)
	head := io.RpcHead
	funcName := rpcPacket.FuncName
	pFunc := this.FindCall(funcName)

	if pFunc != nil {
		f := pFunc.FuncVal
		k := pFunc.FuncType
		rpcPacket.RpcHead.SocketId = io.SocketId
		params := rpc.UnmarshalBody(rpcPacket, k)
		if funcName != "cluster_add" {
			fmt.Println(params)
		}
		if len(params) >= 1 {
			in := make([]reflect.Value, len(params))
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}
			this.Trace(funcName)
			ret := f.Call(in)
			this.Trace("")
			if ret != nil && head.Reply != "" {
				ret = append([]reflect.Value{reflect.ValueOf(&head)}, ret...)
				rpc.GCall.Call(ret)
			}
		} else {
			log.Printf("func [%s] params at least one context", funcName)
		}
	}
}

func (this *Actor) Trace(funcName string) {
	this.mTrace.funcName = funcName
}

func (this *Actor) FindCall(funcName string) *CallFunc {
	funcName = strings.ToLower(funcName)
	fun, exist := this.CallMap[funcName]
	if exist == true {
		return fun
	}
	return nil
}
func AssignActorId() int64 {
	atomic.AddInt64(&IdSeed, 1)
	return int64(IdSeed)
}
func (this *Actor) GetId() int64 {
	return this.id
}

func (this *Actor) GetRpcHead(ctx context.Context) rpc3.RpcHead {
	rpcHead := ctx.Value("rpcHead").(rpc3.RpcHead)
	return rpcHead
}

func (this *Actor) PacketFunc(packet rpc3.Packet) bool {
	rpcPacket, head := rpc.UnmarshalHead(packet.Buff)
	if this.FindCall(rpcPacket.FuncName) != nil {
		head.SocketId = packet.Id
		head.Reply = packet.Reply
		this.Send(head, packet.Buff)
		return true
	}
	return false
}
func (this *Actor) RegisterCall(funcName string, call interface{}) {
	funcName = strings.ToLower(funcName)
	if this.FindCall(funcName) != nil {
		log.Fatalln("actor error [%s] 消息重复定义", funcName)
	}

	callfunc := &CallFunc{Func: call, FuncVal: reflect.ValueOf(call), FuncType: reflect.TypeOf(call), FuncParams: reflect.TypeOf(call).String()}
	this.CallMap[funcName] = callfunc
}
func (this *Actor) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
	head.SocketId = 0
	this.Send(head, rpc.Marshal(head, funcName, params...))
}

func (this *Actor) Send(head rpc3.RpcHead, buff []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
		}
	}()
	var io CallIO
	io.RpcHead = head
	io.Buff = buff
	this.CallChan <- io
}
func (this *traceInfo) Init() {
	_, file, _, bOk := runtime.Caller(2)
	if bOk {
		index := strings.LastIndex(file, "/")
		if index != -1 {
			this.fileName = file[index+1:]
			this.filePath = file[:index]
			indexTow := strings.LastIndex(this.fileName, ".")
			if indexTow != -1 {
				this.className = strings.ToLower(this.fileName[:indexTow])
			}
		}
	}
}
func (this *traceInfo) ToString() string {
	return fmt.Sprintf("trace go file[%s] call[%s]\n", this.fileName, this.funcName)
}
