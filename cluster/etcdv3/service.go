package etcdv3

import (
	"context"
	"encoding/json"
	"log"

	"github.com/erDong01/micro-kit/common"

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
	leaseResp, err := this.lease.Grant(context.Background(), 10)
	if err != nil {
		log.Fatalln(err)
	}
	this.leaseId = leaseResp.ID
	var (
		keepResp     *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)
	if keepRespChan, err = this.lease.KeepAlive(context.TODO(), this.leaseId); err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					log.Println("租约已经失效")
					goto END
				} else { //每秒会续租一次，所以就会受到一次应答
					// log.Println("收到自动续租应答:", keepResp.ID)
				}
			}
		}
	END:
	}()
	key := ETCD_DIR + this.String() + "/" + this.IpString()
	data, _ := json.Marshal(this.ClusterInfo)
	this.client.Put(context.Background(), key, string(data), clientv3.WithLease(this.leaseId))
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
	this.Run()
}
