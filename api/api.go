package api

type Packet struct {
	Id         uint32
	Reply      string
	Buff       []byte
	JsonPacket *JsonPacket
}
type JsonPacket struct {
	FuncName string
	JsonHead *JsonHead
	Data     string
}
type JsonHead struct {
	Id             int64
	SocketId       uint32
	SrcClusterId   uint32
	ClusterId      uint32
	DestServerType int32
	SendType       int32
	ActorName      string
	Reply          string
	Code           int32
	Msg            string
	Token          string
}
