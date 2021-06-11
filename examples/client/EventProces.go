package main

import (
	"context"
	"fmt"
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/examples/message"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
	"google.golang.org/protobuf/proto"
	"sync/atomic"
)

type (
	EventProcess struct {
		actor.Actor
		Client      *network.ClientSocket
		AccountId   int64
		PlayerId    int64
		AccountName string
		Passwd      string
		SimId       int64
		dh          tools.Dh
	}
	IEventProcess interface {
		actor.IActor
		LoginGame()
		LoginAccount()
		SendPacket(message proto.Message)
	}
)

func ToSlat(accountName string, pwd string) string {
	return fmt.Sprintf("%s__%s", accountName, pwd)
}

func ToCrc(accountName string, pwd string, buildNo string, nKey int64) uint32 {
	return tools.GetMessageCode1(fmt.Sprintf("%s_%s_%s_%d", accountName, pwd, buildNo, nKey))
}

func SendPacket(packet proto.Message) {
	buff := message.Encode(packet)
	CLIENT.Send(rpc3.RpcHead{}, buff)
}

func (this *EventProcess) SendPacket(packet proto.Message) {
	buff := message.Encode(packet)
	this.Client.Send(rpc3.RpcHead{}, buff)
}
func (this *EventProcess) PacketFunc(packet1 rpc3.Packet) bool {
	packetId, data := message.Decode(packet1.Buff)
	packet := message.GetPakcet(packetId)
	if packet == nil {
		return true
	}
	err := message.UnmarshalText(packet, data)
	if err == nil {
		this.Send(rpc3.RpcHead{}, rpc.Marshal(rpc3.RpcHead{}, message.GetMessageName(packet), packet))
		return true
	}

	return true
}
func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("W_C_SelectPlayerResponse", func(ctx context.Context, packet *message.W_C_SelectPlayerResponse) {
		this.AccountId = packet.GetAccountId()
		nLen := len(packet.GetPlayerData())
		if nLen == 0 {
			packet1 := &message.C_W_CreatePlayerRequest{PacketHead: message.BuildPacketHead(this.AccountId, rpc3.SERVICE_GATESERVER),
				PlayerName: "我是大坏蛋",
				Sex:        int32(0)}
			this.SendPacket(packet1)
		} else {
			this.PlayerId = packet.GetPlayerData()[0].GetPlayerID()
			this.LoginGame()
		}
	})

}
func (this *EventProcess) LoginGame() {
	packet1 := &message.C_W_Game_LoginRequset{PacketHead: message.BuildPacketHead(this.AccountId, rpc3.SERVICE_GATESERVER),
		PlayerId: this.PlayerId}
	this.SendPacket(packet1)
}

var (
	id int32
)

func (this *EventProcess) LoginAccount() {
	id := atomic.AddInt32(&id, 1)
	this.AccountName = fmt.Sprintf("test32%d", id)
	this.Passwd = tools.MD5(ToSlat(this.AccountName, "123456"))
	//this.AccountName = fmt.Sprintf("test%d", base.RAND.RandI(0, 7000))
	packet1 := &message.C_A_LoginRequest{PacketHead: message.BuildPacketHead(0, rpc3.SERVICE_GATESERVER),
		AccountName: this.AccountName, Password: this.Passwd, BuildNo: "1,5,1,1", Key: this.dh.ShareKey()}
	this.SendPacket(packet1)
}
func (this *EventProcess) LoginGate() {
	packet1 := &message.C_G_LoginResquest{PacketHead: message.BuildPacketHead(0, rpc3.SERVICE_GATESERVER),
		Key: this.dh.PubKey()}
	this.SendPacket(packet1)
}

var (
	PACKET *EventProcess
)
