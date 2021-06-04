package main

import (
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/test/account"
)



func main() {
	var s network.ServerSocket
	s.Init()
	s.StartTcpServer()


	//账号管理类
	AccountMgr := new(account.AccountMgr)
	AccountMgr.Init(1000)

}
