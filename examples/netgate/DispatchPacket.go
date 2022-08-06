package netgate

import (
	"bytes"
	"encoding/gob"
	"github.com/erDong01/micro-kit/examples/message"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/wrong"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	A_C_RegisterResponse = proto.MessageName(&message.A_C_RegisterResponse{})
	A_C_LoginResponse    = proto.MessageName(&message.A_C_LoginResponse{})
)

func SendToClient(socketId uint32, packet proto.Message) {
	SERVER.GetServer().Send(rpc.RpcHead{SocketId: socketId}, message.Encode(packet))
}

func DispatchPacket(packet rpc.Packet) bool {
	defer func() {
		if err := recover(); err != nil {
			wrong.TraceCode(err)
		}
	}()
	rpcPacket, head := rpc.Unmarshal(packet.Buff)
	switch head.DestServerType {
	case rpc.SERVICE_GATESERVER:
		var messageName = ""
		buf := bytes.NewBuffer(rpcPacket.RpcBody)
		dec := gob.NewDecoder(buf)
		dec.Decode(&messageName)
		fn := protoreflect.FullName(messageName)
		messageType, _ := protoregistry.GlobalTypes.FindMessageByName(fn)
		messageType.New().Interface()
		packet := messageType.New().Interface().(proto.Message)
		dec.Decode(packet)
		buff := message.Encode(packet)
		if messageName == string(A_C_RegisterResponse) || messageName == string(A_C_LoginResponse) {
			SERVER.GetServer().Send(rpc.RpcHead{SocketId: head.SocketId}, buff)
		} else {
			socketId := SERVER.GetPlayerMgr().GetSocket(head.Id)
			SERVER.GetServer().Send(rpc.RpcHead{SocketId: socketId}, buff)
		}
	default:
		SERVER.GetCluster().Send(head, packet.Buff)
	}

	return true
}
