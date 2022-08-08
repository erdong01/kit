// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: rpc3.proto

package rpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//服务器类型
type SERVICE int32

const (
	SERVICE_NONE             SERVICE = 0
	SERVICE_CLIENT           SERVICE = 1
	SERVICE_GATESERVER       SERVICE = 2  //网关,转发服务
	SERVICE_ACCOUNTSERVER    SERVICE = 3  //账号
	SERVICE_WORLDSERVER      SERVICE = 4  //世界
	SERVICE_ZONESERVER       SERVICE = 5  //地图
	SERVICE_WORLDDBSERVER    SERVICE = 6  //db
	SERVICE_ServiceWarKings  SERVICE = 7  //王者战记
	SERVICE_ServiceFruitPark SERVICE = 8  //db水果乐园
	SERVICE_GATE             SERVICE = 9  //网关,转发服务
	SERVICE_GM               SERVICE = 10 //gamemgr
	SERVICE_GAME             SERVICE = 11 //game
	SERVICE_ZONE             SERVICE = 12 //地图
	SERVICE_DB               SERVICE = 13 //db
)

// Enum value maps for SERVICE.
var (
	SERVICE_name = map[int32]string{
		0:  "NONE",
		1:  "CLIENT",
		2:  "GATESERVER",
		3:  "ACCOUNTSERVER",
		4:  "WORLDSERVER",
		5:  "ZONESERVER",
		6:  "WORLDDBSERVER",
		7:  "ServiceWarKings",
		8:  "ServiceFruitPark",
		9:  "GATE",
		10: "GM",
		11: "GAME",
		12: "ZONE",
		13: "DB",
	}
	SERVICE_value = map[string]int32{
		"NONE":             0,
		"CLIENT":           1,
		"GATESERVER":       2,
		"ACCOUNTSERVER":    3,
		"WORLDSERVER":      4,
		"ZONESERVER":       5,
		"WORLDDBSERVER":    6,
		"ServiceWarKings":  7,
		"ServiceFruitPark": 8,
		"GATE":             9,
		"GM":               10,
		"GAME":             11,
		"ZONE":             12,
		"DB":               13,
	}
)

func (x SERVICE) Enum() *SERVICE {
	p := new(SERVICE)
	*p = x
	return p
}

func (x SERVICE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SERVICE) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc3_proto_enumTypes[0].Descriptor()
}

func (SERVICE) Type() protoreflect.EnumType {
	return &file_rpc3_proto_enumTypes[0]
}

func (x SERVICE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SERVICE.Descriptor instead.
func (SERVICE) EnumDescriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{0}
}

//发送标志
type SEND int32

const (
	SEND_POINT      SEND = 0 //指定集群id
	SEND_BALANCE    SEND = 1 //负载
	SEND_BOARD_CAST SEND = 2 //广播
)

// Enum value maps for SEND.
var (
	SEND_name = map[int32]string{
		0: "POINT",
		1: "BALANCE",
		2: "BOARD_CAST",
	}
	SEND_value = map[string]int32{
		"POINT":      0,
		"BALANCE":    1,
		"BOARD_CAST": 2,
	}
)

func (x SEND) Enum() *SEND {
	p := new(SEND)
	*p = x
	return p
}

func (x SEND) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SEND) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc3_proto_enumTypes[1].Descriptor()
}

func (SEND) Type() protoreflect.EnumType {
	return &file_rpc3_proto_enumTypes[1]
}

func (x SEND) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SEND.Descriptor instead.
func (SEND) EnumDescriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{1}
}

//STUB类型
type STUB int32

const (
	STUB_Master     STUB = 0 //master
	STUB_PlayerMgr  STUB = 1 //db
	STUB_AccountMgr STUB = 2 //登录
	STUB_ChatMgr    STUB = 3 //聊天
	STUB_END        STUB = 4
)

// Enum value maps for STUB.
var (
	STUB_name = map[int32]string{
		0: "Master",
		1: "PlayerMgr",
		2: "AccountMgr",
		3: "ChatMgr",
		4: "END",
	}
	STUB_value = map[string]int32{
		"Master":     0,
		"PlayerMgr":  1,
		"AccountMgr": 2,
		"ChatMgr":    3,
		"END":        4,
	}
)

