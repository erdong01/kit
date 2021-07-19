package game

import "google.golang.org/protobuf/proto"

func Marshal() {
	proto.Marshal(&ConnectReq{})
}
func Unmarshal(buff []byte) {
	proto.Unmarshal(buff, &ConnectReq{})
}
