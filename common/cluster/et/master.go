package et

import (
	"context"
	"encoding/json"
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/common"
	"github.com/erDong01/micro-kit/rpc"
	client "go.etcd.io/etcd/client/v3"
	"log"
)

// 监控服务器
type (
	Master struct {
		keysAPI *client.Client
		common.IClusterInfo
	}
)

// 监控服务器
func (m *Master) Init(info common.IClusterInfo, Endpoints []string, pActor actor.IActor) {
	cfg := client.Config{
		Endpoints: Endpoints,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	m.keysAPI = etcdClient
	m.Start()
	m.IClusterInfo = info
	m.InitServices()
}

func (m *Master) Start() {
	go m.Run()
}

func (m *Master) addService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster.Cluster_Add", info)
}

func (m *Master) delService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster.Cluster_Del", info)
}

func NodeToService(val []byte) *common.ClusterInfo {
	info := &common.ClusterInfo{}
	err := json.Unmarshal(val, info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (m *Master) Run() {
	watcher := m.keysAPI.Watcher(ETCD_DIR+m.String(), &client.WatcherOptions{
		Recursive: true,
	})

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			continue
		}
		if res.Action == "expire" || res.Action == "delete" {
			info := NodeToService([]byte(res.PrevNode.Value))
			m.delService(info)
		} else if res.Action == "set" || res.Action == "create" {
			info := NodeToService([]byte(res.Node.Value))
			m.addService(info)
		}
	}
}

func (m *Master) InitServices() {
	resp, err := m.keysAPI.Get(context.Background(), ETCD_DIR, &client.GetOptions{Recursive: true})
	if err == nil && (resp != nil && resp.Node != nil) {
		for _, v := range resp.Node.Nodes {
			if v != nil && v.Nodes != nil {
				for _, v1 := range v.Nodes {
					info := NodeToService([]byte(v1.Value))
					m.addService(info)
				}
			}
		}
	}
}
