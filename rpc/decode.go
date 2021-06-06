package rpc

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/tools"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strings"
)

//rpc UnmarshalHead
func UnmarshalHead(buff []byte) (*rpc3.RpcPacket, rpc3.RpcHead) {
	nLen := tools.Clamp(len(buff), 0, 256)
	return Unmarshal(buff[:nLen])
}
func Unmarshal(buff []byte) (*rpc3.RpcPacket, rpc3.RpcHead) {
	rpcPacket := &rpc3.RpcPacket{}
	proto.Unmarshal(buff, rpcPacket)
	if rpcPacket.RpcHead == nil {
		rpcPacket.RpcHead = &rpc3.RpcHead{}
	}
	rpcPacket.FuncName = strings.ToLower(rpcPacket.FuncName)
	return rpcPacket, *(*rpc3.RpcHead)(rpcPacket.RpcHead)
}

//rpc Unmarshal
//pFuncType for RegisterCall func
func UnmarshalBody(rpcPacket *rpc3.RpcPacket, pFuncType reflect.Type) []interface{}{
	nCurLen := pFuncType.NumIn()
	params := make([]interface{}, nCurLen)
	buf := bytes.NewBuffer(rpcPacket.RpcBody)
	dec := gob.NewDecoder(buf)
	for i := 0; i < nCurLen; i++{
		if i == 0{
			params[0] = context.WithValue(context.Background(), "rpcHead", *(*rpc3.RpcHead)(rpcPacket.RpcHead))
			continue
		}
		val := reflect.New(pFuncType.In(i))
		if i < int(rpcPacket.ArgLen + 1) {
			dec.DecodeValue(val)
		}
		params[i] = val.Elem().Interface()
	}
	return params
}

func UnmarshalBodyCall(rpcPacket *rpc3.RpcPacket, pFuncType reflect.Type) (error, []interface{}) {
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
			params[0] = context.WithValue(context.Background(), "rpcHead", *(*rpc3.RpcHead)(rpcPacket.RpcHead))
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
