package account

import (
	"database/sql"
	"github.com/erDong01/micro-kit/cluster/common"
	cluster2 "github.com/erDong01/micro-kit/common/cluster"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/rpc"
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
		cluster    *cluster2.Cluster
		actorDB    *sql.DB
		inited     bool
		accountMgr *AccountMgr
		snowFlake  *cluster2.Snowflake
	}
	IServerMgr interface {
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetServer() *network.ServerSocket
		GetCluster() *cluster2.Cluster
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
	Nats_Cluster := "47.97.219.81:4222"
	//初始化socket
	this.service = new(network.ServerSocket)
	this.service.Init(UserNetIP, port)
	this.service.Start()
	this.accountMgr = new(AccountMgr)
	this.accountMgr.Init(1000)

	//本身账号集群管理
	this.cluster = new(cluster2.Cluster)
	this.cluster.Init(&common.ClusterInfo{Type: rpc.SERVICE_ACCOUNTSERVER, Ip: userNetIP, Port: int32(port)}, etcdEndpoints, Nats_Cluster)
	var packet EventProcess
	packet.Init(1000)
	this.cluster.BindPacketFunc(packet.PacketFunc)
	this.cluster.BindPacketFunc(this.accountMgr.PacketFunc)
	this.snowFlake = cluster2.NewSnowflake(etcdEndpoints)
	return true
}
