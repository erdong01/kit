package actor

import (
	"sync"

	"github.com/erDong01/micro-kit/rpc"
)

type ActorPool struct {
	Actor
	ActorMap  map[int64]IActor
	ActorLock *sync.RWMutex
}
type IActorPool interface {
	GetActor(Id int64) IActor                         //获取actor
	AddActor(Id int64, pActor IActor)                 //添加actor
	DelActor(Id int64)                                //删除actor
	BoardCast(funcName string, params ...interface{}) //广播actor
	GetActorNum() int
}

func (this *ActorPool) Init(chanNum int) {
	this.ActorMap = make(map[int64]IActor)
	this.ActorLock = &sync.RWMutex{}
	this.Actor.Init()
}

func (this *ActorPool) GetActor(Id int64) IActor {
	this.ActorLock.RLock()
	pActor, bEx := this.ActorMap[Id]
	this.ActorLock.RUnlock()
	if bEx {
		return pActor
	}
	return nil
}

func (this *ActorPool) AddActor(Id int64, actor IActor) {
	this.ActorLock.Lock()
	this.ActorMap[Id] = actor
	this.ActorLock.Unlock()
}

func (this *ActorPool) DelActor(Id int64) {
	this.ActorLock.Lock()
	delete(this.ActorMap, Id)
	this.ActorLock.Unlock()
}

func (this *ActorPool) GetActorNum() int {
	nLen := 0
	this.ActorLock.RLock()
	nLen = len(this.ActorMap)
	this.ActorLock.RUnlock()
	return nLen
}

func (this *ActorPool) BoardCast(funcName string, params ...interface{}) {
	this.ActorLock.RLock()
	for _, v := range this.ActorMap {
		v.SendMsg(rpc.RpcHead{}, funcName, params...)
	}
	this.ActorLock.RUnlock()
}
func (this *ActorPool) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	buff := rpc.Marshal(head, funcName, params...)
	head.SocketId = 0
	if head.Id != 0 {
		pActor := this.GetActor(head.Id)
		if pActor != nil && pActor.FindCall(funcName) != nil {
			pActor.Send(head, buff)
			return
		}
	}
	this.Send(head, buff)
}

func (this *ActorPool) PacketFunc(packet rpc.Packet) bool {
	rpcPacket, head := rpc.UnmarshalHead(packet.Buff)
	if this.FindCall(rpcPacket.FuncName) != nil {
		head.SocketId = packet.Id
		head.Reply = packet.Reply
		this.Send(head, packet.Buff)
	} else {
		pActor := this.GetActor(rpcPacket.RpcHead.Id)
		if pActor != nil && pActor.FindCall(rpcPacket.FuncName) != nil {
			head.SocketId = packet.Id
			head.Reply = packet.Reply
			pActor.Send(head, packet.Buff)
		}
	}
	return false
}
