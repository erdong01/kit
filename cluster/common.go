package cluster

import (
	"fmt"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/tools"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
)

const (
	ETCD_DIR = "server/"
)

type (
	ClusterInfo  rpc3.ClusterInfo
	IClusterInfo interface {
		Id() uint32
		String() string
		ServiceType() rpc3.SERVICE
		IpString() string
		RaftIp() string
	}
)

func (this ClusterInfo) IpString() string {
	return fmt.Sprintf("%s:%d", this.Ip, this.Port)
}

func (this *ClusterInfo) RaftIp() string {
	return fmt.Sprintf("%s:%d", this.Ip, this.Port+10000)
}

func (this *ClusterInfo) String() string {
	return strings.ToLower(this.Type.String())
}

func (this *ClusterInfo) Id() uint32 {
	return tools.ToHash(this.IpString())
}

func getRpcChannel(head rpc3.RpcHead) string {
	return fmt.Sprintf("%s/%s/%d", ETCD_DIR, strings.ToLower(head.DestServerType.String()), head.ClusterId)
}

func getRpcTopicChannel(head rpc3.RpcHead) string {
	return fmt.Sprintf("%s/%s", ETCD_DIR, strings.ToLower(head.DestServerType.String()))
}
func getRpcCallChannel(head rpc3.RpcHead) string {
	return fmt.Sprintf("%s/%s/call/%d", ETCD_DIR, strings.ToLower(head.DestServerType.String()), head.ClusterId)
}

func setupNatsConn(connectString string, appDieChan chan bool, options ...nats.Option) (*nats.Conn, error) {
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
func getChannel(clusterInfo ClusterInfo) string {
	return fmt.Sprintf("%s/%s/%d", ETCD_DIR, clusterInfo.String(), clusterInfo.Id())
}

func getTopicChannel(clusterInfo ClusterInfo) string {
	return fmt.Sprintf("%s/%s", ETCD_DIR, clusterInfo.String())
}

func getCallChannel(clusterInfo ClusterInfo) string {
	return fmt.Sprintf("%s/%s/call/%d", ETCD_DIR, clusterInfo.String(), clusterInfo.Id())
}
