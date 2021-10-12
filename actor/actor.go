package actor

import (
	"context"
	"fmt"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools/mpsc"
	"github.com/erDong01/micro-kit/tools/timer"
	"github.com/erDong01/micro-kit/wrong"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	IdSeed int64
)

const (
	AsfNull = iota
	AsfRun  = iota
	AsfStop = iota
)

type (
	Actor struct {
		acotrChan chan int //use for states
		id        int64
		CallMap   map[string]*CallFunc
		state     int32
		mTrace    traceInfo //trace func
		mailBox   *mpsc.Queue
		mailIn    int32
		mailChan  chan bool
		mTimerId  *int64
	}

	IActor interface {
		Init()
		Stop()
		Start()
		FindCall(funcName string) *CallFunc
		RegisterCall(funcName string, call interface{})
		SendMsg(head rpc3.RpcHead, funcName string, params ...interface{})
		Send(head rpc3.RpcHead, buff []byte)
		PacketFunc(packet rpc3.Packet) bool                                       //回调函数
		RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) //注册定时器,时间为纳秒 1000 * 1000 * 1000
		GetId() int64
		GetState() int32
		setState(state int32)
		GetRpcHead(ctx context.Context) rpc3.RpcHead //rpc is safe
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
)

const (
	DESDORY_EVENT = iota
)

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

func (this *Actor) GetState() int32 {
	return atomic.LoadInt32(&this.state)
}

func (this *Actor) setState(state int32) {
	atomic.StoreInt32(&this.state, state)
}

func (this *Actor) Init() {
	this.mailChan = make(chan bool)
	this.mailBox = mpsc.New()
	this.acotrChan = make(chan int, 1)
	this.id = AssignActorId()
	this.CallMap = make(map[string]*CallFunc)
	//trance
	this.RegisterCall("UpdateTimer", func(ctx context.Context, p *int64) {
		func1 := (*func())(unsafe.Pointer(p))
		this.Trace("timer")
		(*func1)()
		this.Trace("")
	})
	this.mTrace.Init()
}

func (this *Actor) RegisterTimer(duration time.Duration, fun func(), opts ...timer.OpOption) {
	if this.mTimerId == nil {
		this.mTimerId = new(int64)
		*this.mTimerId = this.id
	}
	timer.RegisterTimer(this.mTimerId, duration, func() {
		this.SendMsg(rpc3.RpcHead{}, "UpdateTimer", (*int64)(unsafe.Pointer(&fun)))
	}, opts...)
}

func (this *Actor) clear() {
	this.id = 0
	this.setState(AsfNull)
	timer.StopTimer(this.mTimerId)
	this.CallMap = make(map[string]*CallFunc)
}

func (this *Actor) Stop() {
	if atomic.CompareAndSwapInt32(&this.state, AsfRun, AsfStop) {
		this.acotrChan <- DESDORY_EVENT
	}
}

func (this *Actor) Start() {
	if atomic.CompareAndSwapInt32(&this.state, AsfNull, AsfRun) {
		go this.run()
	}
}
func (this *Actor) FindCall(funcName string) *CallFunc {
	funcName = strings.ToLower(funcName)
	fun, exist := this.CallMap[funcName]
	if exist == true {
		return fun
	}
	return nil
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
	this.mailBox.Push(io)
	if atomic.CompareAndSwapInt32(&this.mailIn, 0, 1) {
		this.mailChan <- true
	}
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

func (this *Actor) Trace(funcName string) {
	this.mTrace.funcName = funcName
}

func (this *Actor) call(io CallIO) {
	defer func() {
		if err := recover(); err != nil {
			wrong.TraceCode(this.mTrace.ToString(), err)
		}
	}()
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
			fmt.Println(funcName, params)
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
func (this *Actor) consume() {
	atomic.StoreInt32(&this.mailIn, 0)
	for data := this.mailBox.Pop(); data != nil; data = this.mailBox.Pop() {
		this.call(data.(CallIO))
	}
}

func (this *Actor) loop() bool {
	defer func() {
		if err := recover(); err != nil {
			wrong.TraceCode(this.mTrace.ToString(), err)
		}
	}()
	select {
	case <-this.mailChan:
		this.consume()
	case msg := <-this.acotrChan:
		if msg == DESDORY_EVENT {
			return false
		}
	}
	return true
}
func (this *Actor) run() {
	for {
		if !this.loop() {
			break
		}
	}
	this.clear()
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

//ClientSocket 给客户发送消息
func (this *Actor) ClientSocket(ctx context.Context) *network.ServerSocketClient {
	rpcHead := ctx.Value("rpcHead").(rpc3.RpcHead)
	return network.SocketServer.GetClientById(rpcHead.SocketId)
}
