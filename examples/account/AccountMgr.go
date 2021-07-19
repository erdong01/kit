package account

import (
	"context"
	"fmt"
	"github.com/erDong01/micro-kit/actor"
)

type (
	AccountMgr struct {
		actor.Actor
		accountMap     map[int64]*Account
		accountNameMap map[string]*Account
	}
	IAccountMgr interface {
		actor.IActor

		GetAccount(int64) *Account
		AddAccount(int64) *Account
		RemoveAccount(int64, bool)
		KickAccount(int64)
	}
)

var (
	ACCOUNTMGR AccountMgr
)

func (this *AccountMgr) Init(num int) {
	this.Actor.Init()
	this.accountMap = make(map[int64]*Account)
	this.accountNameMap = make(map[string]*Account)
	//this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	//账号登录处理
	this.RegisterCall("Account_Login", func(ctx context.Context, accountName string, accountId int64, socketId uint32, id uint32) {
		fmt.Println("出来了")
	})

	//账号断开连接
	this.RegisterCall("G_ClientLost", func(ctx context.Context, accountId int64) {
		fmt.Println("出去了")
	})

	this.Actor.Start()
}
