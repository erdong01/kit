package netgate

import (
	"github.com/erDong01/micro-kit/cluster"
	"github.com/erDong01/micro-kit/cluster/common"
	"github.com/erDong01/micro-kit/examples/message"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/tools"
	"time"
)

type (
	ServerMgr struct {
		service        *network.ServerSocket
		inited         bool
		timeTraceTimer *time.Ticker
		cluster        *cluster.Cluster
		playerMgr      *PlayerManager
	}
	IServerMag interface {
		Init() bool
		GetServer() *network.ServerSocket
		GetCluster() *cluster.Service
		OnServerStart()
	}
)

var (
	UserNetIp     string
	UserNetPort   string
	EtcdEndpoints []string
	NatsCluster   string
	SERVER        ServerMgr
)

func (this *ServerMgr) Init() bool {
	UserNetIP := "192.168.2.231"
	port := 31700
	etcdEndpoints := []string{"192.168.2.129:2379"}
	Nats_Cluster := "192.168.2.129:4222"
	this.service = new(network.ServerSocket)
	this.service.Init(UserNetIP, port)
	this.service.SetMaxPacketLen(tools.MAX_CLIENT_PACKET)
	this.service.SetConnectType(network.CLIENT_CONNECT)
	packet := new(UserPrcoess)
	packet.Init(1000)
	this.service.BindPacketFunc(packet.PacketFunc)
	this.service.Start()
	var packet1 EventProcess
	this.cluster = new(cluster.Cluster)
	this.cluster.Init(1000, &common.ClusterInfo{Type: rpc3.SERVICE_GATESERVER, Ip: UserNetIP, Port: int32(port)}, etcdEndpoints, Nats_Cluster)
	this.cluster.BindPacketFunc(packet1.PacketFunc)
	this.cluster.BindPacketFunc(DispatchPacket)
	//初始玩家管理
	this.playerMgr = new(PlayerManager)
	this.playerMgr.Init(1000)

	message.Init()
	return true

}
func (this *ServerMgr) GetServer() *network.ServerSocket {
	return this.service
}

func (this ServerMgr) GetCluster() *cluster.Cluster {
	return this.cluster
}

func (this *ServerMgr) GetPlayerMgr() *PlayerManager {
	return this.playerMgr
}

func (this *ServerMgr) OnServerStart() {
	this.service.Start()
}