func (x STUB) Enum() *STUB {
	p := new(STUB)
	*p = x
	return p
}

func (x STUB) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (STUB) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc3_proto_enumTypes[2].Descriptor()
}

func (STUB) Type() protoreflect.EnumType {
	return &file_rpc3_proto_enumTypes[2]
}

func (x STUB) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use STUB.Descriptor instead.
func (STUB) EnumDescriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{2}
}

//邮件类型
type MAIL int32

const (
	MAIL_Player MAIL = 0 //player
)

// Enum value maps for MAIL.
var (
	MAIL_name = map[int32]string{
		0: "Player",
	}
	MAIL_value = map[string]int32{
		"Player": 0,
	}
)

func (x MAIL) Enum() *MAIL {
	p := new(MAIL)
	*p = x
	return p
}

func (x MAIL) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MAIL) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc3_proto_enumTypes[3].Descriptor()
}

func (MAIL) Type() protoreflect.EnumType {
	return &file_rpc3_proto_enumTypes[3]
}

func (x MAIL) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MAIL.Descriptor instead.
func (MAIL) EnumDescriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{3}
}

type RpcHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             int64   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	SocketId       uint32  `protobuf:"varint,2,opt,name=SocketId,proto3" json:"SocketId,omitempty"`
	SrcClusterId   uint32  `protobuf:"varint,3,opt,name=SrcClusterId,proto3" json:"SrcClusterId,omitempty"`                      //源集群id
	ClusterId      uint32  `protobuf:"varint,4,opt,name=ClusterId,proto3" json:"ClusterId,omitempty"`                            //目标集群id
	DestServerType SERVICE `protobuf:"varint,5,opt,name=DestServerType,proto3,enum=rpc.SERVICE" json:"DestServerType,omitempty"` //目标集群
	SendType       SEND    `protobuf:"varint,6,opt,name=SendType,proto3,enum=rpc.SEND" json:"SendType,omitempty"`
	ActorName      string  `protobuf:"bytes,7,opt,name=ActorName,proto3" json:"ActorName,omitempty"`
	Reply          string  `protobuf:"bytes,8,opt,name=Reply,proto3" json:"Reply,omitempty"` //call sessionid
	Code           int32   `protobuf:"varint,9,opt,name=Code,proto3" json:"Code,omitempty"`
	Msg            string  `protobuf:"bytes,10,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Token          string  `protobuf:"bytes,11,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *RpcHead) Reset() {
	*x = RpcHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcHead) ProtoMessage() {}

func (x *RpcHead) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcHead.ProtoReflect.Descriptor instead.
func (*RpcHead) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{0}
}

func (x *RpcHead) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RpcHead) GetSocketId() uint32 {
	if x != nil {
		return x.SocketId
	}
	return 0
}

func (x *RpcHead) GetSrcClusterId() uint32 {
	if x != nil {
		return x.SrcClusterId
	}
	return 0
}

func (x *RpcHead) GetClusterId() uint32 {
	if x != nil {
		return x.ClusterId
	}
	return 0
}

func (x *RpcHead) GetDestServerType() SERVICE {
	if x != nil {
		return x.DestServerType
	}
	return SERVICE_NONE
}

func (x *RpcHead) GetSendType() SEND {
	if x != nil {
		return x.SendType
	}
	return SEND_POINT
}

func (x *RpcHead) GetActorName() string {
	if x != nil {
		return x.ActorName
	}
	return ""
}

func (x *RpcHead) GetReply() string {
	if x != nil {
		return x.Reply
	}
	return ""
}

func (x *RpcHead) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *RpcHead) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *RpcHead) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RpcPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FuncName string   `protobuf:"bytes,1,opt,name=FuncName,proto3" json:"FuncName,omitempty"`
	ArgLen   int32    `protobuf:"varint,2,opt,name=ArgLen,proto3" json:"ArgLen,omitempty"`
	RpcHead  *RpcHead `protobuf:"bytes,3,opt,name=RpcHead,proto3" json:"RpcHead,omitempty"`
	RpcBody  []byte   `protobuf:"bytes,4,opt,name=RpcBody,proto3" json:"RpcBody,omitempty"`
}

func (x *RpcPacket) Reset() {
	*x = RpcPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcPacket) ProtoMessage() {}

func (x *RpcPacket) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcPacket.ProtoReflect.Descriptor instead.
func (*RpcPacket) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{1}
}

