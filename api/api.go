package api

type Packet struct {
	Id         uint32      `json:"Id,omitempty"`
	Reply      string      `json:"Reply,omitempty"`
	Buff       []byte      `json:"Buff,omitempty"`
	JsonPacket *JsonPacket `json:"JsonPacket,omitempty"`
}

type JsonPacket struct {
	FuncName string      `json:"FuncName,omitempty"`
	Head     *JsonHead   `json:"Head,omitempty"`
	Data     interface{} `json:"Data,omitempty"` // TODO json.RawMessage
}

type JsonHead struct {
	Id             int64  `json:"Id,omitempty"`
	SocketId       uint32 `json:"SocketId,omitempty"`
	SrcClusterId   uint32 `json:"SrcClusterId,omitempty"`
	ClusterId      uint32 `json:"ClusterId,omitempty"`
	DestServerType int32  `json:"DestServerType,omitempty"`
	SendType       int32  `json:"SendType,omitempty"`
	ActorName      string `json:"ActorName,omitempty"`
	Reply          string `json:"Reply,omitempty"`
	Code           int32  `json:"Code,omitempty"`
	Msg            string `json:"Msg,omitempty"`
	Token          string `json:"Token,omitempty"`
}
