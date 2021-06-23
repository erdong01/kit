package game

import "github.com/golang/protobuf/proto"

func Marshal() {
	proto.Marshal(&ConnectReq{})
}
func Unmarshal(buff []byte) {
	proto.Unmarshal(buff, &ConnectReq{})
}
