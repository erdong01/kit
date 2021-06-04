package cluster

import (
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/tools"
	"github.com/erDong01/micro-kit/tools/vector"
	"sync"
)

type (
	//集群包管理
	IClusterPacket interface {
		actor.IActor
		SetClusterId(uint32)
	}
	ClusterNode struct {
		*network.ClientSocket
		*ClusterInfo
	}
	Cluster struct {
		actor.Actor
		clusterMap     map[uint32]*ClusterNode
		ClusterLocker  *sync.RWMutex
		Packet         IClusterPacket
		HashRing       *tools.HashRing
		ClusterInfoMap map[uint32]*ClusterInfo
		PacketFunction *vector.Vector
	}

	ICluster interface {
		//actor.IActor
		Init(num int, info *ClusterInfo, Endpoints []string)
		AddCluster(info *ClusterInfo)
		DelCluster(info *ClusterInfo)
		GetCluster(rpc3.RpcHead) *ClusterNode

		BindPacket(IClusterPacket)
		BindPacketFunc(network.PacketFunc)
		SendMsg(rpc3.RpcHead, string, ...interface{}) //发送给集群特定服务器
		Send(rpc3.RpcHead, []byte)                    //发送给集群特定服务器

		RandomCluster(head rpc3.RpcHead) rpc3.RpcHead ///随机分配

		sendPoint(rpc3.RpcHead, []byte)     //发送给集群特定服务器
		balanceSend(rpc3.RpcHead, []byte)   //负载给集群特定服务器
		boardCastSend(rpc3.RpcHead, []byte) //给集群广播
	}
)