func (x *RpcPacket) GetFuncName() string {
	if x != nil {
		return x.FuncName
	}
	return ""
}

func (x *RpcPacket) GetArgLen() int32 {
	if x != nil {
		return x.ArgLen
	}
	return 0
}

func (x *RpcPacket) GetRpcHead() *RpcHead {
	if x != nil {
		return x.RpcHead
	}
	return nil
}

func (x *RpcPacket) GetRpcBody() []byte {
	if x != nil {
		return x.RpcBody
	}
	return nil
}

type F struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RpcPacket *RpcPacket `protobuf:"bytes,1,opt,name=RpcPacket,proto3" json:"RpcPacket,omitempty"`
	Id        int32      `protobuf:"varint,2,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *F) Reset() {
	*x = F{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *F) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*F) ProtoMessage() {}

func (x *F) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use F.ProtoReflect.Descriptor instead.
func (*F) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{2}
}

func (x *F) GetRpcPacket() *RpcPacket {
	if x != nil {
		return x.RpcPacket
	}
	return nil
}

func (x *F) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type T struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RpcPacket *RpcPacket `protobuf:"bytes,1,opt,name=RpcPacket,proto3" json:"RpcPacket,omitempty"`
	Id        int32      `protobuf:"varint,2,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *T) Reset() {
	*x = T{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *T) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*T) ProtoMessage() {}

func (x *T) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use T.ProtoReflect.Descriptor instead.
func (*T) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{3}
}

func (x *T) GetRpcPacket() *RpcPacket {
	if x != nil {
		return x.RpcPacket
	}
	return nil
}

func (x *T) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

