package common

import (
	"fmt"
	"strings"

	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
)

type (
	ClusterInfo  rpc.ClusterInfo
	IClusterInfo interface {
		Id() uint32
		String() string
		ServiceType() rpc.SERVICE
		IpString() string
		RaftIp() string
	}
)

func (clusterInfo *ClusterInfo) IpString() string {
	return fmt.Sprintf("%s:%d", clusterInfo.Ip, clusterInfo.Port)
}

func (clusterInfo *ClusterInfo) RaftIp() string {
	return fmt.Sprintf("%s:%d", clusterInfo.Ip, clusterInfo.Port+10000)
}

func (clusterInfo *ClusterInfo) String() string {
	return strings.ToLower(clusterInfo.Type.String())
}

func (clusterInfo *ClusterInfo) Id() uint32 {
	return tools.ToHash(clusterInfo.IpString())
}

func (clusterInfo *ClusterInfo) ServiceType() rpc.SERVICE {
	return clusterInfo.Type
}
