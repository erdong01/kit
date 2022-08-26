package rpc

import (
	"reflect"
	"strings"
)

type (
	ICluster interface {
		SendMsg(head RpcHead, funcName string, params ...interface{})
		Call(parmas ...interface{})
		Id() uint32
	}
)

var MGR ICluster

func Route(head *RpcHead, funcName string) string {
	serverArgs := strings.Split(funcName, "<-")
	if len(serverArgs) == 2 {
		switch strings.ToLower(serverArgs[0]) {
		case "client":
			head.DestServerType = SERVICE_CLIENT
			// case "gate":
			// 	head.DestServerType = SERVICE_GATE
			// case "gm":
			// 	head.DestServerType = SERVICE_GM
			// case "game":
			// 	head.DestServerType = SERVICE_GAME
			// case "zone":
			// 	head.DestServerType = SERVICE_ZONE
			// case "db":
			// 	head.DestServerType = SERVICE_DB
		}
		funcName = serverArgs[1]
	}

	actorArgs := strings.Split(funcName, ".")
	if len(actorArgs) == 2 {
		head.ActorName = actorArgs[0]
		funcName = actorArgs[1]
	}

	return funcName
}

// params[0]:rpc.RpcHead
// params[1]:error
func Call(parmas ...interface{}) {
	/*head := *parmas[0].(*RpcHead)
	if parmas[1] == nil{
		parmas[1] = ""
	}else{
		parmas[1] = parmas[1].(error).Error()
	}*/
}

var GCall = reflect.ValueOf(Call)
