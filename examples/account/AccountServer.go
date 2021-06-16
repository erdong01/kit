package account

import (
	"database/sql"
	"github.com/erDong01/micro-kit/cluster"
	"github.com/erDong01/micro-kit/cluster/common"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
)

var (
	UserNetIP     string
	UserNetPort   string
	WorkID        int
	DB_Server     string
	DB_Name       string
	DB_UserId     string
	DB_Password   string
	EtcdEndpoints []string
	Nats_Cluster  string
	SERVER        ServerMgr
)

type (
	ServerMgr struct {
		service    *network.ServerSocket
		cluster    *cluster.Cluster
		actorDB    *sql.DB
		inited     bool
		accountMgr *AccountMgr
		snowFlake  *cluster.Snowflake
	}
	IServerMgr interface {
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetServer() *network.ServerSocket
		GetCluster() *cluster.Cluster
		GetAccountMgr() *AccountMgr
	}
)

func (this *ServerMgr) Init() bool {
	if this.inited {
		return true
	}
	userNetIP := "192.168.2.231"
	port := 31701
	etcdEndpoints := []string{"192.168.2.129:2379"}
	Nats_Cluster := "192.168.2.129:4222"
	//初始化socket
	this.service = new(network.ServerSocket)
	this.service.Init(UserNetIP, port)
	this.service.Start()
	this.accountMgr = new(AccountMgr)
	this.accountMgr.Init(1000)

	//本身账号集群管理
	this.cluster = new(cluster.Cluster)
	this.cluster.Init(1000, &common.ClusterInfo{Type: rpc3.SERVICE_ACCOUNTSERVER, Ip: userNetIP, Port: int32(port)}, etcdEndpoints, Nats_Cluster)
	var packet EventProcess
	packet.Init(1000)
	this.cluster.BindPacketFunc(packet.PacketFunc)
	this.cluster.BindPacketFunc(this.accountMgr.PacketFunc)
	this.snowFlake = cluster.NewSnowflake(etcdEndpoints)
	return true
}
