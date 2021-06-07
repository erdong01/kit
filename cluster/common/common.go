package common

import (
	"fmt"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/tools"
	"strings"
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
