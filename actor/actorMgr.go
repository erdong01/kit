package actor

import (
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/tools"
	"log"
	"strings"
)

type (
	ActorMgr struct {
		ActorMap map[string]IActor
	}
	IActorMgr interface {
		Init()
		AddActor(IActor, ...string)
		GetActor(string) IActor
		SendMsg(rpc3.RpcHead, string, ...interface{})
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

func (this *ActorMgr) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
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
