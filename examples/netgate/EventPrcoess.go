package netgate

import (
	"context"
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/pb/rpc3"
)

type (
	EventProcess struct {
		actor.Actor
	}

	IEventProcess interface {
		actor.IActor
	}
)

func (this *EventProcess) Init(num int) {
	this.Actor.Init()

	this.RegisterCall("A_G_Account_Login", func(ctx context.Context, accountId int64, socketId uint32) {
		SERVER.GetPlayerMgr().SendMsg(rpc3.RpcHead{}, "ADD_ACCOUNT", accountId, socketId)
	})

	this.Actor.Start()
}
