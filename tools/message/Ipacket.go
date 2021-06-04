package message

import (
	"github.com/erDong01/micro-kit/tools"
	"github.com/golang/protobuf/proto"
	"hash/crc32"
	"strings"
)

var (
	Packet_CreateFactorStringMap map[string]func() proto.Message
	Packet_CreateFactorMap       map[uint32]func() proto.Message
	Packet_CrcNamesMap           map[uint32]string
)

func GetMessageName(packet proto.Message) string {
	sType := strings.ToLower(proto.MessageName(packet))
	index := strings.Index(sType, ".")
	if index != -1 {
		sType = sType[index+1:]
	}
	return sType
}

func Encode(packet proto.Message) []byte {
	packetId := crc32.ChecksumIEEE([]byte(GetMessageName(packet)))
	buff, _ := proto.Marshal(packet)
	data := append(tools.IntToBytes(int(packetId)), buff...)
	return data
}

func Decode(buff []byte) (uint32, []byte) {
	packetId := uint32(tools.BytesToInt(buff[0:4]))
	return packetId, buff[4:]
}

func GetPakcet(packetId uint32) proto.Message {
	packetFunc, exist := Packet_CreateFactorMap[packetId]
	if exist {
		return packetFunc()
	}
	return nil
}
