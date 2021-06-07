package etcdv3

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/cluster/common"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"log"
)

type Master struct {
	serviceMap map[uint32]*common.ClusterInfo
	client     *clientv3.Client
	actor      actor.IActor
	common.IClusterInfo
}

func (this *Master) Init(info common.IClusterInfo, Endpoints []string, pActor actor.IActor) {
	cfg := clientv3.Config{
		Endpoints: Endpoints,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	this.serviceMap = make(map[uint32]*common.ClusterInfo)
	this.client = etcdClient
	this.BindActor(pActor)
	this.Start()
}
func (this *Master) Start() {
	go this.Run()
}

func (this *Master) BindActor(pActor actor.IActor) {
	this.actor = pActor
}

func (this *Master) addService(info *common.ClusterInfo) {
	this.actor.SendMsg(rpc3.RpcHead{}, "cluster_add", info)
	this.serviceMap[info.Id()] = info
}

func (this *Master) delService(info *common.ClusterInfo) {
	delete(this.serviceMap, info.Id())
	this.actor.SendMsg(rpc3.RpcHead{}, "Cluster_Del", info)
}

func (this *Master) InitService(info *common.ClusterInfo) {
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
