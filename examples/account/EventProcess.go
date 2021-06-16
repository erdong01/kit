package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/examples/message"
)

type (
	EventProcess struct {
		actor.Actor
		m_db *sql.DB
	}

	IEventProcess interface {
		actor.IActor
	}
)

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	//创建账号
	this.RegisterCall("C_A_RegisterRequest", func(ctx context.Context, packet *message.C_A_RegisterRequest) {
		fmt.Println("创建账号")
	})

	//登录账号
	this.RegisterCall("C_A_LoginRequest", func(ctx context.Context, packet *message.C_A_LoginRequest) {
		fmt.Println("登录账号EventProcess")
	})

	//创建玩家
	this.RegisterCall("W_A_CreatePlayer", func(ctx context.Context, accountId int64, playername string, sex int32, gClusterId uint32) {
		fmt.Println("创建玩家")
	})

	//删除玩家
	this.RegisterCall("W_A_DeletePlayer", func(ctx context.Context, accountId int64, playerId int64) {
		fmt.Println("删除玩家")
	})

	this.RegisterCall("test", func(ctx context.Context, aa int, bb string) (error, int, string) {
		return errors.New("test"), aa, bb
	})

	this.Actor.Start()
}
