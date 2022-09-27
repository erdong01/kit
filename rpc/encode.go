package rpc

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"

	"google.golang.org/protobuf/proto"
)

// rpc  Marshal
func Marshal(head RpcHead, funcName string, params ...interface{}) []byte {
	data, _ := marshal(head, funcName, params...)
	return data
}

// rpc  marshal
func marshal(head RpcHead, funcName string, params ...interface{}) ([]byte, *RpcPacket) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
		}
	}()

	rpcPacket := &RpcPacket{FuncName: strings.ToLower(funcName), ArgLen: int32(len(params)), RpcHead: (*RpcHead)(&head)}
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	for _, param := range params {
		enc.Encode(param)
	}
	rpcPacket.RpcBody = buf.Bytes()
	dat, _ := proto.Marshal(rpcPacket)
	return dat, rpcPacket
}

// rpc  MarshalPacket
func MarshalPacket(head RpcHead, funcName string, packet proto.Message) []byte {
	data, _ := marshalPacket(head, funcName, packet)
	return data
}

// rpc  marshal
func marshalPacket(head RpcHead, funcName string, packet proto.Message) ([]byte, *RpcPacket) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
		}
	}()

	rpcPacket := &RpcPacket{FuncName: strings.ToLower(funcName), RpcHead: (*RpcHead)(&head)}
	buff, err := proto.Marshal(packet)
	if err != nil {
		log.Println("marshalPacket proto.Marshal(packet)", err)
	}
	rpcPacket.RpcBody = buff
	dat, _ := proto.Marshal(rpcPacket)
	return dat, rpcPacket
}
