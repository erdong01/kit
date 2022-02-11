package rpc

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"

	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/wrong"
	"google.golang.org/protobuf/proto"
)

//rpc  Marshal
func Marshal(head rpc3.RpcHead, funcName string, params ...interface{}) rpc3.Packet {
	return marshal(head, funcName, params...)
}

//rpc  marshal
func marshal(head rpc3.RpcHead, funcName string, params ...interface{}) rpc3.Packet {
	defer func() {
		if err := recover(); err != nil {
			wrong.TraceCode(err)
		}
	}()

	rpcPacket := &rpc3.RpcPacket{FuncName: strings.ToLower(funcName), ArgLen: int32(len(params)), RpcHead: (*rpc3.RpcHead)(&head)}
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	for _, param := range params {
		enc.Encode(param)
	}
	rpcPacket.RpcBody = buf.Bytes()
	dat, _ := proto.Marshal(rpcPacket)
	return rpc3.Packet{Buff: dat, RpcPacket: rpcPacket}
}

//rpc  MarshalPacket
func MarshalPacket(head rpc3.RpcHead, funcName string, packet proto.Message) []byte {
	data, _ := marshalPacket(head, funcName, packet)
	return data
}

//rpc  marshal
func marshalPacket(head rpc3.RpcHead, funcName string, packet proto.Message) ([]byte, *rpc3.RpcPacket) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
		}
	}()

	rpcPacket := &rpc3.RpcPacket{FuncName: strings.ToLower(funcName), RpcHead: (*rpc3.RpcHead)(&head)}
	buff, _ := proto.Marshal(packet)
	rpcPacket.RpcBody = buff
	dat, _ := proto.Marshal(rpcPacket)
	return dat, rpcPacket
}
