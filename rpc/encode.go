package rpc

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"

	"github.com/erDong01/micro-kit/base"
	"google.golang.org/protobuf/proto"
)

// rpc  Marshal
func Marshal(head *RpcHead, funcName *string, params ...interface{}) Packet {
	return marshal(head, funcName, params...)
}

// rpc  marshal
func marshal(head *RpcHead, funcName *string, params ...interface{}) Packet {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	*funcName = Route(head, *funcName)
	rpcPacket := &RpcPacket{FuncName: *funcName, ArgLen: int32(len(params)), RpcHead: (*RpcHead)(head)}
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	for _, param := range params {
		enc.Encode(param)
	}
	rpcPacket.RpcBody = buf.Bytes()
	dat, _ := proto.Marshal(rpcPacket)
	return Packet{Buff: dat, RpcPacket: rpcPacket}
}

// rpc  MarshalPB
func marshalPB(bitstream *base.BitStream, packet proto.Message) {
	buf, _ := proto.Marshal(packet)
	nLen := len(buf)
	bitstream.WriteInt(nLen, 32)
	bitstream.WriteBits(buf, nLen<<3)
}

// rpc  MarshalPacket
func MarshalPacket(head RpcHead, funcName string, packet proto.Message) ([]byte, *RpcPacket) {
	data, rpcPacket := marshalPacket(head, funcName, packet)
	return data, rpcPacket
}

// rpc  marshal
func marshalPacket(head RpcHead, funcName string, packet proto.Message) ([]byte, *RpcPacket) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
		}
	}()

	rpcPacket := &RpcPacket{FuncName: strings.ToLower(funcName), RpcHead: (*RpcHead)(&head)}
	buff, _ := proto.Marshal(packet)
	rpcPacket.RpcBody = buff
	dat, _ := proto.Marshal(rpcPacket)
	return dat, rpcPacket
}
