package actor

import (
	"log"
	"strings"

	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
)

type (
	ActorMgr struct {
		ActorMap map[string]IActor
	}
	IActorMgr interface {
		Init()
		AddActor(IActor, ...string)
		GetActor(string) IActor
		InitActorHandle(ICluster)
		SendMsg(rpc.RpcHead, string, ...interface{})
	}
	ICluster interface {
		BindPacketFunc(packetFunc network.PacketFunc)
	}
)

func (this *ActorMgr) Init() {
	this.ActorMap = make(map[string]IActor)
}
func (this *ActorMgr) AddActor(pActor IActor, names ...string) {
	name := ""
	if len(names) == 0 {
		name = tools.GetClassName(pActor)
		_, exist := this.ActorMap[name]
		if exist {
			log.Printf("Register an existed GobalActor")
			return
		}
	} else {
		name = names[0]
	}

	this.ActorMap[name] = pActor
}

func (this *ActorMgr) GetActor(name string) IActor {
	name = strings.ToLower(name)
	pActor, exist := this.ActorMap[name]
	if exist {
		return pActor
	}
	return nil
}
func (this *ActorMgr) InitActorHandle(pCluster ICluster) {
	for _, v := range this.ActorMap {
		pCluster.BindPacketFunc(v.PacketFunc)
	}
}
func (this *ActorMgr) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	name := strings.ToLower(head.ActorName)
	pActor, exist := this.ActorMap[name]
	if exist {
		pActor.SendMsg(head, funcName, params...)
	}
}

var (
	MGR *ActorMgr
)

func init() {
	MGR = &ActorMgr{}
	MGR.Init()
}
