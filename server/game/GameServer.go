package game

import (
	"github.com/erdong01/kit/base"
	"github.com/erdong01/kit/common"
	"github.com/erdong01/kit/common/cluster"
	"github.com/erdong01/kit/network"
)

type (
	ServerMgr struct {
		service    *network.ServerSocket
		isInited   bool
		snowFilake *cluster.Snowflake
	}
	IServerMgr interface {
		Init() bool
		InitDB() bool
		GetServer() *network.ServerSocket
	}
	Config struct {
		common.Server    `yaml:"game"`
		common.Db        `yaml:"DB"`
		common.Etcd      `yaml:"etcd"`
		common.SnowFlake `yaml:"snowflake"`
		common.Raft      `yaml:"raft"`
		common.Nats      `yaml:"nats"`
		common.Stub      `yaml:"stub"`
	}
)

var (
	CONF   Config
	SERVER ServerMgr
	RdID   int
)

func (s *ServerMgr) Init() bool {
	if s.isInited {
		return true
	}
	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)
	ShowMessage := func() {
		base.LOG.Println("**********************************************************")
		base.LOG.Printf("\tGAME Version:\t%s", base.BUILD_NO)
		base.LOG.Printf("\tGAME IP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
		base.LOG.Printf("\tDBServer(LAN):\t%s", CONF.Db.Ip)
		base.LOG.Printf("\tDBName:\t\t%s", CONF.Db.Name)
		base.LOG.Println("**********************************************************")
	}
	ShowMessage()
	//初始化socket
	s.service = new(network.ServerSocket)
	s.service.Init(CONF.Server.Ip, CONF.Server.Port)
	s.service.Start()

	return false
}
