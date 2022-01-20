package etcdv3

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/erDong01/micro-kit/cluster/common"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ETCD_DIR = "server/"
)

//注册服务器
type Service struct {
	*common.ClusterInfo
	client   *clientv3.Client
	lease    clientv3.Lease
	leasseId clientv3.LeaseID
}

func (this *Service) Run() {
	for {
		leaseResp, _ := this.lease.Grant(context.Background(), 10)
		this.leasseId = leaseResp.ID
		key := ETCD_DIR + this.String() + "/" + this.IpString()
		data, _ := json.Marshal(this.ClusterInfo)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		this.client.Put(ctx, key, string(data), clientv3.WithLease(this.leasseId))
		cancel()
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
