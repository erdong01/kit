package rpc

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/erdong01/kit/api"
	"github.com/erdong01/kit/base"
	"github.com/erdong01/kit/tools/mapstructure"
	"google.golang.org/protobuf/proto"
)

// rpc UnmarshalHead
func UnmarshalHead(buff []byte) (*RpcPacket, RpcHead) {
	nLen := base.Clamp(len(buff), 0, 256)
	return Unmarshal(buff[:nLen])
}

func Unmarshal(buff []byte) (*RpcPacket, RpcHead) {
	rpcPacket := &RpcPacket{}
	proto.Unmarshal(buff, rpcPacket)
	if rpcPacket.RpcHead == nil {
		rpcPacket.RpcHead = &RpcHead{}
	}
	// actor funcname
	/*actorArgs := strings.Split(rpcPacket.FuncName, ".")
	if len(actorArgs) == 2 {
		rpcPacket.RpcHead.ActorName = actorArgs[0]
		rpcPacket.FuncName = actorArgs[1]
	} else {
		rpcPacket.FuncName = actorArgs[0]
	}*/

	return rpcPacket, *(*RpcHead)(rpcPacket.RpcHead)
}

func UnmarshalBodyJson(jsonPacket api.JsonPacket, pFuncType reflect.Type) []interface{} {
	nCurLen := pFuncType.NumIn()
	params := make([]interface{}, nCurLen)
	if jsonPacket.JsonHead != nil {
		params[0] = context.WithValue(context.Background(), "rpcHead", *(*api.JsonHead)(jsonPacket.JsonHead))
	} else {
		params[0] = context.Background()
	}
	for i := 1; i < nCurLen; i++ {
		if i == 0 {
			continue
		}
		val := reflect.New(pFuncType.In(i))
		var ii = val.Elem().Interface()
		// json.Unmarshal(buff, &ii)
		mapstructure.Decode(jsonPacket.Data, &ii)
		fmt.Println("ii", ii)
		fmt.Println(" val.Elem().Interface()", val.Elem().Interface())
		params[i] = ii
	}
	return params
}

// rpc Unmarshal
// pFuncType for  (this *X)func(conttext, params)
func UnmarshalBody(rpcPacket *RpcPacket, pFuncType reflect.Type) []interface{} {
	nCurLen := pFuncType.NumIn()
	params := make([]interface{}, nCurLen)
	buf := bytes.NewBuffer(rpcPacket.RpcBody)
	dec := gob.NewDecoder(buf)
	for i := 1; i < nCurLen; i++ {
		if i == 1 {
			params[1] = context.WithValue(context.Background(), "rpcHead", *(*RpcHead)(rpcPacket.RpcHead))
			continue
		}

		val := reflect.New(pFuncType.In(i))
		if i < int(rpcPacket.ArgLen+2) {
			dec.DecodeValue(val)
		}
		params[i] = val.Elem().Interface()
	}
	return params
}

func UnmarshalBodyCall(rpcPacket *RpcPacket, pFuncType reflect.Type) (error, []interface{}) {
	strErr := ""
	nCurLen := pFuncType.NumIn()
	params := make([]interface{}, nCurLen)
	buf := bytes.NewBuffer(rpcPacket.RpcBody)
	dec := gob.NewDecoder(buf)
	dec.Decode(&strErr)
	if strErr != "" {
		return errors.New(strErr), params
	}
	for i := 0; i < nCurLen; i++ {
		if i == 0 {
			params[0] = context.WithValue(context.Background(), "rpcHead", *(*RpcHead)(rpcPacket.RpcHead))
			continue
		}

		val := reflect.New(pFuncType.In(i))
		if i < int(rpcPacket.ArgLen+1) {
			dec.DecodeValue(val)
		}
		params[i] = val.Elem().Interface()
	}
	return nil, params
}

// ------- 旧代码 -------------

func Unmarshal2(buff []byte) (*RpcPacket, RpcHead) {
	rpcPacket := &RpcPacket{}
	proto.Unmarshal(buff, rpcPacket)
	if rpcPacket.RpcHead == nil {
		rpcPacket.RpcHead = &RpcHead{}
	}
	rpcPacket.FuncName = strings.ToLower(rpcPacket.FuncName)
	return rpcPacket, *(*RpcHead)(rpcPacket.RpcHead)
}

func UnmarshalBody2(rpcPacket *RpcPacket, pFuncType reflect.Type) []interface{} {
	nCurLen := pFuncType.NumIn()
	params := make([]interface{}, nCurLen)
	var dec *gob.Decoder
	if rpcPacket.ArgLen > 0 {
		buf := bytes.NewBuffer(rpcPacket.RpcBody)
		dec = gob.NewDecoder(buf)
	}
	for i := 0; i < nCurLen; i++ {
		if i == 0 {
			params[0] = context.WithValue(context.Background(), "rpcHead", *(*RpcHead)(rpcPacket.RpcHead))
			continue
		}
		if i < int(rpcPacket.ArgLen+1) {
			val := reflect.New(pFuncType.In(i))
			dec.DecodeValue(val)
			params[i] = val.Elem().Interface()
		}
		if rpcPacket.ArgLen == 0 {
			val := reflect.New(pFuncType.In(i).Elem())
			m := val.Interface().(proto.Message)
			proto.Unmarshal(rpcPacket.RpcBody, m)
			params[i] = m
		}
	}
	return params
}
