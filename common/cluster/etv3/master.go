package etv3

import (
	"context"
	"encoding/json"
	"log"

	"github.com/erDong01/micro-kit/common"
	"github.com/erDong01/micro-kit/rpc"

	"github.com/erDong01/micro-kit/actor"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Master struct {
	serviceMap map[uint32]*common.ClusterInfo
	client     *clientv3.Client
	actor      actor.IActor
	common.IClusterInfo
}

// 监控服务器
func (m *Master) Init(info common.IClusterInfo, Endpoints []string) {
	cfg := clientv3.Config{
		Endpoints: Endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	m.client = etcdClient
	m.Start()
	m.IClusterInfo = info
	m.InitServices()
}
func (this *Master) Start() {
	go this.Run()
}

func (this *Master) BindActor(pActor actor.IActor) {
	this.actor = pActor
}

func (this *Master) addService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster_Add", info)
	this.serviceMap[info.Id()] = info
}

func (this *Master) delService(info *common.ClusterInfo) {
	delete(this.serviceMap, info.Id())
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster_Del", info)

}

func NodeToService(val []byte) *common.ClusterInfo {
	info := &common.ClusterInfo{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (this *Master) Run() {
	wch := this.client.Watch(context.Background(), ETCD_DIR+this.String(), clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wch {
		for _, v1 := range v.Events {
			if v1.Type.String() == "PUT" {
				info := NodeToService(v1.Kv.Value)
				this.addService(info)
			} else {
				info := NodeToService(v1.PrevKv.Value)
				this.delService(info)
			}
		}
	}
}

func (this *Master) GetServices() []*common.ClusterInfo {
	services := []*common.ClusterInfo{}
	resp, err := this.client.Get(context.Background(), ETCD_DIR+this.String(), clientv3.WithPrefix(), clientv3.WithPrevKV())
	if err == nil && (resp != nil && resp.Kvs != nil) {
		for _, v := range resp.Kvs {
			info := NodeToService(v.Value)
			if info.Id() != this.Id() {
				services = append(services, info)
			}
		}
	}
	return services
}
func (m *Master) InitServices() {
	resp, err := m.client.Get(context.Background(), ETCD_DIR, clientv3.WithPrefix())
	if err == nil && (resp != nil && resp.Kvs != nil) {
		for _, v := range resp.Kvs {
			info := NodeToService(v.Value)
			m.addService(info)
		}
	}
}