//原始包
type Packet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`      //socketid
	Reply string `protobuf:"bytes,2,opt,name=Reply,proto3" json:"Reply,omitempty"` //call sessionid
	Buff  []byte `protobuf:"bytes,3,opt,name=Buff,proto3" json:"Buff,omitempty"`   //buff
}

func (x *Packet) Reset() {
	*x = Packet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Packet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Packet) ProtoMessage() {}

func (x *Packet) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Packet.ProtoReflect.Descriptor instead.
func (*Packet) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{4}
}

func (x *Packet) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Packet) GetReply() string {
	if x != nil {
		return x.Reply
	}
	return ""
}

func (x *Packet) GetBuff() []byte {
	if x != nil {
		return x.Buff
	}
	return nil
}

//集群信息
type ClusterInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     SERVICE `protobuf:"varint,1,opt,name=Type,proto3,enum=rpc.SERVICE" json:"Type,omitempty"`
	Ip       string  `protobuf:"bytes,2,opt,name=Ip,proto3" json:"Ip,omitempty"`
	Port     int32   `protobuf:"varint,3,opt,name=Port,proto3" json:"Port,omitempty"`
	Weight   int32   `protobuf:"varint,4,opt,name=Weight,proto3" json:"Weight,omitempty"`
	SocketId uint32  `protobuf:"varint,5,opt,name=SocketId,proto3" json:"SocketId,omitempty"`
	Url      string  `protobuf:"bytes,6,opt,name=Url,proto3" json:"Url,omitempty"`
}

func (x *ClusterInfo) Reset() {
	*x = ClusterInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterInfo) ProtoMessage() {}

func (x *ClusterInfo) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterInfo.ProtoReflect.Descriptor instead.
func (*ClusterInfo) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{5}
}

func (x *ClusterInfo) GetType() SERVICE {
	if x != nil {
		return x.Type
	}
	return SERVICE_NONE
}

func (x *ClusterInfo) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *ClusterInfo) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *ClusterInfo) GetWeight() int32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *ClusterInfo) GetSocketId() uint32 {
	if x != nil {
		return x.SocketId
	}
	return 0
}

func (x *ClusterInfo) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

//玩家集群信息
type PlayerClusterInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	LeaseId    int64  `protobuf:"varint,2,opt,name=LeaseId,proto3" json:"LeaseId,omitempty"`
	WClusterId uint32 `protobuf:"varint,3,opt,name=WClusterId,proto3" json:"WClusterId,omitempty"`
	ZClusterId uint32 `protobuf:"varint,4,opt,name=ZClusterId,proto3" json:"ZClusterId,omitempty"`
}

func (x *PlayerClusterInfo) Reset() {
	*x = PlayerClusterInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerClusterInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerClusterInfo) ProtoMessage() {}

func (x *PlayerClusterInfo) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerClusterInfo.ProtoReflect.Descriptor instead.
func (*PlayerClusterInfo) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{6}
}

func (x *PlayerClusterInfo) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PlayerClusterInfo) GetLeaseId() int64 {
	if x != nil {
		return x.LeaseId
	}
	return 0
}

func (x *PlayerClusterInfo) GetWClusterId() uint32 {
	if x != nil {
		return x.WClusterId
	}
	return 0
}

func (x *PlayerClusterInfo) GetZClusterId() uint32 {
	if x != nil {
		return x.ZClusterId
	}
	return 0
}

//邮箱
type MailBox struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	LeaseId   int64  `protobuf:"varint,2,opt,name=LeaseId,proto3" json:"LeaseId,omitempty"`
	MailType  MAIL   `protobuf:"varint,3,opt,name=MailType,proto3,enum=rpc.MAIL" json:"MailType,omitempty"`
	ClusterId uint32 `protobuf:"varint,5,opt,name=ClusterId,proto3" json:"ClusterId,omitempty"` //集群id
}

func (x *MailBox) Reset() {
	*x = MailBox{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailBox) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailBox) ProtoMessage() {}

func (x *MailBox) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailBox.ProtoReflect.Descriptor instead.
func (*MailBox) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{7}
}

func (x *MailBox) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MailBox) GetLeaseId() int64 {
	if x != nil {
		return x.LeaseId
	}
	return 0
}

func (x *MailBox) GetMailType() MAIL {
	if x != nil {
		return x.MailType
	}
	return MAIL_Player
}

func (x *MailBox) GetClusterId() uint32 {
	if x != nil {
		return x.ClusterId
	}
	return 0
}

//玩法集群信息
type StubMailBox struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	LeaseId   int64  `protobuf:"varint,2,opt,name=LeaseId,proto3" json:"LeaseId,omitempty"`
	StubType  STUB   `protobuf:"varint,3,opt,name=StubType,proto3,enum=rpc.STUB" json:"StubType,omitempty"`
	ClusterId uint32 `protobuf:"varint,5,opt,name=ClusterId,proto3" json:"ClusterId,omitempty"` //集群id
}

func (x *StubMailBox) Reset() {
	*x = StubMailBox{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc3_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StubMailBox) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StubMailBox) ProtoMessage() {}

func (x *StubMailBox) ProtoReflect() protoreflect.Message {
	mi := &file_rpc3_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StubMailBox.ProtoReflect.Descriptor instead.
func (*StubMailBox) Descriptor() ([]byte, []int) {
	return file_rpc3_proto_rawDescGZIP(), []int{8}
}

func (x *StubMailBox) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *StubMailBox) GetLeaseId() int64 {
	if x != nil {
		return x.LeaseId
	}
	return 0
}

func (x *StubMailBox) GetStubType() STUB {
	if x != nil {
		return x.StubType
	}
	return STUB_Master
}

func (x *StubMailBox) GetClusterId() uint32 {
	if x != nil {
		return x.ClusterId
	}
	return 0
}

var File_rpc3_proto protoreflect.FileDescriptor

var file_rpc3_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x72, 0x70, 0x63, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x72, 0x70,
	0x63, 0x22, 0xc4, 0x02, 0x0a, 0x07, 0x52, 0x70, 0x63, 0x48, 0x65, 0x61, 0x64, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x08, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x72, 0x63,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0c, 0x53, 0x72, 0x63, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x0e, 0x44,
	0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43,
	0x45, 0x52, 0x0e, 0x44, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x25, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x45, 0x4e, 0x44, 0x52, 0x08,
	0x53, 0x65, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x41, 0x63, 0x74, 0x6f,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x41, 0x63, 0x74,
	0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04,
	0x43, 0x6f, 0x64, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4d,
	0x73, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x81, 0x01, 0x0a, 0x09, 0x52, 0x70, 0x63,
	0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x75, 0x6e, 0x63, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x75, 0x6e, 0x63, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x72, 0x67, 0x4c, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x41, 0x72, 0x67, 0x4c, 0x65, 0x6e, 0x12, 0x26, 0x0a, 0x07, 0x52, 0x70,
	0x63, 0x48, 0x65, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x52, 0x70, 0x63, 0x48, 0x65, 0x61, 0x64, 0x52, 0x07, 0x52, 0x70, 0x63, 0x48, 0x65,
	0x61, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x70, 0x63, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x52, 0x70, 0x63, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x41, 0x0a, 0x01,
	0x46, 0x12, 0x2c, 0x0a, 0x09, 0x52, 0x70, 0x63, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x70, 0x63, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x09, 0x52, 0x70, 0x63, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x64, 0x22,
	0x41, 0x0a, 0x01, 0x54, 0x12, 0x2c, 0x0a, 0x09, 0x52, 0x70, 0x63, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x70,
	0x63, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x09, 0x52, 0x70, 0x63, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02,
	0x49, 0x64, 0x22, 0x42, 0x0a, 0x06, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x75, 0x66, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x42, 0x75, 0x66, 0x66, 0x22, 0x99, 0x01, 0x0a, 0x0b, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x45, 0x52, 0x56, 0x49,
	0x43, 0x45, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x70, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x6f, 0x72, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x57, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x55, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55,
	0x72, 0x6c, 0x22, 0x7d, 0x0a, 0x11, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x4c, 0x65, 0x61, 0x73, 0x65,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x49,
	0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x57, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x57, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x5a, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x5a, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x78, 0x0a, 0x07, 0x4d, 0x61, 0x69, 0x6c, 0x42, 0x6f, 0x78, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x4c, 0x65, 0x61, 0x73, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x4c,
	0x65, 0x61, 0x73, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x08, 0x4d, 0x61, 0x69, 0x6c, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x4d,
	0x41, 0x49, 0x4c, 0x52, 0x08, 0x4d, 0x61, 0x69, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x22, 0x7c, 0x0a, 0x0b, 0x53,
	0x74, 0x75, 0x62, 0x4d, 0x61, 0x69, 0x6c, 0x42, 0x6f, 0x78, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x4c, 0x65,
	0x61, 0x73, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x4c, 0x65, 0x61,
	0x73, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x08, 0x53, 0x74, 0x75, 0x62, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x54, 0x55,
	0x42, 0x52, 0x08, 0x53, 0x74, 0x75, 0x62, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x2a, 0xcf, 0x01, 0x0a, 0x07, 0x53, 0x45,
	0x52, 0x56, 0x49, 0x43, 0x45, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12,
	0x0a, 0x0a, 0x06, 0x43, 0x4c, 0x49, 0x45, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x47,
	0x41, 0x54, 0x45, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x41,
	0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x10, 0x03, 0x12, 0x0f,
	0x0a, 0x0b, 0x57, 0x4f, 0x52, 0x4c, 0x44, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x10, 0x04, 0x12,
	0x0e, 0x0a, 0x0a, 0x5a, 0x4f, 0x4e, 0x45, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x10, 0x05, 0x12,
	0x11, 0x0a, 0x0d, 0x57, 0x4f, 0x52, 0x4c, 0x44, 0x44, 0x42, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52,
	0x10, 0x06, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x57, 0x61, 0x72,
	0x4b, 0x69, 0x6e, 0x67, 0x73, 0x10, 0x07, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x46, 0x72, 0x75, 0x69, 0x74, 0x50, 0x61, 0x72, 0x6b, 0x10, 0x08, 0x12, 0x08, 0x0a,
	0x04, 0x47, 0x41, 0x54, 0x45, 0x10, 0x09, 0x12, 0x06, 0x0a, 0x02, 0x47, 0x4d, 0x10, 0x0a, 0x12,
	0x08, 0x0a, 0x04, 0x47, 0x41, 0x4d, 0x45, 0x10, 0x0b, 0x12, 0x08, 0x0a, 0x04, 0x5a, 0x4f, 0x4e,
	0x45, 0x10, 0x0c, 0x12, 0x06, 0x0a, 0x02, 0x44, 0x42, 0x10, 0x0d, 0x2a, 0x2e, 0x0a, 0x04, 0x53,
	0x45, 0x4e, 0x44, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x4f, 0x49, 0x4e, 0x54, 0x10, 0x00, 0x12, 0x0b,
	0x0a, 0x07, 0x42, 0x41, 0x4c, 0x41, 0x4e, 0x43, 0x45, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x42,
	0x4f, 0x41, 0x52, 0x44, 0x5f, 0x43, 0x41, 0x53, 0x54, 0x10, 0x02, 0x2a, 0x47, 0x0a, 0x04, 0x53,
	0x54, 0x55, 0x42, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x10, 0x00, 0x12,
	0x0d, 0x0a, 0x09, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x67, 0x72, 0x10, 0x01, 0x12, 0x0e,
	0x0a, 0x0a, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4d, 0x67, 0x72, 0x10, 0x02, 0x12, 0x0b,
	0x0a, 0x07, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x67, 0x72, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03, 0x45,
	0x4e, 0x44, 0x10, 0x04, 0x2a, 0x12, 0x0a, 0x04, 0x4d, 0x41, 0x49, 0x4c, 0x12, 0x0a, 0x0a, 0x06,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x10, 0x00, 0x32, 0x2a, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x07, 0x54, 0x65, 0x73, 0x74, 0x53,
	0x6f, 0x6e, 0x12, 0x06, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x46, 0x1a, 0x06, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x54, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2e, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc3_proto_rawDescOnce sync.Once
	file_rpc3_proto_rawDescData = file_rpc3_proto_rawDesc
)

func file_rpc3_proto_rawDescGZIP() []byte {
	file_rpc3_proto_rawDescOnce.Do(func() {
		file_rpc3_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc3_proto_rawDescData)
	})
	return file_rpc3_proto_rawDescData
}

var file_rpc3_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_rpc3_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_rpc3_proto_goTypes = []interface{}{
	(SERVICE)(0),              // 0: rpc.SERVICE
	(SEND)(0),                 // 1: rpc.SEND
	(STUB)(0),                 // 2: rpc.STUB
	(MAIL)(0),                 // 3: rpc.MAIL
	(*RpcHead)(nil),           // 4: rpc.RpcHead
	(*RpcPacket)(nil),         // 5: rpc.RpcPacket
	(*F)(nil),                 // 6: rpc.F
	(*T)(nil),                 // 7: rpc.T
	(*Packet)(nil),            // 8: rpc.Packet
	(*ClusterInfo)(nil),       // 9: rpc.ClusterInfo
	(*PlayerClusterInfo)(nil), // 10: rpc.PlayerClusterInfo
	(*MailBox)(nil),           // 11: rpc.MailBox
	(*StubMailBox)(nil),       // 12: rpc.StubMailBox
}
var file_rpc3_proto_depIdxs = []int32{
	0, // 0: rpc.RpcHead.DestServerType:type_name -> rpc.SERVICE
	1, // 1: rpc.RpcHead.SendType:type_name -> rpc.SEND
	4, // 2: rpc.RpcPacket.RpcHead:type_name -> rpc.RpcHead
	5, // 3: rpc.F.RpcPacket:type_name -> rpc.RpcPacket
	5, // 4: rpc.T.RpcPacket:type_name -> rpc.RpcPacket
	0, // 5: rpc.ClusterInfo.Type:type_name -> rpc.SERVICE
	3, // 6: rpc.MailBox.MailType:type_name -> rpc.MAIL
	2, // 7: rpc.StubMailBox.StubType:type_name -> rpc.STUB
	6, // 8: rpc.TestService.TestSon:input_type -> rpc.F
	7, // 9: rpc.TestService.TestSon:output_type -> rpc.T
	9, // [9:10] is the sub-list for method output_type
	8, // [8:9] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_rpc3_proto_init() }
func file_rpc3_proto_init() {
	if File_rpc3_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc3_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcHead); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc3_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcPacket); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc3_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*F); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc3_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*T); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc3_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Packet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc3_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc3_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerClusterInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc3_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailBox); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc3_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StubMailBox); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc3_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc3_proto_goTypes,
		DependencyIndexes: file_rpc3_proto_depIdxs,
		EnumInfos:         file_rpc3_proto_enumTypes,
		MessageInfos:      file_rpc3_proto_msgTypes,
	}.Build()
	File_rpc3_proto = out.File
	file_rpc3_proto_rawDesc = nil
	file_rpc3_proto_goTypes = nil
	file_rpc3_proto_depIdxs = nil
}
