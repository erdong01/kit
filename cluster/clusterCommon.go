package cluster

import (
	"fmt"
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/cluster/common"
	"github.com/erDong01/micro-kit/cluster/etcdv3"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
)

type (
	Service   etcdv3.Service
	Master    etcdv3.Master
	Snowflake etcdv3.Snowflake
)

//NewService 注册服务器
func NewService(info *common.ClusterInfo, Endpoints []string) *Service {
	service := &etcdv3.Service{}
	service.Init(info, Endpoints)
	return (*Service)(service)
}

//NewMaster 监控服务器
func NewMaster(info common.IClusterInfo, Endpoints []string, pActor actor.IActor) *Master {
	master := &etcdv3.Master{}
	master.Init(info, Endpoints, pActor)
	return (*Master)(master)
}

//NewSnowflake uuid生成器
func NewSnowflake(Endpoints []string) *Snowflake {
	uuid := &etcdv3.Snowflake{}
	uuid.Init(Endpoints)
	return (*Snowflake)(uuid)
}
func GetRpcChannel(head rpc3.RpcHead) string {
	return fmt.Sprintf("%s/%s/%d", etcdv3.ETCD_DIR, strings.ToLower(head.DestServerType.String()), head.ClusterId)
}

func GetRpcTopicChannel(head rpc3.RpcHead) string {
	return fmt.Sprintf("%s/%s", etcdv3.ETCD_DIR, strings.ToLower(head.DestServerType.String()))
}
func GetRpcCallChannel(head rpc3.RpcHead) string {
	return fmt.Sprintf("%s/%s/call/%d", etcdv3.ETCD_DIR, strings.ToLower(head.DestServerType.String()), head.ClusterId)
}

func SetupNatsConn(connectString string, appDieChan chan bool, options ...nats.Option) (*nats.Conn, error) {
	natsOptions := append(
		options,
		nats.DisconnectHandler(func(_ *nats.Conn) {
			log.Println("disconnected from nats!")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("reconnected to nats server %s with address %s in cluster %s!", nc.ConnectedServerId(), nc.ConnectedAddr(), nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			err := nc.LastError()
			if err == nil {
				log.Println("nats connection closed with no error.")
				return
			}

			log.Fatalf("nats connection closed. reason: %q", nc.LastError())
			if appDieChan != nil {
				appDieChan <- true
			}
		}),
	)

	nc, err := nats.Connect(connectString, natsOptions...)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
func GetChannel(clusterInfo common.ClusterInfo) string {
	return fmt.Sprintf("%s/%s/%d", etcdv3.ETCD_DIR, clusterInfo.String(), clusterInfo.Id())
}

func GetTopicChannel(clusterInfo common.ClusterInfo) string {
	return fmt.Sprintf("%s/%s", etcdv3.ETCD_DIR, clusterInfo.String())
}

func GetCallChannel(clusterInfo common.ClusterInfo) string {
	return fmt.Sprintf("%s/%s/call/%d", etcdv3.ETCD_DIR, clusterInfo.String(), clusterInfo.Id())
}
