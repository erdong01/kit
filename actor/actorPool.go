package actor

import (
	"sync"

	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
)

type ActorPool struct {
	Actor
	actorMap map[int64]IActor
	ctorLock *sync.RWMutex
}
type IActorPool interface {
	GetActor(Id int64) IActor                         //获取actor
	AddActor(Id int64, pActor IActor)                 //添加actor
	DelActor(Id int64)                                //删除actor
	BoardCast(funcName string, params ...interface{}) //广播actor
	GetActorNum() int
}

func (this *ActorPool) Init(chanNum int) {
	this.actorMap = make(map[int64]IActor)
	this.ctorLock = &sync.RWMutex{}
	this.Actor.Init()
}

func (this *ActorPool) GetActor(Id int64) IActor {
	this.ctorLock.RLock()
	pActor, bEx := this.actorMap[Id]
	this.ctorLock.RUnlock()
	if bEx {
		return pActor
	}
	return nil
}

func (this *ActorPool) AddActor(Id int64, actor IActor) {
	this.ctorLock.Lock()
	this.actorMap[Id] = actor
	this.ctorLock.Unlock()
}

func (this *ActorPool) DelActor(Id int64) {
	this.ctorLock.Lock()
	delete(this.actorMap, Id)
	this.ctorLock.Unlock()
}

func (this *ActorPool) GetActorNum() int {
	nLen := 0
	this.ctorLock.RLock()
	nLen = len(this.actorMap)
	this.ctorLock.RUnlock()
	return nLen
}

func (this *ActorPool) BoardCast(funcName string, params ...interface{}) {
	this.ctorLock.RLock()
	for _, v := range this.actorMap {
		v.SendMsg(rpc3.RpcHead{}, funcName, params...)
	}
	this.ctorLock.RUnlock()
}
func (this *ActorPool) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
	buff := rpc.Marshal(head, funcName, params...)
	head.SocketId = 0
	if head.Id != 0 {
		pActor := this.GetActor(head.Id)
		if pActor != nil && pActor.HasRpc(funcName) {
			pActor.Send(head, buff)
			return
		}
	}
	this.Send(head, buff)
}

func (this *ActorPool) PacketFunc(packet rpc3.Packet) bool {
	rpcPacket, head := rpc.UnmarshalHead(packet.Buff)
	if this.HasRpc(rpcPacket.FuncName) {
		head.SocketId = packet.Id
		head.Reply = packet.Reply
		this.Send(head, packet)
	} else {
		pActor := this.GetActor(rpcPacket.RpcHead.Id)
		if pActor != nil && pActor.HasRpc(rpcPacket.FuncName) {
			head.SocketId = packet.Id
			head.Reply = packet.Reply
			pActor.Send(head, packet)
		}
	}
	return false
}
