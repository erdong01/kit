package etcdv3

import (
	"context"
	"encoding/json"
	"github.com/erDong01/micro-kit/common"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ETCD_DIR = "server/"
)

//注册服务器
type Service struct {
	*common.ClusterInfo
	client  *clientv3.Client
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
}

func (this *Service) Run() {
	for {
		leaseResp, _ := this.lease.Grant(context.Background(), 10)
		this.leaseId = leaseResp.ID
		key := ETCD_DIR + this.String() + "/" + this.IpString()
		data, _ := json.Marshal(this.ClusterInfo)
		this.client.Put(context.Background(), key, string(data), clientv3.WithLease(this.leaseId))
		time.Sleep(time.Second * 3)
	}
}

func (this *Service) Init(info *common.ClusterInfo, endpoints []string) {
	cfg := clientv3.Config{Endpoints: endpoints}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	this.client = etcdClient
	this.lease = lease
	this.ClusterInfo = info
	this.Start()
}
func (this *Service) Start() {
	go this.Run()
}
