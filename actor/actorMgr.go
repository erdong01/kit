package actor

import (
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
)

type ACTOR_TYPE uint32

const (
	ACTOR_TYPE_SINGLETON ACTOR_TYPE = iota //单列
	ACTOR_TYPE_PLAYER    ACTOR_TYPE = iota //玩家 必须初始一个全局的actor 作为类型判断
) //ACTOR_TYPE

const (
	MAX_RPC_TAG = 10
)

//一些全局的actor,不可删除的,不用锁考虑性能
//不是全局的actor,请使用actor pool
type (
	Op struct {
		name         string //name
		aType        ACTOR_TYPE
		rpcMethodMap map[string]string
	}

	OpOption func(*Op)
	ActorMgr struct {
		actorMap     map[reflect.Type]IActor
		actorNameMap map[string]IActor
		msgMap       map[string]IActor
		rpcMethodMap map[reflect.Type]map[string]string
		playerMap    map[int64]IActor
		playerLock   *sync.RWMutex
		bStart       bool
	}
	IActorMgr interface {
		Init()
		RegisterActor(IActor, ...string) //注册回调
		PacketFunc(rpc3.Packet) bool     //回调函数
		SendMsg(rpc3.RpcHead, string, ...interface{})
	}
	ICluster interface {
		BindPacketFunc(packetFunc network.PacketFunc)
	}
)

func (op *Op) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func (op *Op) IsActorType(actorType ACTOR_TYPE) bool {
	return op.aType == actorType
}

func WithName(name string) OpOption {
	return func(op *Op) {
		op.name = name
	}
}

func WithRpcMethodMap(rpcMethodMap map[string]string) OpOption {
	return func(op *Op) {
		op.rpcMethodMap = rpcMethodMap
	}
}
func (this *ActorMgr) Init() {
	this.actorMap = make(map[reflect.Type]IActor)
	this.actorNameMap = make(map[string]IActor)
	this.msgMap = make(map[string]IActor)
	this.rpcMethodMap = map[reflect.Type]map[string]string{}
	this.playerMap = make(map[int64]IActor)
	this.playerLock = &sync.RWMutex{}
}

func (this *ActorMgr) Start() {
	this.bStart = true
}

func (this *ActorMgr) RegisterActor(pActor IActor, params ...OpOption) {
	op := Op{}
	op.applyOpts(params)
	if len(op.name) == 0 {
		op.name = tools.GetClassName(pActor)
	}
	rType := reflect.TypeOf(pActor)
	_, bEx := this.actorMap[rType]
	if bEx {
		log.Panicf("InitActor actor[%s] must  global variable", op.name)
		return
	}
	//rpc
	shareRpcMethodMap := GetRpcMethodMap(rType, "share_rpc")
	methodNum := rType.NumMethod()
	this.rpcMethodMap[rType] = map[string]string{}
	for i := 0; i < methodNum; i++ {
		m := rType.Method(i)
		if m.Type.NumIn() >= 2 {
			if m.Type.In(1).String() == "context.Context" {
				funcName := strings.ToLower(m.Name)
				methodName := m.Name
				_, bInShare := shareRpcMethodMap[funcName]
				if !bInShare {
					pMsgHandle, bEx := this.msgMap[funcName]
					if bEx && pMsgHandle != nil {
						log.Panicf("RegisterFuncName [%s] exist_actor [%s] actor [%s]", methodName, pMsgHandle.GetName(), op.name)
						return
					}
					this.msgMap[funcName] = pActor
				}
				this.rpcMethodMap[rType][funcName] = methodName
			}
		}
	}
	op.rpcMethodMap = this.rpcMethodMap[rType]
	pActor.Register(pActor, op)
	this.actorMap[rType] = pActor
	this.actorNameMap[op.name] = pActor
}

func (this *ActorMgr) AddPlayer(pActor IActor) {
	rType := reflect.TypeOf(pActor)
	op := Op{aType: ACTOR_TYPE_PLAYER, name: this.actorMap[rType].GetName(), rpcMethodMap: this.rpcMethodMap[rType]}
	pActor.Register(pActor, op)
	this.playerLock.Lock()
	this.playerMap[pActor.GetId()] = pActor
	this.playerLock.Unlock()
}
func (this *ActorMgr) DelPlayer(Id int64) {
	this.playerLock.Lock()
	delete(this.playerMap, Id)
	this.playerLock.Unlock()
}

func (this *ActorMgr) GetPlayer(Id int64) IActor {
	this.playerLock.RLock()
	pActor, bEx := this.playerMap[Id]
	this.playerLock.RUnlock()
	if bEx {
		return pActor
	}
	return nil
}

func (this *ActorMgr) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
	head.SocketId = 0
	this.SendActor(funcName, head, rpc.Marshal(head, funcName, params...))
}
func (this *ActorMgr) SendActor(funcName string, head rpc3.RpcHead, packet rpc3.Packet) bool {
	var pActor IActor
	funcName = strings.ToLower(funcName)
	bEx := false
	if head.ActorName != "" {
		pActor, bEx = this.actorNameMap[head.ActorName]
	} else {
		pActor, bEx = this.msgMap[funcName]
	}
	if bEx && pActor != nil {
		if pActor.HasRpc(funcName) {
			switch pActor.GetActorType() {
			case ACTOR_TYPE_SINGLETON:
				pActor.GetAcotr().Send(head, packet)
				return true
			case ACTOR_TYPE_PLAYER:
				if head.Id != 0 {
					pActor := this.GetPlayer(head.Id)
					if pActor != nil {
						pActor.GetAcotr().Send(head, packet)
						return true
					}

				}
			}
		}
	}
	return false
}
func (this *ActorMgr) PacketFunc(packet rpc3.Packet) bool {
	rpcPacket, head := rpc.Unmarshal(packet.Buff)
	packet.RpcPacket = rpcPacket
	head.SocketId = packet.Id
	head.Reply = packet.Reply
	return this.SendActor(rpcPacket.FuncName, head, packet)
}

var (
	MGR *ActorMgr
)

func init() {
	MGR = &ActorMgr{}
	MGR.Init()
}
