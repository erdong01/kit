package cluster

import (
	"fmt"
	"github.com/erDong01/micro-kit/cluster/etcdv3"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/tools"
	"strings"
)

const (
	ETCD_DIR = "server/"
)

type (
	Service      etcdv3.Service
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
